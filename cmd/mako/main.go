package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/fabiobrug/mako.git/internal/stream"
)

func main() {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	fmt.Printf("Mako starting with shell: %s\n", shell)

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
		fmt.Printf("\nCaptured %d lines of output\n", len(lines))

		if len(lines) > 0 {
			fmt.Println("Last 5 lines captured:")
			recent := interceptor.GetRecentLines(5)
			for _, line := range recent {
				fmt.Printf(" > %s\n", line)
			}
		}
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
