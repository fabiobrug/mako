//go:build linux || darwin

package main

import (
	"syscall"
	"unsafe"
)

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
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, getTermiosCmd(), uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	if err != 0 {
		return nil, err
	}
	return termios, nil
}

func SetTermios(fd uintptr, termios *syscall.Termios) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, setTermiosCmd(), uintptr(unsafe.Pointer(termios)), 0, 0, 0)

	if err != 0 {
		return err
	}
	return nil
}

func FD_ZERO(set *syscall.FdSet) {
	for i := range set.Bits {
		set.Bits[i] = 0
	}
}

func FD_SET(fd int, set *syscall.FdSet) {
	set.Bits[fd/64] |= 1 << (uint(fd) % 64)
}
