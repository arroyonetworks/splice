// Copyright 2017 Arroyo Networks, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package splice_test

import (
	"gitlab.com/ArroyoNetworks/splice"
	"net"
	"testing"
)

// ============================================================================
//	RouteExistsTo
// ============================================================================

func TestRouteExistsTo_ValidRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Stage a Route into the Routing Table
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4Route(t, config.loopbackIntf)

	// (2)	Test RouteExistsTo
	//			Expected: true
	// ------------------------------------------------------------------------

	if retval := splice.RouteExistsTo(routeNet.IP); retval != true {
		t.Error("RouteExistsTo Returned 'false' for an Existing Route")
	}
}

func TestRouteExistsTo_MissingRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Test RouteExistsTo
	//			Expected: false
	// ------------------------------------------------------------------------

	ip := net.ParseIP("100.100.5.5")
	if retval := splice.RouteExistsTo(ip); retval != false {
		t.Error("RouteExistsTo Returned 'true' for a Non-Existing Route")
	}
}

// ============================================================================
//	RouteHasEntry
// ============================================================================

func TestRouteHasEntry_ValidRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Stage a Route into the Routing Table
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4Route(t, config.loopbackIntf)

	// (2)	Test RouteHasEntry
	//			Expected: true
	// ------------------------------------------------------------------------

	if retval := splice.RouteHasEntry(routeNet); retval != true {
		t.Error("RouteHasEntry Returned 'false' for an Existing Route")
	}
}

func TestRouteHasEntry_MissingRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Test RouteHasEntry
	//			Expected: false
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4()
	if retval := splice.RouteHasEntry(routeNet); retval != false {
		t.Error("RouteHasEntry Returned 'true' for a Non-Existing Route")
	}
}

func TestRouteHasEntry_InvalidDestinationValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Test RouteHasEntry
	//			Expected: false
	// ------------------------------------------------------------------------

	if retval := splice.RouteHasEntry(&net.IPNet{}); retval != false {
		t.Error("RouteHasEntry Returned 'true' for an Invalid Destination Value")
	}
}

// Tests to ensure that a subnet of the route destination network does not
// match the route entry.
func TestRouteHasEntry_SubnetRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Stage a Route into the Routing Table

	routeNet := RandomIPv4Route(t, config.loopbackIntf)

	// (2)	Test RouteHasEntry
	//			Expected: false

	subnet := &net.IPNet{
		IP:   routeNet.IP,
		Mask: net.IPv4Mask(255, 255, 255, 254),
	}
	if retval := splice.RouteHasEntry(subnet); retval != false {
		t.Error("RouteHasEntry Returned 'true' for a Subnet of a Route")
	}
}

// ============================================================================
//	RouteAddViaGateway
// ============================================================================

func TestRouteAddViaGateway_ReachableGW(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add a Route via RouteAddViaGateway
	//			Expect: No error
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4()
	gw := IPv4LoopbackAddr.IP

	if err := splice.RouteAddViaGateway(routeNet, gw); err != nil {
		t.Fatal("RouteAddViaGateway Returned Error: ", err)
	}

	// (2)	Expect: The route was added to the routing table
	// ------------------------------------------------------------------------

	if !RouteExists(t, routeNet) {
		t.Fatal("Added Route Does Not Exist in the Routing Table")
	}

}

func TestRouteAddViaGateway_UnreachableGW(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Route via RouteAddViaGateway
	//			Expect: Error (Unreachable Network)
	// ------------------------------------------------------------------------

	_, destination, _ := net.ParseCIDR("100.100.0.0/16")
	gw := net.ParseIP("25.0.0.1")

	if err := splice.RouteAddViaGateway(destination, gw); err == nil {
		t.Fatal("No Error Returned for Unreachable Gateway")
	}

	// (2)	Expect: The route was NOT added to the routing table
	// ------------------------------------------------------------------------

	if RouteExists(t, destination) {
		t.Fatal("Route Exists in Routing Table when it Shouldn't")
	}
}

// ============================================================================
//	RouteAddViaInterface
// ============================================================================

func TestRouteAddViaInterface(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Route via RouteAddViaInterface
	//			Expect: No error
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4()

	if err := splice.RouteAddViaInterface(routeNet, config.loopbackIntf); err != nil {
		t.Fatal("RouteAddViaInterface Returned Error: ", err)
	}

	// (2)	Expect: The route was added to the routing table
	// ------------------------------------------------------------------------

	if !RouteExists(t, routeNet) {
		t.Fatal("Added Route Does Not Exist in the Routing Table")
	}
}
