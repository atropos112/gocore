package errs

import (
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"

	"github.com/atropos112/atrogolib/pubsub"
	"github.com/atropos112/atrogolib/types"
	"github.com/negrel/assert"
)

func Alert(publisher pubsub.Publisher, source, msg string, args ...any) error {
	argsMap := make(map[string]any)
	for i := 0; i < len(args); i += 2 {
		argAsStr, ok := args[i].(string)

		if !ok {
			return errors.New("can't deserialize argument to string")
		}

		argsMap[argAsStr] = args[i+1]
	}

	event, err := types.MakeSimpleEvent(
		source,
		"error",
		"",
		map[string]any{
			"Message": msg,
			"Args":    argsMap,
		},
	)
	if err != nil {
		return err
	}

	return publisher.Publish("errors", *event, nil)
}

// AlertWrapError does the best attempt to capture as mouch information as possible, make a resonable source and alert.
// If it successfuly alerts it will return the error that was passed in (hence the "Wrap" in the name), so it can continue
// to propagate through the program.
func AlertWrapError(publisher pubsub.Publisher, err error) error {
	hostname, errHostname := os.Hostname()
	if errHostname != nil {
		hostname = "unknown"
	}

	iPAddresses, errIPAddresses := net.LookupIP(hostname)
	if errIPAddresses != nil {
		iPAddresses = []net.IP{}
	}

	serviceName := os.Getenv("ATRO_SERVICE_NAME")
	assert.NotEmpty(serviceName, "ATRO_SERVICE_NAME is empty.")

	var source string
	if serviceName != "" {
		source = serviceName
	} else {
		// This is a sad to end up. But its better to produce a not so good source than not to alert at all.
		source = "go_app_" + hostname
	}

	event, er := types.MakeSimpleEvent(
		source,
		"error",
		"",
		map[string]any{
			"error_message": err.Error(),
			"error_type":    fmt.Sprintf("%T", err),
			"hostname":      hostname,
			"ip_addresses":  iPAddresses,
			"go_version":    runtime.Version(),
			"go_os":         runtime.GOOS,
			"go_arch":       runtime.GOARCH,
			"goroutine_id":  runtime.NumGoroutine(),
		},
	)

	if er != nil {
		return er
	}

	er = publisher.Publish("errors", *event, nil)

	if er != nil {
		return er
	}

	return err
}

// AlertErrAndDie will take the error provided and if its not nil it will alert and panic.
func AlertErrAndDie(publisher pubsub.Publisher, err error) {
	assert.NotEmpty(err)

	panic(AlertWrapError(publisher, err))
}

func AlertAndDie(publisher pubsub.Publisher, msg string, args ...any) {
	serviceName := os.Getenv("ATRO_SERVICE_NAME")
	assert.NotEmpty(serviceName, "ATRO_SERVICE_NAME is empty.")

	assert.NotEmpty(msg)
	err := Alert(publisher, serviceName, msg, args...)
	if err != nil {
		panic(err)
	}

	panic(msg)
}
