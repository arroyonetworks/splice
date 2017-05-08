// Copyright 2017 Arroyo Networks, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splice

import (
	"github.com/vishvananda/netlink"
	"net"
)

// Provides route table manipulation for Linux using netlink.

// Determines if a route to the destination IP is available.
// This will always return true if a default route exists, regardless if the
// gateway can actually reach the destination.
func RouteExistsTo(destination net.IP) bool {

	if _, err := netlink.RouteGet(destination); err != nil {
		return false
	}
	return true
}

// Determines if the routing table has a specific entry for the given
// destination network.
func RouteHasEntry(destination *net.IPNet) bool {

	filter := &netlink.Route{
		Dst: destination,
	}
	if routes, err := netlink.RouteListFiltered(netlink.FAMILY_ALL, filter, 0); err == nil {
		for _, route := range routes {
			if route.Dst != nil {
				if route.Dst.String() == destination.String() {
					return true
				}
			}
		}
	}

	return false
}

// Adds a new route to the given IP network, routed by the given gateway.
// This is equivalent to 'ip route add <destination> via <gateway>'.
func RouteAddViaGateway(destination *net.IPNet, gateway net.IP) error {

	route := &netlink.Route{
		Dst: destination,
		Gw:  gateway,
	}
	return netlink.RouteAdd(route)
}

// Adds a new route to the given IP network, send out the given interface.
// This is equivalent to 'ip route add <destination> dev <intf.Name>'.
func RouteAddViaInterface(destination *net.IPNet, intf *net.Interface) error {

	route := &netlink.Route{
		Dst:       destination,
		LinkIndex: intf.Index,
		Scope:     netlink.SCOPE_LINK,
	}
	return netlink.RouteAdd(route)
}
