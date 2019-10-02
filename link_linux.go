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

package splice

import (
	"github.com/vishvananda/netlink"
	"net"
)

// Provides network link manipulation for Linux using netlink.

// Administratively brings up the given network interface.
func LinkBringUp(intf *net.Interface) error {

	var err error
	var link netlink.Link

	if link, err = netlink.LinkByIndex(intf.Index); err == nil {
		return netlink.LinkSetUp(link)
	}

	return err

}

// Administratively brings down the given network interface.
func LinkBringDown(intf *net.Interface) error {

	var err error
	var link netlink.Link

	if link, err = netlink.LinkByIndex(intf.Index); err == nil {
		return netlink.LinkSetDown(link)
	}

	return err

}
