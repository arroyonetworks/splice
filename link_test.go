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
//	LinkBringUp
// ============================================================================

func TestLinkBringUp(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Get a Downed Link
	// ------------------------------------------------------------------------

	intf := GetDummyDownIntf(t)

	// (2)	Bring the Link Up
	//			Expect: The link is actually up.
	// ------------------------------------------------------------------------

	err := splice.LinkBringUp(intf)
	if err != nil {
		t.Fatal("Failed to Bring Link Up: ", err)
	}

	intf, err = net.InterfaceByIndex(intf.Index)

	if err != nil || !IntfIsUp(intf) {
		t.Fatal("Downed Interface did not Come Up")
	}

}

func TestLinkBringUp_InvalidIntfValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Bring the Link Up
	//			Expect: Error since the interface is invalid
	// ------------------------------------------------------------------------

	intf := &net.Interface{Index: -1}
	if err := splice.LinkBringUp(intf); err == nil {
		t.Fatal("LinkBringUp Did Not Return an Error with Invalid Interface value")
	}

}

// ============================================================================
//	LinkBringDown
// ============================================================================

func TestLinkBringDown(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Get an Upped Link
	// ------------------------------------------------------------------------

	intf := GetDummyUpIntf(t)

	// (2)	Bring the Link Down
	//			Expect: The link is actually down.
	// ------------------------------------------------------------------------

	err := splice.LinkBringDown(intf)
	if err != nil {
		t.Fatal("Failed to Bring Link Down: ", err)
	}

	intf, err = net.InterfaceByIndex(intf.Index)

	if err != nil || !IntfIsDown(intf) {
		t.Fatal("Upped Interface did not Come Down")
	}
}

func TestLinkBringDown_InvalidIntfValue(t *testing.T) {

	config := SetUpTest(t)
	defer config.tearDownTest()

	// (1)	Bring the Link Down
	//			Expect: Error since the interface is invalid
	// ------------------------------------------------------------------------

	intf := &net.Interface{Index: -1}
	if err := splice.LinkBringDown(intf); err == nil {
		t.Fatal("LinkBringDown Did Not Return an Error with Invalid Interface value")
	}

}
