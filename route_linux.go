package splice

import (
	"github.com/vishvananda/netlink"
	"net"
)

// Provides route table manipulation for Linux using netlink.

// Determines if a route to the destination IP is available.
// This will always return true if a default route exists, regardless if the
// gateway can actually reach the destination.
func CanRouteTo(destination net.IP) bool {

	routes, err := netlink.RouteGet(destination)
	if err != nil {
		return false
	}

	if len(routes) > 0 {
		return true
	}

	return false
}

// Determines if the routing table has a specific entry for the given
// destination network.
func HasRoute(destination *net.IPNet) bool {

	filter := &netlink.Route{
		Dst: destination,
	}
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_ALL, filter, 0)
	if err != nil {
		return false
	}

	for _, route := range routes {

		if route.Dst == nil {
			continue
		}

		if route.Dst.String() == destination.String() {
			return true
		}
	}

	return false
}

// Adds a new route to the given IP network, routed by the given gateway.
// This is equivalent of 'ip address add <destination> via <gateway>'.
func AddRouteViaGateway(destination *net.IPNet, gateway net.IP) error {

	route := &netlink.Route{
		Dst: destination,
		Gw:  gateway,
	}
	return netlink.RouteAdd(route)
}

// Adds a new route to the given IP network, send out the given interface.
// This is equivalent of 'ip address add <destination> dev <intf.Name>'.
func AddRouteViaInterface(destination *net.IPNet, intf *net.Interface) error {

	route := &netlink.Route{
		Dst:       destination,
		LinkIndex: intf.Index,
		Scope:     netlink.SCOPE_LINK,
	}
	return netlink.RouteAdd(route)
}
