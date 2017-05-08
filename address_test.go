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
//	AddressList
// ============================================================================

func TestAddressList(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Get the IP Addresses
	//			Expect: No error
	// ------------------------------------------------------------------------

	addrs, err := splice.AddressList(config.loopbackIntf)
	if err != nil {
		t.Fatal("AddressList Returned Error: ", err)
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

func TestAddressList_InvalidIntfValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Get the IP Addresses
	//			Expect: Error since the interface is invalid
	// ------------------------------------------------------------------------

	intf := &net.Interface{Index: -1}
	if _, err := splice.AddressList(intf); err == nil {
		t.Fatal("AddressList Did Not Return an Error With Invalid Interface Value")
	}
}

// ============================================================================
//	AddressAdd
// ============================================================================

func TestAddressAdd(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Another Loopback Address
	//			Expect: No error
	// ------------------------------------------------------------------------

	newAddr := RandomIPv4()

	if err := splice.AddressAdd(config.loopbackIntf, newAddr); err != nil {
		t.Fatal("AddressAdd Returned Error: ", err)
	}

	// (2)	Expect: Address is Present on Interface
	// ------------------------------------------------------------------------

	if !IntfHasAddress(t, config.loopbackIntf, newAddr) {
		t.Fatal("Loopback Does Not Have New Address")
	}
}

func TestAddressAdd_InvalidAddressValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Add Another Loopback Address
	//			Expect: Error since the address is invalid
	// ------------------------------------------------------------------------

	newAddr := &net.IPNet{
		IP: net.IP{},
	}

	if err := splice.AddressAdd(config.loopbackIntf, newAddr); err == nil {
		t.Fatal("AddressAdd Did Not Return an Error With Invalid Address")
	}
}

// ============================================================================
//	AddressDelete
// ============================================================================

func TestAddressDelete(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Delete the Loopback Address
	//			Expect: No error
	// ------------------------------------------------------------------------

	if err := splice.AddressDelete(config.loopbackIntf, IPv4LoopbackAddr); err != nil {
		t.Fatal("AddressDelete Returned Error: ", err)
	}

	// (2)	Expect: Address is NOT Present on Interface
	// ------------------------------------------------------------------------

	if IntfHasAddress(t, config.loopbackIntf, IPv4LoopbackAddr) {
		t.Fatal("Loopback Still Has Loopback Address")
	}
}

func TestAddressDelete_InvalidAddressValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Delete the Loopback Address
	//			Expect: Error since the address is invalid
	// ------------------------------------------------------------------------

	badAddr := &net.IPNet{
		IP: net.IP{},
	}

	if err := splice.AddressDelete(config.loopbackIntf, badAddr); err == nil {
		t.Fatal("AddressDelete Did Not Return an Error With Invalid Address")
	}
}
