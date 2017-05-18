# Splice
### A high-level and multi-os library for manipulating network interfaces, links, and routes.


[![build status](https://gitlab.com/ArroyoNetworks/splice/badges/master/build.svg)](https://gitlab.com/ArroyoNetworks/splice/commits/master)
[![coverage report](https://gitlab.com/ArroyoNetworks/splice/badges/master/coverage.svg)](https://gitlab.com/ArroyoNetworks/splice/commits/master)
[![GoDoc](https://godoc.org/gitlab.com/ArroyoNetworks/splice?status.svg)](https://godoc.org/gitlab.com/ArroyoNetworks/splice)

**This library is still considered beta, the interfaces may change until 1.0 is released.**

Splice provides a high-level library for manipulating network interfaces, network links, and routes. Splice provides a common interface
for multiple operating systems.

## Getting Started

### Add Splice to Your Workspace

    go get gitlab.com/ArroyoNetworks/splice

### Examples

#### Add an Address to an Interface

```go
package main

import (
    "fmt"
    "net"
    "gitlab.com/ArroyoNetworks/splice"
)

func main() {

    var err error

    _, address, err := net.ParseCIDR("127.1.1.0/24")
    intf, err := net.InterfaceByName("eth0")    
    err = splice.AddressAdd(intf, address)

    if err != nil {
        fmt.Println("Failed to Add Address")
        fmt.Println(err)
    }

}
```

#### Bring an Interface Down

```go
package main

import (
    "fmt"
    "net"
    "gitlab.com/ArroyoNetworks/splice"
)

func main() {

    var err error

    intf, err := net.InterfaceByName("wlan0")    
    err = splice.LinkBringDown(intf)

    if err != nil {
        fmt.Println("Failed to Bring Down Interface")
        fmt.Println(err)
    }

}
```

#### Add a Route to the Routing Table

```go
package main

import (
    "fmt"
    "net"
    "gitlab.com/ArroyoNetworks/splice"
)

func main() {

    var err error

    _, dest, err := net.ParseCIDR("172.10.0.0/24")
    intf, err := net.InterfaceByName("wlan0")    

    err = splice.RouteAddViaInterface(dest, intf)

    if err != nil {
        fmt.Println("Failed to Add Route")
        fmt.Println(err)
    }

}
```

## Supported Operating Systems

### Linux

The following are supported on Linux systems:

- IP Address Configuration
- Interface Link Manipulation
- Route Manipulation

#### Dependencies

The following are third-party dependencies used for providing Linux support:

- [github.com/vishvananda/netlink](https://github.com/vishvananda/netlink)

#### Unit Tests

> Linux unit tests are automatically ran in a temporary networking namespace in order to prevent
> accidental alteration of the system's networking configuration. Because of this, test for
> Linux require the `cap_net_admin` capability.

1. Download the Test Dependencies

        go get -t gitlab.com/ArroyoNetworks/splice

2. Run the Unit Tests

        sudo -E go test github.com/ArroyoNetworks/splice

### macOS

The following are supported on Darwin systems:

- IP Address Configuration
- Interface Link Manipulation

The following are *NOT* yet supported:

- Route Manipulation (Accepting Merge Requests)

#### Dependencies

The following are third-party dependencies used for providing Darwin support:

- [golang.org/x/sys/unix](https://godoc.org/golang.org/x/sys/unix)

#### Unit Tests

Unit tests are not yet available for Darwin.

### Windows

Not yet implemented.

## License

Copyright 2017 Arroyo Networks, LLC. All rights reserved.<br/>
This project is governed by a BSD-style license. See [LICENSE](https://gitlab.com/ArroyoNetworks/splice/raw/master/LICENSE) for the full license text.
