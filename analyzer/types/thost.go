package types

import (
	"fmt"
	"regexp"
	"strings"

	probing "github.com/prometheus-community/pro-bing"
)

var tHostMethodsDescriptions map[string]string = map[string]string{
	"reachable": "host.reachable() bool : Checks that the host is reachable",
	"addPort":   "host.addPort(p port) host_port : Adds a port to the host to form a host_port type",
	"toString":  "host.toString() string : Converts the value to a string",
}

type tHost struct {
	value string
}

func hostFactory(value interface{}) (IType, error) {
	if value, ok := value.(string); ok {
		if !isValidHostname(value) {
			return nil, fmt.Errorf("value is not a valid host name")
		}

		return &tHost{value: value}, nil
	}

	return nil, fmt.Errorf("value is not an host name (string)")
}

func (t tHost) TypeName() string {
	return "host"
}

func (t tHost) Value() interface{} {
	return t.value
}

func (t tHost) GetMethod(method string) Method {
	tHostMethods := map[string]Method{
		"reachable": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("host.reachable expects 0 arguments")
			}

			// Try to ping the host
			if err := t.ping(); err != nil {
				return &tBool{value: false}, err
			}

			return &tBool{value: true}, nil
		},
		"addPort": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 1 {
				return nil, fmt.Errorf("host.addPort expects 1 argument")
			}

			// Cast argument to port type
			var port IType
			var ok bool
			port, ok = args[0].(*tPort)
			if !ok {
				// Cast argument to int type
				i, ok := args[0].(*tInt)
				if !ok {
					return nil, fmt.Errorf("host.addPort expects a port or int argument")
				}

				// Make port from the int
				var err error
				port, err = portFactory(i.value)
				if err != nil {
					return nil, err
				}
			}

			return &tHostPort{host: t.value, port: port.(*tPort).value}, nil
		},
		"toString": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("host.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: t.value}, nil
		},
	}

	// Check if method doesn't exist
	if _, ok := tHostMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("host does not have method %s", method)
		}
	}

	return tHostMethods[method]
}

func isValidHostname(hostname string) bool {
	hostname = strings.Trim(hostname, " ")
	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	return re.MatchString(hostname)
}

func (t *tHost) ping() error {
	pinger, err := probing.NewPinger(t.value)
	if err != nil {
		return err
	}

	pinger.Count = 4
	pinger.Timeout = 10000000000         // 10 seconds
	if err := pinger.Run(); err != nil { // Blocks until finished.
		return err
	}

	stats := pinger.Statistics()
	if stats.PacketLoss == 0 {
		return nil
	} else if stats.PacketLoss > 0 && stats.PacketLoss < 100 {
		return fmt.Errorf("host %s is unstable, %f%% packet loss", t.value, stats.PacketLoss)
	}
	return fmt.Errorf("host %s is unreachable", t.value)
}
