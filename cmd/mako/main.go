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
			fmt.Println(" Mako v0.1.0 - AI-Native Shell Orchestrator")
			return
		case "ask", "history", "stats":
			fmt.Printf(" '%s' command should be used inside Mako shell\n\n", os.Args[1])
			fmt.Println("Start Mako with: ./mako")
			fmt.Printf("Then inside Mako: mako %s <args>\n", os.Args[1])
			return
		}
	}

	runShellWrapper()
}

func showHelp() {
	help := `
 Mako - AI-Native Shell Orchestrator

USAGE:
    mako                                    Start Mako shell wrapper
    mako ask <question>                     Generate shell command from natural language
    mako history                            Show recent command history
    mako history <keyword>                  Search history by keyword
    mako history semantic <query>           Search history by meaning
    mako stats                              Show usage statistics
    mako help                               Show this help message
    mako version                            Show version

EXAMPLES:
    mako ask "find all files larger than 100MB"
    mako history semantic "compress video"
    mako history grep

INSIDE MAKO SHELL:
    Type commands normally - they're automatically saved with embeddings
    Use Ctrl+D or 'exit' to leave Mako

FEATURES:
     AI-powered command generation
     Semantic command search
     Automatic command history
     Full-text search
     Beautiful custom prompt

ENVIRONMENT:
    GEMINI_API_KEY    Your Gemini API key (required for AI features)

For more info: https://github.com/fabiobrug/mako
`
	fmt.Println(help)
}

func runShellWrapper() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	biolumeCyan := "\033[38;2;0;209;255m"
	deepAtlantic := "\033[38;2;93;173;226m"
	dorsalGrey := "\033[38;2;149;165;166m"
	reset := "\033[0m"

	fmt.Printf("%s Mako starting with shell: %s%s\n", biolumeCyan, shell, reset)

	fmt.Printf("%s", dorsalGrey)
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

	fmt.Printf("%s", deepAtlantic)
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

# Shark-themed PS1
PS1='\[\033[1;94m\]\u\[\033[0;90m\]@\[\033[0;96m\]\h\[\033[0;37m\]:\[\033[1;94m\]\w\[\033[0;90m\]=> \[\033[0m\]'

echo ""
echo " Mako shell active - type 'mako help' for commands"
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
