//go:build darwin

package main

import (
	"syscall"
	"unsafe"
)

// selectWithTimeout wraps syscall.Select with proper error handling for macOS
// On Darwin, Select returns only an error, not (int, error)
func selectWithTimeout(nfd int, readfds *syscall.FdSet, timeout *syscall.Timeval) (int, error) {
	// Use raw syscall for Darwin
	n, _, errno := syscall.Syscall6(
		syscall.SYS_SELECT,
		uintptr(nfd),
		uintptr(unsafe.Pointer(readfds)),
		0, // writefds
		0, // exceptfds
		uintptr(unsafe.Pointer(timeout)),
		0,
	)
	
	if errno != 0 {
		return 0, errno
	}
	
	return int(n), nil
}
