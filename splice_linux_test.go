package splice_test

import (
	"fmt"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"math/rand"
	"net"
	"os"
	"testing"
)

// ============================================================================
//	Test Scaffolding for Linux
// ============================================================================

func init() {
	IPv4LoopbackAddr = &net.IPNet{
		IP:   net.ParseIP("127.0.0.1"),
		Mask: net.IPv4Mask(255, 0, 0, 0),
	}
}

// Sets up a new Linux Test.
// In Linux, we are able to create a new Network Namespace to run tests
// in a fresh network stack, and prevent any interruptions to the host.
func _platformSetup(t *testing.T) func() {

	if os.Getuid() != 0 {
		SkipWithReason(t, "Test Setup Failed: Root Privileges are Required")
	}

	ns, err := netns.New()
	if err != nil {
		SkipWithReason(t, "Test Setup Failed: Failed to Created Network Namespace: "+err.Error())
	}

	return func() {
		ns.Close()
	}
}

// Sets up the Linux Loopback Adapter.
func _platformSetupLoopback(t *testing.T) *net.Interface {

	intf, err := net.InterfaceByName("lo")
	if err != nil {
		SkipWithReason(t, "Test Setup Failed: Failed to Prepare Loopback: "+err.Error())
	}

	loopback, err := netlink.LinkByIndex(intf.Index)
	if err != nil {
		SkipWithReason(t, "Test Setup Failed: Failed to Prepare Loopback: "+err.Error())
	}

	if err := netlink.LinkSetUp(loopback); err != nil {
		SkipWithReason(t, "Test Setup Failed: Failed to Prepare Loopback: "+err.Error())
	}

	return intf
}

func _platformRandomIPv4Route(intf *net.Interface) (*net.IPNet, error) {

	link, err := netlink.LinkByIndex(intf.Index)
	if err != nil {
		return nil, err
	}

	destination := RandomIPv4()

	err = netlink.RouteAdd(&netlink.Route{
		Dst:       destination,
		LinkIndex: link.Attrs().Index,
	})
	if err != nil {
		return nil, err
	}

	return destination, nil
}

func _platformIntfHasAddress(intf *net.Interface, address *net.IPNet) bool {

	link, err := netlink.LinkByIndex(intf.Index)
	if err != nil {
		return false
	}

	addrs, err := netlink.AddrList(link, netlink.FAMILY_ALL)
	if err != nil {
		return false
	}

	for _, addr := range addrs {
		if addr.IPNet.String() == address.String() {
			return true
		}
	}
	return false
}

func _platformRouteExists(destination *net.IPNet) bool {

	filter := &netlink.Route{
		Dst: destination,
	}

	routes, err := netlink.RouteListFiltered(netlink.FAMILY_ALL, filter, 0)
	if err != nil {
		return false
	}

	return len(routes) != 0
}

func _platformGetDummyUpIntf() (*net.Interface, error) {

	// Find a free interface name
	intfName := ""
	for true {
		intfName = fmt.Sprintf("dummy%d", rand.Intn(127))
		if !IntfExists(intfName) {
			break
		}
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = intfName

	link := &netlink.Dummy{LinkAttrs: attrs}

	err := netlink.LinkAdd(link)
	if err != nil {
		return nil, err
	}

	err = netlink.LinkSetUp(link)
	if err != nil {
		return nil, err
	}

	return net.InterfaceByName(attrs.Name)

}

func _platformGetDummyDownIntf() (*net.Interface, error) {

	// Find a free interface name
	intfName := ""
	for true {
		intfName = fmt.Sprintf("dummy%d", rand.Intn(127))
		if !IntfExists(intfName) {
			break
		}
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = intfName
	attrs.OperState = netlink.OperDown

	link := &netlink.Dummy{LinkAttrs: attrs}

	err := netlink.LinkAdd(link)
	if err != nil {
		return nil, err
	}

	return net.InterfaceByName(attrs.Name)

}
