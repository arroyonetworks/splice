package splice_test

import (
	"github.com/vishvananda/netlink"
	"gitlab.com/ArroyoNetworks/splice"
	"net"
	"testing"
)

// ============================================================================
//	CanRouteTo
// ============================================================================

func TestCanRouteTo_ValidRoute(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Stage a Route into the Routing Table

	_, destination, _ := net.ParseCIDR("127.100.0.0/16")
	ip := net.ParseIP("127.100.5.5")

	err := addRouteToLoopback(destination.String())
	if err != nil {
		t.Fatal("Failed Staging Route to Loopback: ", err)
	}

	// (2)	Test CanRouteTo
	//			Expected: true

	if retval := splice.CanRouteTo(ip); retval != true {
		t.Error("CanRouteTo Returned 'false' for an Existing Route")
	}
}

func TestCanRouteTo_MissingRoute(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Test CanRouteTo
	//			Expected: false

	ip := net.ParseIP("127.100.5.5")
	if retval := splice.CanRouteTo(ip); retval != false {
		t.Error("CanRouteTo Returned 'true' for a Non-Existing Route")
	}
}

// ============================================================================
//	HasRoute
// ============================================================================

func TestHasRoute_ValidRoute(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Stage a Route into the Routing Table

	_, destination, _ := net.ParseCIDR("127.100.0.0/16")

	err := addRouteToLoopback(destination.String())
	if err != nil {
		t.Fatal("Failed Staging Route to Loopback: ", err)
	}

	// (2)	Test HasRoute
	//			Expected: true

	if retval := splice.HasRoute(destination); retval != true {
		t.Error("HasRoute Returned 'false' for an Existing Route")
	}
}

func TestHasRoute_MissingRoute(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Test HasRoute
	//			Expected: false

	_, destination, _ := net.ParseCIDR("127.100.0.0/16")
	if retval := splice.HasRoute(destination); retval != false {
		t.Error("HasRoute Returned 'true' for a Non-Existing Route")
	}
}

// Tests to ensure that a subnet of the route destination network does not
// match the route entry.
func TestHasRoute_SubnetRoute(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Stage a Route into the Routing Table

	_, destination, _ := net.ParseCIDR("127.100.0.0/16")

	err := addRouteToLoopback(destination.String())
	if err != nil {
		t.Fatal("Failed Staging Route to Loopback: ", err)
	}

	// (2)	Test HasRoute
	//			Expected: false

	_, subnet, _ := net.ParseCIDR("127.100.0.0/24")
	if retval := splice.HasRoute(subnet); retval != false {
		t.Error("HasRoute Returned 'true' for a Subnet of a Route")
	}
}

// ============================================================================
//	AddRouteViaGateway
// ============================================================================

func TestAddRouteViaGateway_ReachableGW(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Bring Loopback Interface Up

	if err := setLoopbackUp(); err != nil {
		t.Fatal("Unable to Set Loopback Up: ", err)
	}

	// (2)	Add a Route via AddRouteViaGateway
	//			Expect: No error

	_, destination, _ := net.ParseCIDR("100.100.0.0/16")
	gw := net.ParseIP("127.100.100.100")

	if err := splice.AddRouteViaGateway(destination, gw); err != nil {
		t.Fatal("AddRouteViaGateway Returned Error: ", err)
	}

	// (3)	Verify the Route was Added

	r := &netlink.Route{
		Dst: destination,
	}
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_ALL, r, 0)
	if err != nil {
		t.Fatal("Unable to Verify Routes: ", err)
	}

	if len(routes) == 0 {
		t.Error("No Route Found in Table")
	}

}

func TestAddRouteViaGateway_UnreachableGW(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Bring Loopback Interface Up

	if err := setLoopbackUp(); err != nil {
		t.Fatal("Unable to Set Loopback Up: ", err)
	}

	// (2)	Add Route via AddRouteViaGateway
	//			Expect: Error (Unreachable Network)

	_, destination, _ := net.ParseCIDR("100.100.0.0/16")
	gw := net.ParseIP("25.0.0.1")

	if err := splice.AddRouteViaGateway(destination, gw); err == nil {
		t.Fatal("No Error Returned for Unreachable Gateway")
	}

	// (3)	Verify Route is NOT in Routing Table

	r := &netlink.Route{
		Dst: destination,
	}
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_ALL, r, 0)
	if err != nil {
		t.Fatal("Unable to Verify Routes: ", err)
	}

	if len(routes) != 0 {
		t.Error("Matching Route Found in Table")
	}
}

// ============================================================================
//	AddRouteViaInterface
// ============================================================================

func TestAddRouteViaInterface(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Bring Loopback Interface Up and Get Loopback Interface Handle

	if err := setLoopbackUp(); err != nil {
		t.Fatal("Unable to Set Loopback Up: ", err)
	}

	lo, _ := net.InterfaceByName("lo")

	// (2)	Add Route via AddRouteViaInterface
	//			Expect: No error

	_, destination, _ := net.ParseCIDR("100.100.0.0/16")

	if err := splice.AddRouteViaInterface(destination, lo); err != nil {
		t.Fatal("AddRouteViaInterface Returned Error: ", err)
	}

	// (3)	Verify the Route was Added

	r := &netlink.Route{
		Dst: destination,
	}
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_ALL, r, 0)
	if err != nil {
		t.Fatal("Unable to Verify Routes: ", err)
	}

	if len(routes) == 0 {
		t.Error("No Route Found in Table")
	}
}
