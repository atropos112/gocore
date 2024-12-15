package pubsub

import (
	"testing"

	"github.com/atropos112/gocore/logging"
)

func TestSlogFatalWithAlerting(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	l := logging.InitSlogLogger()

	NewNATSAlertContext(l, "test").ErrorAndDie("test", "a", 1, "b", 2)
}
