package splice

import (
	"github.com/vishvananda/netlink"
	"net"
)

// Provides network link manipulation for Linux using netlink.

// Administratively brings up the given network interface.
func LinkBringUp(intf *net.Interface) error {

	var err error
	var link netlink.Link

	if link, err = netlink.LinkByIndex(intf.Index); err == nil {
		return netlink.LinkSetUp(link)
	}

	return err

}

// Administratively brings down the given network interface.
func LinkBringDown(intf *net.Interface) error {

	var err error
	var link netlink.Link

	if link, err = netlink.LinkByIndex(intf.Index); err == nil {
		return netlink.LinkSetDown(link)
	}

	return err

}
