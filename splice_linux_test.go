package splice_test

import (
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"log"
	"net"
	"os"
	"testing"
)

// ============================================================================
//	Test Setup and Teardown Helpers
// ============================================================================

type namespaceTestTeardown func()

func skipUnlessRoot(t *testing.T) {

	if os.Getuid() != 0 {
		msg := "Skipped test: root privileges are required."
		log.Printf(msg)
		t.Skip(msg)
	}
}

func setupNamespaceTest(t *testing.T) namespaceTestTeardown {

	skipUnlessRoot(t)

	ns, err := netns.New()
	if err != nil {
		t.Fatal("Failed to Create Network Namespace", ns)
	}

	return func() {
		ns.Close()
	}
}

// ============================================================================
//	Test Utility Functions
// ============================================================================

func setLoopbackUp() error {
	loopback, err := netlink.LinkByName("lo")
	if err != nil {
		return err
	}

	netlink.LinkSetUp(loopback)

	return nil
}

func addRouteToLoopback(network string) error {

	loopback, err := netlink.LinkByName("lo")
	if err != nil {
		return err
	}

	if err := setLoopbackUp(); err != nil {
		return err
	}

	_, destination, err := net.ParseCIDR(network)
	if err != nil {
		return err
	}

	err = netlink.RouteAdd(&netlink.Route{
		Dst:       destination,
		LinkIndex: loopback.Attrs().Index,
	})
	if err != nil {
		return err
	}

	return nil
}
