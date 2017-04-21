package splice

import (
	"github.com/vishvananda/netlink"
	"net"
)

// Provides ip address manipulation for Linux using netlink.

// Returns a list of IP addresses configured on the given interface.
// This is equivalent to 'ip address show <interface>'
func GetIPAddresses(intf *net.Interface) ([]*net.IPNet, error) {

	var ipAddresses []*net.IPNet

	link, err := netlink.LinkByIndex(intf.Index)
	if err != nil {
		return ipAddresses, err
	}

	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return ipAddresses, err
	}

	for _, addr := range addrs {
		ipAddresses = append(ipAddresses, addr.IPNet)
	}

	return ipAddresses, nil
}

// Adds an IP address to an interface.
// This is equivalent to 'ip address add <address> dev <intf.Name>'
func AddIPAddress(intf *net.Interface, address *net.IPNet) error {

	link, err := netlink.LinkByIndex(intf.Index)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(address.String())
	if err != nil {
		return err
	}

	return netlink.AddrAdd(link, addr)
}

// Removes an IP address from an interface.
// This is equivalent to 'ip address del <address> dev <intf.Name>'
func DeleteIPAddress(intf *net.Interface, address *net.IPNet) error {

	link, err := netlink.LinkByIndex(intf.Index)
	if err != nil {
		return err
	}

	addr, err := netlink.ParseAddr(address.String())
	if err != nil {
		return err
	}

	return netlink.AddrDel(link, addr)
}
