// Copyright 2017 Arroyo Networks, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splice

import (
	"errors"
	"golang.org/x/sys/unix"
	"os"
)

const (
	// Missing IOCTL Requests
	ioctl_SIOCAIFADDR_IN6 = 0x8080691a // Add IPv6 Address
	ioctl_SIOCGIFADDR_IN6 = 0xC1206921 // Get IPv6 Address
	ioctl_SIOCDIFADDR_IN6 = 0x81206919 // Del IPv6 Address

	// Other IPv6 Constants
	nd6_INFINITE_LIFETIME = 0xFFFFFFFF
	nd6_MAX_LIFETIME      = 0x7FFFFFFF
)

// struct ifreq with the ifru_addr union member
type ifReqInet4Addr struct {
	Name [unix.IFNAMSIZ]byte
	Addr unix.RawSockaddrInet4
}

// struct ifreq with the ifru_flags union member
type ifReqInet4Flags struct {
	Name  [unix.IFNAMSIZ]byte
	Flags uint16
}

// struct in6_ifreq (see netinet6/in6_var.h) with the ifru_addr union member
type ifReqInet6Addr struct {
	Name [unix.IFNAMSIZ]byte
	Addr unix.RawSockaddrInet6
}

// struct ifalias_req
type ifAliasReqInet4 struct {
	Name       [unix.IFNAMSIZ]byte
	Addr       unix.RawSockaddrInet4
	Broadcast  unix.RawSockaddrInet4
	PrefixMask unix.RawSockaddrInet4
}

// struct in6_addrlifetime (see netinet6/in6_var.h)
type addrLifetimeInet6 struct {
	Expire         int64
	Preferred      int64
	ValidLifetime  uint32
	PrefixLifetime uint32
}

// struct in6_aliasreq (see netinet6/in6_var.h)
type ifAliasReqInet6 struct {
	Name       [unix.IFNAMSIZ]byte
	Addr       unix.RawSockaddrInet6
	DstAddr    unix.RawSockaddrInet6
	PrefixMask unix.RawSockaddrInet6
	Flags      int32
	Lifetime   addrLifetimeInet6
}

// Syscall wrapper for calling ioctl requests
func ioctl(fd int, request int, argp uintptr) error {
	_, _, errorp := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(request), argp)
	if errorp == 0 {
		return os.NewSyscallError("ioctl", nil)
	} else {
		return os.NewSyscallError("ioctl", errors.New(errorp.Error()))
	}
}
