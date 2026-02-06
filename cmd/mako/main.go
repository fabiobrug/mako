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
		case "help", "-h", "--help":
			showHelp()
			return
		case "version", "-v", "--version":
			cyan := "\033[38;2;0;209;255m"
			lightBlue := "\033[38;2;93;173;226m"
			dimBlue := "\033[38;2;120;150;180m"
			reset := "\033[0m"
			fmt.Printf("\n%s▸ Mako - AI-Native Shell Orchestrator - v0.1.1 %s%s\n", lightBlue, cyan, reset)
			fmt.Printf("%s", dimBlue)
			return
		case "ask", "history", "stats":
			lightBlue := "\033[38;2;93;173;226m"
			cyan := "\033[38;2;0;209;255m"
			reset := "\033[0m"
			fmt.Printf("\n%sℹ  '%s' command must be used inside Mako shell%s\n\n", lightBlue, os.Args[1], reset)
			fmt.Printf("%sStart Mako with:%s %s./mako%s\n", lightBlue, reset, cyan, reset)
			fmt.Printf("%sThen inside Mako:%s %smako %s <args>%s\n\n", lightBlue, reset, cyan, os.Args[1], reset)
			return
		}
	}

	runShellWrapper()
}

func showHelp() {
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"

	help := fmt.Sprintf(`
%s╭─ Mako%s - AI-Native Shell Orchestrator
%s│%s
%s│%s %sUSAGE:%s
%s│%s   %smako%s                                Start Mako shell wrapper
%s│%s   %smako ask <question>%s                 Generate shell command from natural language
%s│%s   %smako history%s                        Show recent command history
%s│%s   %smako history <keyword>%s              Search history by keyword
%s│%s   %smako history semantic <query>%s       Search history by meaning
%s│%s   %smako stats%s                          Show usage statistics
%s│%s   %smako help%s                           Show this help message
%s│%s   %smako version%s                        Show version
%s│%s
%s│%s %sEXAMPLES:%s
%s│%s   %smako ask "find all files larger than 100MB"%s
%s│%s   %smako history semantic "compress video"%s
%s│%s   %smako history grep%s
%s│%s
%s│%s %sINSIDE MAKO SHELL:%s
%s│%s   Type commands normally - they're automatically saved with embeddings
%s│%s   Use Ctrl+D or 'exit' to leave Mako
%s│%s
%s│%s %sFEATURES:%s
%s│%s   ▸ AI-powered command generation
%s│%s   ▸ Semantic command search
%s│%s   ▸ Automatic command history
%s│%s   ▸ Full-text search
%s│%s   ▸ Beautiful custom prompt
%s│%s
%s│%s %sENVIRONMENT:%s
%s│%s   %sGEMINI_API_KEY%s    Your Gemini API key (required for AI features)
%s│%s
	%s╰─%s %shttps://github.com/fabiobrug/mako%s

`, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, dimBlue, reset,
		lightBlue, reset, dimBlue, reset,
		lightBlue, reset, dimBlue, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset, cyan, reset,
		lightBlue, reset,
		lightBlue, reset, dimBlue, reset)

	fmt.Println(help)
}

func runShellWrapper() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	cyan := "\033[38;2;0;209;255m"
	lightBlue := "\033[38;2;93;173;226m"
	dimBlue := "\033[38;2;120;150;180m"
	reset := "\033[0m"
	fmt.Printf("\n%s▸ Starting Mako%s%s\n", lightBlue, cyan, reset)
	fmt.Printf("%s", dimBlue)
	fmt.Println(`
 ███        ███      ████        ███  ███     █████████    
 ████      ████     ██████       ███ ███    ███     ███░   
 ███ ██  ██ ███    ███  ███      ██████     ███     ███░   
 ███  ████  ███   ██████████     ███ ███    ███     ███░   
 ███   ██   ███  ███      ███    ███  ███   ███     ███░   
 ███        ███ ███        ███   ███   ███   █████████░    
 ░░░        ░░░ ░░░        ░░░   ░░░   ░░░    ░░░░░░░░░    
    `)
	fmt.Print(reset)
	fmt.Printf("%s", lightBlue)
	fmt.Println(`
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢠⣾⣿⣿⣿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⡀⣀⢀⣀⣀⣀⣀⣀⣀⣀⣤⣤⣤⠤⠤⠤⠤⠤⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⡀⣼⣿⣿⣿⣿⢿⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⢰⣿⣿⣿⣿⡉⠉⠉⠉⠉⠉⠉⠉⠉⠉⣾⣿⣆⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠛⠻⠻⡿⢿⣿⣿⣇⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣶⣶⠀⠀
⠈⢹⣿⣿⣿⣿⣿⣦⣼⣷⣦⣀⠀⠀⠀⠈⠛⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠉⠘⠛⢧⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣼⣿⣿⠀⠀⠀
⠀⠀⠛⢿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣦⣤⣀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣰⣿⣿⣿⡿⠇⠀⠀⠀
⠀⠀⠀⠈⠘⠿⣿⣿⣿⣿⣿⣿⠛⢿⠛⠟⠛⠋⠋⠉⠋⠙⠛⠳⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣿⣿⣿⣿⣿⠃⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠈⠉⠻⡿⡇⢀⣤⣤⣶⣶⣶⣶⣏⠉⠉⠉⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣤⣆⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠠⠶⠿⣄⠀⠀⠀⠀⠀⠀⠀⢤⣾⣿⣿⣿⣿⣿⠙⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠱⢾⡿⣿⣿⣿⣿⣿⣿⣿⣦⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣶⣿⣿⣿⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠐⠢⠤⢤⣀⣀⣀⣀⣀⣀⣀⣀⣀⣠⣤⣤⣤⣤⣤⣭⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠛⠻⡿⣿⣿⣿⣿⣿⣿⣷⣶⣤⣤⣄⣀⣀⣀⣠⣤⣶⣾⣿⣿⣿⣿⣿⣶⣤⣄⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣀⣠⣤⣶⣿⣿⣿⣿⣿⣿⢿⠟⠿⠛⠛⠛⠙⠿⣿⣿⣿⣿⣿⣾⡀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠁⠙⠛⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣷⣶⣶⣆⣀⣀⣀⣿⠿⢿⣿⣿⣿⡿⠿⠟⠟⠛⠋⠈⠀⠀⠀⠀⠀⠀⠀⠙⠿⣿⣿⣿⣿⣷⣆⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣤⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⢿⣿⣿⣿⣿⣿⣷⣶⣦⣬⣭⣄⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠛⢿⢿⣿⣿⣀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⣿⣿⣿⣿⣿⡏⠿⠙⠃⠋⠀⠀⠙⠛⠘⠟⠿⠻⠟⢿⠿⡿⠿⡿⢽⠿⡿⠻⠷⠆⠀⠁⠉⠀⠘⠋⠛⠘⠛⠿⠻⠿⠿⠿⣿⢶⣶⣤⣤⣄⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠋⠻⠻⠦⠤
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⡿⠿⠟⠛⠁⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠉⠙⠛⠓⠒⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀`)
	fmt.Print(reset)
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
	interceptor := stream.NewInterceptor(500)
	if db != nil {
		interceptor.SetDatabase(db)
	}
	fmt.Printf("\n%s▸ Mako shell ready%s\n", lightBlue, reset)
	makoRcPath := createMakoRc()
	defer os.Remove(makoRcPath)
	cmd := exec.Command(shell, "--rcfile", makoRcPath, "-i")
	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start PTY: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		ptmx.Close()
		// Clean exit - clear the line and print exit message
		fmt.Print("\r\033[K") // Clear current line
		fmt.Printf("\n%s▸ Mako session ended%s\n", lightBlue, reset)
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
	// Exit message moved to defer block above
}

func createMakoRc() string {
	homeDir := os.Getenv("HOME")
	makoDir := filepath.Join(homeDir, ".mako")
	cmdFile := filepath.Join(makoDir, "last_command.txt")

	os.MkdirAll(makoDir, 0755)

	content := fmt.Sprintf(`
# Source user's normal bashrc first
if [ -f %s/.bashrc ]; then
    source %s/.bashrc
fi

# Mako customizations
export MAKO_ACTIVE=1
MAKO_CMD_FILE="%s"

# Create a 'mako' shell function
mako() {
    # Write command to file for interceptor to read
    echo "mako $@" > "$MAKO_CMD_FILE"
    # Print a unique marker
    echo "<<<MAKO_EXECUTE>>>"
}

# PS1
PS1='\[\033[0;36m\]\w\[\033[1;37m\] ❯ \[\033[0m\]'

echo ""
`, homeDir, homeDir, cmdFile)

	tmpFile, err := os.CreateTemp("", "makorc-*.sh")
	if err != nil {
		return ""
	}

	tmpFile.WriteString(content)
	tmpFile.Close()

	return tmpFile.Name()
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
