// Copyright 2017 Arroyo Networks, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splice

import (
	"golang.org/x/sys/unix"
	"net"
	"unsafe"
)

// Provides network link manipulation for macOS using System Calls.

func getIntfFlags(fd int, intf *net.Interface) (uint16, error) {

	var ifReq *ifReqInet4Flags
	var err error

	// First ------------------------------------------------------------------
	//	Prepare the ioctl Request Argument
	// ------------------------------------------------------------------------
	ifReq = new(ifReqInet4Flags)

	copy(ifReq.Name[:], intf.Name)

	// Second -----------------------------------------------------------------
	//	Call ioctl to get the Interface Flags
	// ------------------------------------------------------------------------
	err = ioctl(fd, unix.SIOCGIFFLAGS, uintptr(unsafe.Pointer(ifReq)))
	if err != nil {
		return 0, err
	}

	return ifReq.Flags, nil
}

// Administratively brings up the given network interface.
func LinkBringUp(intf *net.Interface) error {

	var fd int
	var flags uint16
	var ifReq *ifReqInet4Flags
	var err error

	// First ------------------------------------------------------------------
	//	Open an AF_INET Socket and Get Current Flags
	// ------------------------------------------------------------------------
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}

	flags, err = getIntfFlags(fd, intf)
	if err != nil {
		return err
	}

	// Second -----------------------------------------------------------------
	//	Prepare the ioctl Request Argument
	// ------------------------------------------------------------------------
	ifReq = new(ifReqInet4Flags)

	copy(ifReq.Name[:], intf.Name)

	ifReq.Flags = flags | uint16(unix.IFF_UP)

	// Third ------------------------------------------------------------------
	//	Call ioctl to set the Interface Flags
	// ------------------------------------------------------------------------
	return ioctl(fd, unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(ifReq)))
}

// Administratively brings down the given network interface.
func LinkBringDown(intf *net.Interface) error {

	var fd int
	var flags uint16
	var ifReq *ifReqInet4Flags
	var err error

	// First ------------------------------------------------------------------
	//	Open an AF_INET Socket and Get Current Flags
	// ------------------------------------------------------------------------
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}

	flags, err = getIntfFlags(fd, intf)
	if err != nil {
		return err
	}

	// Second -----------------------------------------------------------------
	//	Prepare the ioctl Request Argument
	// ------------------------------------------------------------------------
	ifReq = new(ifReqInet4Flags)

	copy(ifReq.Name[:], intf.Name)

	ifReq.Flags = flags & ^uint16(unix.IFF_UP)

	// Third ------------------------------------------------------------------
	//	Call ioctl to set the Interface Flags
	// ------------------------------------------------------------------------
	return ioctl(fd, unix.SIOCSIFFLAGS, uintptr(unsafe.Pointer(ifReq)))
}
