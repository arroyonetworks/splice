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
