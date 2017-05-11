// Copyright 2017 Arroyo Networks, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splice

import (
	"errors"
	"golang.org/x/sys/unix"
	"net"
	"unsafe"
)

// Provides ip address manipulation for macOS using System Calls.

const (
	// Missing IOCTL Requests
	ioctl_SIOCAIFADDR_IN6 = 0x8080691a // Add IPv6 Address
	ioctl_SIOCGIFADDR_IN6 = 0xC1206921 // Get IPv6 Address
	ioctl_SIOCDIFADDR_IN6 = 0x81206919 // Del IPv6 Address

	// Other IPv6 Constants
	nd6_INFINITE_LIFETIME = 0xFFFFFFFF
	nd6_MAX_LIFETIME      = 0x7FFFFFFF
)

// struct ifreq with the ifra_addr union member
type ifReqInet4Addr struct {
	Name [unix.IFNAMSIZ]byte
	Addr unix.RawSockaddrInet4
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

// Returns a list of IP addresses configured on the given interface.
func AddressList(intf *net.Interface) ([]*net.IPNet, error) {

	var ipAddresses []*net.IPNet
	var addrs []net.Addr
	var err error

	addrs, err = intf.Addrs()
	if err != nil {
		return ipAddresses, err
	}

	for _, a := range addrs {
		ip, ipnet, err := net.ParseCIDR(a.String())
		if err != nil {
			continue
		}
		ipnet.IP = ip
		ipAddresses = append(ipAddresses, ipnet)
	}

	return ipAddresses, nil
}

// Implementation: Adds an IPv4 address to an interface.
func addressAdd4(intf *net.Interface, address *net.IPNet) error {

	var fd int
	var ifReq *ifAliasReqInet4
	var err error

	// First ------------------------------------------------------------------
	//	Open an AF_INET Socket
	// ------------------------------------------------------------------------
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}

	// Second -----------------------------------------------------------------
	//	Prepare the ioctl Request Argument
	// ------------------------------------------------------------------------
	ifReq = new(ifAliasReqInet4)

	copy(ifReq.Name[:], intf.Name)

	copy(ifReq.Addr.Addr[:], []byte(address.IP.To4()))
	ifReq.Addr.Family = unix.AF_INET
	ifReq.Addr.Len = uint8(unsafe.Sizeof(ifReq.Addr))

	copy(ifReq.PrefixMask.Addr[:], []byte(address.Mask))
	ifReq.PrefixMask.Family = unix.AF_INET
	ifReq.PrefixMask.Len = uint8(unsafe.Sizeof(ifReq.PrefixMask))

	// Third ------------------------------------------------------------------
	//	Call ioctl to set the Address
	// ------------------------------------------------------------------------
	return ioctl(fd, unix.SIOCAIFADDR, uintptr(unsafe.Pointer(ifReq)))
}

// Implementation: Adds an IPv6 address to an interface.
func addressAdd6(intf *net.Interface, address *net.IPNet) error {

	var fd int
	var ifReq *ifAliasReqInet6
	var err error

	// First ------------------------------------------------------------------
	//	Open an AF_INET6 Socket
	// ------------------------------------------------------------------------
	fd, err = unix.Socket(unix.AF_INET6, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}

	// Second -----------------------------------------------------------------
	//	Prepare the ioctl Request Argument
	// ------------------------------------------------------------------------
	ifReq = new(ifAliasReqInet6)

	copy(ifReq.Name[:], intf.Name)

	copy(ifReq.Addr.Addr[:], []byte(address.IP.To16()))
	ifReq.Addr.Family = unix.AF_INET6
	ifReq.Addr.Len = uint8(unsafe.Sizeof(ifReq.Addr))

	copy(ifReq.PrefixMask.Addr[:], []byte(address.Mask))
	ifReq.PrefixMask.Family = unix.AF_INET6
	ifReq.PrefixMask.Len = uint8(unsafe.Sizeof(ifReq.PrefixMask))

	ifReq.Lifetime.PrefixLifetime = nd6_MAX_LIFETIME
	ifReq.Lifetime.ValidLifetime = nd6_INFINITE_LIFETIME

	// Third ------------------------------------------------------------------
	//	Call ioctl to set the Address
	// ------------------------------------------------------------------------
	return ioctl(fd, ioctl_SIOCAIFADDR_IN6, uintptr(unsafe.Pointer(ifReq)))
}

// Adds an IP address to an interface.
func AddressAdd(intf *net.Interface, address *net.IPNet) error {

	if address.IP.To4() != nil {
		return addressAdd4(intf, address)
	}

	if address.IP.To16() != nil {
		return addressAdd6(intf, address)
	}

	return errors.New("Unknown address type")
}

// Implementation: Removes an IPv4 address from an interface.
func addressDelete4(intf *net.Interface, address *net.IPNet) error {

	var fd int
	var ifReq *ifReqInet4Addr
	var err error

	// First ------------------------------------------------------------------
	//	Open an AF_INET Socket
	// ------------------------------------------------------------------------
	fd, err = unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}

	// Second -----------------------------------------------------------------
	//	Prepare the ioctl Request Argument
	// ------------------------------------------------------------------------
	ifReq = new(ifReqInet4Addr)

	copy(ifReq.Name[:], intf.Name)

	copy(ifReq.Addr.Addr[:], []byte(address.IP.To4()))
	ifReq.Addr.Family = unix.AF_INET
	ifReq.Addr.Len = uint8(unsafe.Sizeof(ifReq.Addr))

	// Third ------------------------------------------------------------------
	//	Call ioctl to delete the Address
	// ------------------------------------------------------------------------
	return ioctl(fd, unix.SIOCDIFADDR, uintptr(unsafe.Pointer(ifReq)))
}

// Implementation: Removes an IPv6 address from an interface.
func addressDelete6(intf *net.Interface, address *net.IPNet) error {

	var fd int
	var ifReq *ifReqInet6Addr
	var err error

	// First ------------------------------------------------------------------
	//	Open an AF_INET6 Socket
	// ------------------------------------------------------------------------
	fd, err = unix.Socket(unix.AF_INET6, unix.SOCK_DGRAM, 0)
	if err != nil {
		return err
	}

	// Second -----------------------------------------------------------------
	//	Prepare the ioctl Request Argument
	// ------------------------------------------------------------------------
	ifReq = new(ifReqInet6Addr)

	copy(ifReq.Name[:], intf.Name)

	copy(ifReq.Addr.Addr[:], []byte(address.IP.To16()))
	ifReq.Addr.Family = unix.AF_INET6
	ifReq.Addr.Len = uint8(unsafe.Sizeof(ifReq.Addr))

	// Third ------------------------------------------------------------------
	//	Call ioctl to delete the Address
	// ------------------------------------------------------------------------
	return ioctl(fd, ioctl_SIOCDIFADDR_IN6, uintptr(unsafe.Pointer(ifReq)))
}

// Removes an IP address from an interface.
func AddressDelete(intf *net.Interface, address *net.IPNet) error {

	if address.IP.To4() != nil {
		return addressDelete4(intf, address)
	}

	if address.IP.To16() != nil {
		return addressDelete6(intf, address)
	}

	return errors.New("Unknown address type")
}
