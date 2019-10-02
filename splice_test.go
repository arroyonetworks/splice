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

// Provides unit tests for the Splice library.
package splice_test

import (
	"log"
	"math/rand"
	"net"
	"testing"
)

// This file provides platform agnostic functions and helpers which can be
// used by unit tests.

var (
	IPv4LoopbackAddr *net.IPNet
)

type TestConfig struct {
	tearDownTest func()
	loopbackIntf *net.Interface
}

// Prepares a Unit Test.
// Each platform provides its own setup.
func SetUpTest(t *testing.T) *TestConfig {
	contents := &TestConfig{
		tearDownTest: _platformSetup(t),
		loopbackIntf: _platformSetupLoopback(t),
	}

	// Put a No-op for Tear Down if Platform did not Provide One
	if contents.tearDownTest == nil {
		contents.tearDownTest = func() {}
	}

	return contents
}

// ============================================================================
//	Helper Functions
// ============================================================================

// Returns a random /24 IPv4 address.
func RandomIPv4() *net.IPNet {
	return &net.IPNet{
		IP:   net.IPv4(byte(rand.Intn(127)), byte(rand.Intn(256)), byte(rand.Intn(256)), 0),
		Mask: net.IPv4Mask(255, 255, 255, 0),
	}
}

// Determines if an interface already exists.
func IntfExists(name string) bool {
	_, err := net.InterfaceByName(name)
	if err != nil {
		return false
	}
	return true
}

func IntfIsDown(intf *net.Interface) bool {
	return intf.Flags&net.FlagUp == 0
}

func IntfIsUp(intf *net.Interface) bool {
	return intf.Flags&net.FlagUp != 0
}

// ============================================================================
//	Test Utility Functions
// ============================================================================

func SkipWithReason(t *testing.T, reason string) {
	log.Print(reason)
	t.Skip(reason)
}

// Adds a random route to the given interface.
func RandomIPv4Route(t *testing.T, intf *net.Interface) *net.IPNet {
	routeNet, err := _platformRandomIPv4Route(intf)
	if err != nil {
		t.Fatal("Failed to Add Random Route: ", err)
	}
	return routeNet
}

// Determines if an interface has the given IP address configured.
func IntfHasAddress(t *testing.T, intf *net.Interface, address *net.IPNet) bool {
	return _platformIntfHasAddress(intf, address)
}

// Determines if a route exists in the system's routing table.
func RouteExists(t *testing.T, destination *net.IPNet) bool {
	return _platformRouteExists(destination)
}

// Returns a new dummy interface in the down state.
func GetDummyDownIntf(t *testing.T) *net.Interface {
	intf, err := _platformGetDummyDownIntf()
	if err != nil {
		t.Fatal("Failed to get a Downed Dummy Interface: ", err)
	}

	if !IntfIsDown(intf) {
		t.Fatalf("Interface %s Appears to be Up when Requsted Down", intf.Name)
	}

	return intf
}

// Returns a new dummy interface in the up date.
func GetDummyUpIntf(t *testing.T) *net.Interface {
	intf, err := _platformGetDummyUpIntf()
	if err != nil {
		t.Fatal("Failed to get an Up Dummy Interface: ", err)
	}

	if !IntfIsUp(intf) {
		t.Fatalf("Interface %s Appears to be Down when Requsted Up", intf.Name)
	}

	return intf
}
