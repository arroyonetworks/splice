package splice_test

import (
	"gitlab.com/ArroyoNetworks/splice"
	"net"
	"testing"
)

// ============================================================================
//	GetIPAddresses
// ============================================================================

func TestGetIPAddresses(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Bring Loopback Interface Up and Get Loopback Interface Handle

	if err := setLoopbackUp(); err != nil {
		t.Fatal("Unable to Set Loopback Up: ", err)
	}

	lo, _ := net.InterfaceByName("lo")

	// (2)	Get the IP Addresses
	//			Expect: Return Loopback Address, No error

	addrs, err := splice.GetIPAddresses(lo)
	if err != nil {
		t.Fatal("GetIPAddresses Returned Error: ", err)
	}

	if len(addrs) == 0 {
		t.Fatal("No Addresses Returned")
	}

	// (3)	Verify that the Loopback Address (127.0.0.1) was Returned

	found := false
	for _, addr := range addrs {
		if addr.IP.String() == v4LoopbackAddr.IP.String() {
			found = true
		}
	}
	if !found {
		t.Fatal("Loopback Address Not Returned")
	}
}

// ============================================================================
//	AddIPAddress
// ============================================================================

func TestAddIPAddress(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Bring Loopback Interface Up and Get Loopback Interface Handle

	if err := setLoopbackUp(); err != nil {
		t.Fatal("Unable to Set Loopback Up: ", err)
	}

	lo, _ := net.InterfaceByName("lo")

	// (2)	Add Another Loopback Address
	//			Expect: No error

	newAddr := &net.IPNet{
		IP:   net.ParseIP("5.5.5.5"),
		Mask: v4LoopbackAddr.Mask,
	}

	if err := splice.AddIPAddress(lo, newAddr); err != nil {
		t.Fatal("AddIPAddress Returned Error: ", err)
	}

	// (3)	Verify the Address was Actually Added

	if !loopbackHasAddress(newAddr.String()) {
		t.Fatal("Loopback Does Not Have New Address")
	}
}

// ============================================================================
//	DeleteIPAddress
// ============================================================================

func TestDeleteIPAddress(t *testing.T) {

	teardown := setupNamespaceTest(t)
	defer teardown()

	// (1)	Bring Loopback Interface Up and Get Loopback Interface Handle

	if err := setLoopbackUp(); err != nil {
		t.Fatal("Unable to Set Loopback Up: ", err)
	}

	lo, _ := net.InterfaceByName("lo")

	// (2)	Delete the Loopback Address
	//			Expect: No error

	if err := splice.DeleteIPAddress(lo, v4LoopbackAddr); err != nil {
		t.Fatal("DeleteIPAddress Returned Error: ", err)
	}

	// (3)	Verify the Address was Actually Removed

	if loopbackHasAddress(v4LoopbackAddr.String()) {
		t.Fatal("Loopback Still Has Loopback Address")
	}

}
