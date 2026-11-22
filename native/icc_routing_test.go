package main

import (
	"net"
	"strings"
	"testing"
	"time"
)

// The synthetic destination is deliberately reachable on the host. A tunnel
// failure must not silently send an AllowedIPs destination through that host.
func TestAllowedDestinationNeverFallsBackToDirect(t *testing.T) {
	id := startEngine("private_key="+strings.Repeat("1", 64)+"\n", "192.0.2.99/32", "", 1280, -1)
	if id == 0 {
		t.Fatal("engine startup failed")
	}
	defer stopEngine(id)
	destination, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer destination.Close()
	reserve, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := reserve.Addr().(*net.TCPAddr).Port
	reserve.Close()
	if msg := startProxy(id, "127.0.0.1", port, "test", "secret", "127.0.0.1/32"); msg != "" {
		t.Fatal(msg)
	}
	c := authenticate(t, getEngine(id).socks, "secret")
	defer c.Close()
	c.SetDeadline(time.Now().Add(20 * time.Second))
	expect(t, c, []byte{1, 0})
	request(t, c, 1, addressBytes(destination.Addr()))
	expect(t, c, []byte{5, 5, 0})
	destination.(*net.TCPListener).SetDeadline(time.Now().Add(100 * time.Millisecond))
	leaked, err := destination.Accept()
	if err == nil {
		leaked.Close()
		t.Fatal("AllowedIPs traffic leaked to direct network")
	}
	if ne, ok := err.(net.Error); !ok || !ne.Timeout() {
		t.Fatalf("unexpected listener error: %v", err)
	}
}

func TestProxyRejectsPublicBindAndMissingRoutes(t *testing.T) {
	id := startEngine("private_key="+strings.Repeat("1", 64)+"\n", "192.0.2.99/32", "", 1280, -1)
	if id == 0 {
		t.Fatal("engine startup failed")
	}
	defer stopEngine(id)
	for _, c := range []struct{ host, routes string }{{"0.0.0.0", "192.0.2.0/24"}, {"::", "192.0.2.0/24"}, {"127.0.0.1", ""}, {"127.0.0.1", "invalid"}} {
		if msg := startProxy(id, c.host, 12345, "test", "secret", c.routes); msg == "" {
			t.Fatalf("accepted invalid proxy settings: %+v", c)
		}
	}
}
