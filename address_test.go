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

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Get the IP Addresses
	//			Expect: No error
	// ------------------------------------------------------------------------

	addrs, err := splice.GetIPAddresses(config.loopbackIntf)
	if err != nil {
		t.Fatal("GetIPAddresses Returned Error: ", err)
	}

	// (2)	Search the Returned IP Addresses
	//			Expect: The loopback address is returned
	// ------------------------------------------------------------------------

	if len(addrs) == 0 {
		t.Fatal("No Addresses Returned")
	}

	found := false
	for _, addr := range addrs {
		if addr.IP.String() == IPv4LoopbackAddr.IP.String() {
			found = true
		}
	}
	if !found {
		t.Fatal("Loopback Address Not Returned")
	}
}

func TestGetIPAddresses_InvalidIntfValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Get the IP Addresses
	//			Expect: Error since the interface is invalid
	// ------------------------------------------------------------------------

	intf := &net.Interface{Index: -1}
	if _, err := splice.GetIPAddresses(intf); err == nil {
		t.Fatal("GetIPAddresses Did Not Return an Error With Invalid Interface Value")
	}
}

// ============================================================================
//	AddIPAddress
// ============================================================================

func TestAddIPAddress(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Another Loopback Address
	//			Expect: No error
	// ------------------------------------------------------------------------

	newAddr := RandomIPv4()

	if err := splice.AddIPAddress(config.loopbackIntf, newAddr); err != nil {
		t.Fatal("AddIPAddress Returned Error: ", err)
	}

	// (2)	Expect: Address is Present on Interface
	// ------------------------------------------------------------------------

	if !IntfHasAddress(t, config.loopbackIntf, newAddr) {
		t.Fatal("Loopback Does Not Have New Address")
	}
}

func TestAddIPAddress_InvalidAddressValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Another Loopback Address
	//			Expect: Error since the address is invalid
	// ------------------------------------------------------------------------

	newAddr := &net.IPNet{
		IP: net.IP{},
	}

	if err := splice.AddIPAddress(config.loopbackIntf, newAddr); err == nil {
		t.Fatal("AddIPAddress Did Not Return an Error With Invalid Address")
	}
}

// ============================================================================
//	DeleteIPAddress
// ============================================================================

func TestDeleteIPAddress(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Delete the Loopback Address
	//			Expect: No error
	// ------------------------------------------------------------------------

	if err := splice.DeleteIPAddress(config.loopbackIntf, IPv4LoopbackAddr); err != nil {
		t.Fatal("DeleteIPAddress Returned Error: ", err)
	}

	// (2)	Expect: Address is NOT Present on Interface
	// ------------------------------------------------------------------------

	if IntfHasAddress(t, config.loopbackIntf, IPv4LoopbackAddr) {
		t.Fatal("Loopback Still Has Loopback Address")
	}
}

func TestDeleteIPAddress_InvalidAddressValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Delete the Loopback Address
	//			Expect: Error since the address is invalid
	// ------------------------------------------------------------------------

	badAddr := &net.IPNet{
		IP: net.IP{},
	}

	if err := splice.DeleteIPAddress(config.loopbackIntf, badAddr); err == nil {
		t.Fatal("DeleteIPAddress Did Not Return an Error With Invalid Address")
	}
}
