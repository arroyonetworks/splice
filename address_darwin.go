/*
Copyright 2016 Arroyo Networks, LLC
Copyright 2019 Project Contributors (See AUTHORS File).

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package splice

import (
	"errors"
	"golang.org/x/sys/unix"
	"net"
	"unsafe"
)

// Provides ip address manipulation for macOS using System Calls.

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
