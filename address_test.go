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

package splice_test

import (
	"github.com/arroyonetworks/splice"
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
