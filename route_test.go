package splice_test

import (
	"gitlab.com/ArroyoNetworks/splice"
	"net"
	"testing"
)

// ============================================================================
//	CanRouteTo
// ============================================================================

func TestCanRouteTo_ValidRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Stage a Route into the Routing Table
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4Route(t, config.loopbackIntf)

	// (2)	Test CanRouteTo
	//			Expected: true
	// ------------------------------------------------------------------------

	if retval := splice.CanRouteTo(routeNet.IP); retval != true {
		t.Error("CanRouteTo Returned 'false' for an Existing Route")
	}
}

func TestCanRouteTo_MissingRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Test CanRouteTo
	//			Expected: false
	// ------------------------------------------------------------------------

	ip := net.ParseIP("100.100.5.5")
	if retval := splice.CanRouteTo(ip); retval != false {
		t.Error("CanRouteTo Returned 'true' for a Non-Existing Route")
	}
}

// ============================================================================
//	HasRoute
// ============================================================================

func TestHasRoute_ValidRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Stage a Route into the Routing Table
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4Route(t, config.loopbackIntf)

	// (2)	Test HasRoute
	//			Expected: true
	// ------------------------------------------------------------------------

	if retval := splice.HasRoute(routeNet); retval != true {
		t.Error("HasRoute Returned 'false' for an Existing Route")
	}
}

func TestHasRoute_MissingRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Test HasRoute
	//			Expected: false
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4()
	if retval := splice.HasRoute(routeNet); retval != false {
		t.Error("HasRoute Returned 'true' for a Non-Existing Route")
	}
}

func TestHasRoute_InvalidDestinationValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Test HasRoute
	//			Expected: false
	// ------------------------------------------------------------------------

	if retval := splice.HasRoute(&net.IPNet{}); retval != false {
		t.Error("HasRoute Returned 'true' for an Invalid Destination Value")
	}
}

// Tests to ensure that a subnet of the route destination network does not
// match the route entry.
func TestHasRoute_SubnetRoute(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Stage a Route into the Routing Table

	routeNet := RandomIPv4Route(t, config.loopbackIntf)

	// (2)	Test HasRoute
	//			Expected: false

	subnet := &net.IPNet{
		IP:   routeNet.IP,
		Mask: net.IPv4Mask(255, 255, 255, 254),
	}
	if retval := splice.HasRoute(subnet); retval != false {
		t.Error("HasRoute Returned 'true' for a Subnet of a Route")
	}
}

// ============================================================================
//	AddRouteViaGateway
// ============================================================================

func TestAddRouteViaGateway_ReachableGW(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add a Route via AddRouteViaGateway
	//			Expect: No error
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4()
	gw := IPv4LoopbackAddr.IP

	if err := splice.AddRouteViaGateway(routeNet, gw); err != nil {
		t.Fatal("AddRouteViaGateway Returned Error: ", err)
	}

	// (2)	Expect: The route was added to the routing table
	// ------------------------------------------------------------------------

	if !RouteExists(t, routeNet) {
		t.Fatal("Added Route Does Not Exist in the Routing Table")
	}

}

func TestAddRouteViaGateway_UnreachableGW(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Route via AddRouteViaGateway
	//			Expect: Error (Unreachable Network)
	// ------------------------------------------------------------------------

	_, destination, _ := net.ParseCIDR("100.100.0.0/16")
	gw := net.ParseIP("25.0.0.1")

	if err := splice.AddRouteViaGateway(destination, gw); err == nil {
		t.Fatal("No Error Returned for Unreachable Gateway")
	}

	// (2)	Expect: The route was NOT added to the routing table
	// ------------------------------------------------------------------------

	if RouteExists(t, destination) {
		t.Fatal("Route Exists in Routing Table when it Shouldn't")
	}
}

// ============================================================================
//	AddRouteViaInterface
// ============================================================================

func TestAddRouteViaInterface(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Route via AddRouteViaInterface
	//			Expect: No error
	// ------------------------------------------------------------------------

	routeNet := RandomIPv4()

	if err := splice.AddRouteViaInterface(routeNet, config.loopbackIntf); err != nil {
		t.Fatal("AddRouteViaInterface Returned Error: ", err)
	}

	// (2)	Expect: The route was added to the routing table
	// ------------------------------------------------------------------------

	if !RouteExists(t, routeNet) {
		t.Fatal("Added Route Does Not Exist in the Routing Table")
	}
}
