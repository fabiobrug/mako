//go:build linux

package main

import "syscall"

const (
	tcgets = syscall.TCGETS
	tcsets = syscall.TCSETS
)

func getTermiosCmd() uintptr {
	return tcgets
}

func setTermiosCmd() uintptr {
	return tcsets
}
