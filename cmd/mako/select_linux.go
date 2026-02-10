//go:build linux

package main

import "syscall"

// selectWithTimeout wraps syscall.Select with proper error handling for Linux
func selectWithTimeout(nfd int, readfds *syscall.FdSet, timeout *syscall.Timeval) (int, error) {
	return syscall.Select(nfd, readfds, nil, nil, timeout)
}
