package overlay

import (
	"io"
	"net/netip"

	"github.com/slackhq/nebula/routing"
)

type Device interface {
	io.ReadWriteCloser
	Activate() error
	Networks() []netip.Prefix
	Name() string
	RoutesFor(netip.Addr) routing.Gateways
	SetMTU(mtu int) error
	SupportsMultiqueue() bool
	NewMultiQueueReader() (io.ReadWriteCloser, error)
}
