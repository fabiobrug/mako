package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"github.com/fabiobrug/mako.git/internal/ai"
	"github.com/fabiobrug/mako.git/internal/database"
	"github.com/fabiobrug/mako.git/internal/stream"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "ask":
			if len(os.Args) >= 3 {
				handleAskCommand(os.Args[2:])
				return
			}
		case "history":
			handleHistoryCommand(os.Args[2:])
			return
		case "stats":
			handleStatsCommand()
			return
		}
	}

	runShellWrapper()
}

func handleHistoryCommand(args []string) {
	dbPath := filepath.Join(os.Getenv("HOME"), ".mako", "history.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if len(args) > 0 && args[0] == "semantic" {
		if len(args) < 2 {
			fmt.Println("Usage: mako history semantic <query>")
			fmt.Println("Example: mako history semantic 'compress video'")
			return
		}

		handleSemanticSearch(db, args[1:])
		return
	}

	if len(args) == 0 {
		commands, err := db.GetRecentCommands(20)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}

		if len(commands) == 0 {
			fmt.Println("No command history yet. Run some commands in Mako first!")
			return
		}

		fmt.Println("Recent command history:")
		for _, cmd := range commands {
			fmt.Printf("\n[%s] %s\n",
				cmd.Timestamp.Format("2006-01-02 15:04:05"),
				cmd.Command)
			if cmd.ExitCode != 0 {
				fmt.Printf(" Exit code: %d\n", cmd.ExitCode)
			}
		}
	} else {
		query := strings.Join(args, " ")
		commands, err := db.SearchCommands(query, 10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Search error: %v\n", err)
			os.Exit(1)
		}

		if len(commands) == 0 {
			fmt.Printf("No commands found matching: %s\n", query)
			return
		}

		fmt.Printf("Found %d commands matching: '%s':\n", len(commands), query)
		for _, cmd := range commands {
			fmt.Printf("\n[%s] %s\n",
				cmd.Timestamp.Format("2006-01-02 15:04:05"),
				cmd.Command)
		}
	}
}

func handleSemanticSearch(db *database.DB, args []string) {
	query := strings.Join(args, " ")

	fmt.Printf(" Semantic search for: '%s'\n", query)
	fmt.Println(" Generating embedding...")

	embedService, err := ai.NewEmbeddingService()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	queryVec, err := embedService.Embed(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating embedding: %v\n", err)
		os.Exit(1)
	}

	queryBytes := ai.VectorToBytes(queryVec)

	commands, err := db.SearchCommandsSemantic(queryBytes, 10, 0.5)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Search error: %v\n", err)
		os.Exit(1)
	}

	if len(commands) == 0 {
		fmt.Printf("\n No commands found with semantic similarity to '%s'\n", query)
		fmt.Println(" Tip: Try running more commands first, or use regular search:")
		fmt.Printf("   mako history %s\n", query)
		return
	}

	fmt.Printf("\n Found %d semantically similar commands:\n", len(commands))
	for _, cmd := range commands {
		fmt.Printf("\n[%s] %s\n",
			cmd.Timestamp.Format("2006-01-02 15:04:05"),
			cmd.Command)
	}
}

func handleStatsCommand() {
	dbPath := filepath.Join(os.Getenv("HOME"), ".mako", "history.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	stats, err := db.GetStats()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting stats: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Mako Statistics:")
	fmt.Printf("  Total commands: %d\n", stats["total_commands"])
	fmt.Printf("  Commands today: %d\n", stats["commands_today"])
	fmt.Printf("  Avg duration: %.0fms\n", stats["avg_duration_ms"])
}

func handleAskCommand(args []string) {
	userRequest := strings.Join(args, " ")
	fmt.Printf("Mako is thinking about: %s\n", userRequest)

	client, err := ai.NewGeminiClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure GEMINI_API_KEY is set in  your environment.\n")
		os.Exit(1)
	}

	context := ai.GetSystemContext([]string{})

	command, err := client.GenerateCommand(userRequest, context)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating command: %v\n", err)
	}

	fmt.Printf("\nSuggested command:\n")
	fmt.Printf("   %s\n\n", command)

	fmt.Printf("Execute this command? [y/N]: ")
	var response string
	fmt.Scanln(&response)

	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Printf("Command not executed\n")
		return
	}

	fmt.Println("Executing...")
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Command failde: %v\n", err)
		os.Exit(1)
	}
}

func runShellWrapper() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	fmt.Printf(" Mako starting with shell: %s\n	", shell)
	fmt.Println(`
 ███        ███      ████        ███  ███     █████████    
 ████      ████     ██████       ███ ███    ███     ███░   
 ███ ██  ██ ███    ███  ███      ██████     ███     ███░   
 ███  ████  ███   ██████████     ███ ███    ███     ███░   
 ███   ██   ███  ███      ███    ███  ███   ███     ███░   
 ███        ███ ███        ███   ███   ███   █████████░    
 ░░░        ░░░ ░░░        ░░░   ░░░   ░░░    ░░░░░░░░░    
    `)
	fmt.Println(`
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢠⣾⣿⣏⠉⠉⠉⠉⠉⠉⢡⣶⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠘⠻⢿⣿⣿⣿⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⡄⠀
⠈⣿⣿⣿⣿⣦⣽⣦⡀⠀⠀⠛⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠛⢧⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⣿⠀⠀
⠀⠘⢿⣿⣿⣿⣿⣿⣿⣦⣄⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣿⠇⠀⠀
⠀⠀⠈⠻⣿⣿⣿⣿⡟⢿⠻⠛⠙⠉⠋⠛⠳⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣿⣿⣿⡟⠀⠀⠀
⠀⠀⠀⠀⠈⠙⢿⡇⣠⣤⣶⣶⣾⡉⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣰⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⠾⢇⠀⠀⠀⠀⠀⣴⣿⣿⣿⣿⠃⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠱⣿⣿⣿⣿⣿⣿⣦⡀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠤⢤⣀⣀⣀⣀⣀⣀⣠⣤⣤⣤⣬⣭⣿⣿⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠈⠛⢿⣿⣿⣿⣿⣿⣶⣤⣄⣀⣀⣠⣴⣾⣿⣿⣿⣷⣤⣀⡀⠀⠀⠀⠀⠀⠀⣀⣀⣤⣾⣿⣿⣿⣿⡿⠿⠛⠛⠻⣿⣿⣿⣿⣇⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠙⠻⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣶⣤⣤⣘⡛⠿⢿⡿⠟⠛⠉⠁⠀⠀⠀⠀⠀⠈⠻⣿⣿⣿⣦⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣴⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠿⢿⣿⣿⣿⣿⣿⣶⣦⣤⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠻⣿⣿⡄⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⠿⠛⠉⠁⠀⠈⠉⠙⠛⠛⠻⠿⠿⠿⠿⠟⠛⠃⠀⠀⠀⠉⠉⠉⠛⠛⠛⠿⠿⠿⣶⣦⣄⡀⠀⠀⠀⠀⠀⠈⠙⠛⠂
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⠿⠛⠋⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠉⠁⠀⠀⠀⠀⠀⠀⠀⠀`)

	dbPath := filepath.Join(os.Getenv("HOME"), ".mako", "history.db")
	db, err := database.NewDB(dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not open database: %v\n", err)
	}
	defer func() {
		if db != nil {
			syncBashHistory(db)
			db.Close()
		}
	}()

	workingDir, _ := os.Getwd()
	_ = workingDir

	interceptor := stream.NewInterceptor(500)

	cmd := exec.Command(shell)

	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start PTY: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		ptmx.Close()

		lines := interceptor.GetAllLines()
		fmt.Printf("\n Captured %d lines of output\n", len(lines))
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				fmt.Fprintf(os.Stderr, "Error resizing PTY: %v\n", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH

	oldState, err := MakeRaw(os.Stdin.Fd())
	if err != nil {
		panic(err)
	}
	defer Restore(os.Stdin.Fd(), oldState)

	go func() { io.Copy(ptmx, os.Stdin) }()

	interceptor.Tee(os.Stdout, ptmx)

	fmt.Println("\n Mako exiting...")
}

func syncBashHistory(db *database.DB) {
	histFile := filepath.Join(os.Getenv("HOME"), ".bash_history")

	data, err := os.ReadFile(histFile)
	if err != nil {
		return
	}

	lines := strings.Split(string(data), "\n")

	recent, err := db.GetRecentCommands(20)
	if err != nil {
		return
	}

	existing := make(map[string]bool)
	for _, cmd := range recent {
		existing[cmd.Command] = true
	}

	workingDir, _ := os.Getwd()

	embedService, _ := ai.NewEmbeddingService()

	startIdx := len(lines) - 10
	if startIdx < 0 {
		startIdx = 0
	}

	for i := startIdx; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		if existing[line] {
			continue
		}

		if line == "exit" || line == "clear" || line == "history" {
			continue
		}

		var embeddingBytes []byte
		if embedService != nil {
			vec, err := embedService.Embed(line)
			if err == nil {
				embeddingBytes = ai.VectorToBytes(vec)
			}
		}

		cmd := database.Command{
			Command:    line,
			Timestamp:  time.Now(),
			ExitCode:   0,
			Duration:   0,
			WorkingDir: workingDir,
			Embedding:  embeddingBytes,
		}

		db.SaveCommand(cmd)
		existing[line] = true
	}
}

func MakeRaw(fd uintptr) (*syscall.Termios, error) {
	termios, err := GetTermios(fd)
	if err != nil {
		return nil, err
	}

	oldState := *termios

	termios.Iflag &^= syscall.IGNBRK | syscall.BRKINT | syscall.PARMRK |
		syscall.ISTRIP | syscall.INLCR | syscall.IGNCR |
		syscall.ICRNL | syscall.IXON
	termios.Oflag &^= syscall.OPOST
	termios.Lflag &^= syscall.ECHO | syscall.ECHONL | syscall.ICANON |
		syscall.ISIG | syscall.IEXTEN
	termios.Cflag &^= syscall.CSIZE | syscall.PARENB
	termios.Cflag |= syscall.CS8
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	if err := SetTermios(fd, termios); err != nil {
		return nil, err
	}

	return &oldState, nil
}

func Restore(fd uintptr, oldState *syscall.Termios) error {
	return SetTermios(fd, oldState)
}

func GetTermios(fd uintptr) (*syscall.Termios, error) {
	termios := &syscall.Termios{}
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCGETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	if err != 0 {
		return nil, err
	}
	return termios, nil
}

func SetTermios(fd uintptr, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TCSETS, uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	if err != 0 {
		return err
	}
	return nil
}
