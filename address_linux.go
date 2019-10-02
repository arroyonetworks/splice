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
	"github.com/vishvananda/netlink"
	"net"
)

// Provides ip address manipulation for Linux using netlink.

// Returns a list of IP addresses configured on the given interface.
// This is equivalent to 'ip address show <interface>'
func AddressList(intf *net.Interface) ([]*net.IPNet, error) {

	var ipAddresses []*net.IPNet
	var err error

	var link netlink.Link
	var addrs []netlink.Addr

	if link, err = netlink.LinkByIndex(intf.Index); err == nil {
		if addrs, err = netlink.AddrList(link, netlink.FAMILY_ALL); err == nil {
			for _, addr := range addrs {
				ipAddresses = append(ipAddresses, addr.IPNet)
			}
			return ipAddresses, nil
		}
	}

	return ipAddresses, err
}

// Adds an IP address to an interface.
// This is equivalent to 'ip address add <address> dev <intf.Name>'
func AddressAdd(intf *net.Interface, address *net.IPNet) error {

	var err error

	var link netlink.Link
	var addr *netlink.Addr

	if link, err = netlink.LinkByIndex(intf.Index); err == nil {
		if addr, err = netlink.ParseAddr(address.String()); err == nil {
			return netlink.AddrAdd(link, addr)
		}
	}

	return err
}

// Removes an IP address from an interface.
// This is equivalent to 'ip address del <address> dev <intf.Name>'
func AddressDelete(intf *net.Interface, address *net.IPNet) error {

	var err error

	var link netlink.Link
	var addr *netlink.Addr

	if link, err = netlink.LinkByIndex(intf.Index); err == nil {
		if addr, err = netlink.ParseAddr(address.String()); err == nil {
			return netlink.AddrDel(link, addr)
		}
	}

	return err
}
