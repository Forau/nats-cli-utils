package natssrv_test

import (
	"github.com/Forau/nats-cli-utils/services"
	_ "github.com/Forau/nats-cli-utils/services/natssrv"
	"testing"
)

func TestServicedCreated(t *testing.T) {
	s, err := services.FindService("nats-srv")
	if err != nil {
		t.Fatal("Got error ", err)
	}
	if s == nil {
		t.Error("Service was nil")
	}
}
