package bus_test

import (
	"testing"

	"github.com/fredrikaugust/otelly/bus"
)

// TestNewTransportBus is just a sanity test so the github action
// will have something to munch on
func TestNewTransportBus(t *testing.T) {
	bus := bus.NewTransportBus()

	if bus.TraceBus == nil {
		t.Fatal("traceBus channel is nil")
	}
}
