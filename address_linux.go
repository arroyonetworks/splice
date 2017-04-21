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
