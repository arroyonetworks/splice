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
