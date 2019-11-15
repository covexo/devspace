package services

import (
	"testing"

	"github.com/devspace-cloud/devspace/pkg/devspace/config/versions/latest"

	"gotest.tools/assert"
)

func TestStartPortForwarding(t *testing.T) {
	client := &client{
		config: &latest.Config{
			Dev: &latest.DevConfig{},
		},
	}
	portForwarder, err := client.StartPortForwarding()
	if err != nil {
		t.Fatalf("Error starting port forwarding with nil ports to forward: %v", err)
	}
	assert.Equal(t, true, portForwarder == nil, "Portforwarder returned despite nil port given to forward.")

	client.config = &latest.Config{
		Dev: &latest.DevConfig{
			Ports: []*latest.PortForwardingConfig{},
		},
	}
	portForwarder, err = client.StartPortForwarding()
	if err != nil {
		t.Fatalf("Error starting port forwarding with 0 ports to forward: %v", err)
	}
	assert.Equal(t, 0, len(portForwarder), "Ports forwarded despite 0 ports given to forward.")
}
