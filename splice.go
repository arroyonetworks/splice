// Copyright 2017 Arroyo Networks, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Provides a high-level library for manipulating network interfaces, links,
// and routes. Splice provides a unified interface for multiple operating
// systems.
package splice

import (
	"errors"
	"golang.org/x/sys/unix"
	"os"
)

// Syscall wrapper for calling ioctl requests
func ioctl(fd int, request int, argp uintptr) error {
	_, _, errorp := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(request), argp)
	if errorp == 0 {
		return os.NewSyscallError("ioctl", nil)
	} else {
		return os.NewSyscallError("ioctl", errors.New(errorp.Error()))
	}
}
