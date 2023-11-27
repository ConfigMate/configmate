package types

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type tPort struct {
	value int
}

func portFactory(value interface{}) (IType, error) {
	if value, ok := value.(int); ok {
		if !isValidPort(value) {
			return nil, fmt.Errorf("value is not a valid port")
		}

		return &tPort{value: value}, nil
	}

	return nil, fmt.Errorf("value is not a port (int)")
}

func (t tPort) TypeName() string {
	return "port"
}

func (t tPort) Value() interface{} {
	return t.value
}

func (t tPort) Methods() []string {
	return []string{
		"open",
		"live",
		"toInt",
	}
}

func (t tPort) MethodDescription(method string) string {
	tPortMethodsDescriptions := map[string]string{
		"open":  "port.open() bool : Checks that the port is open",
		"live":  "port.live() bool : Checks that the port is live",
		"toInt": "port.toInt() int : Converts the value to an int",
	}

	return tPortMethodsDescriptions[method]
}

func (t tPort) GetMethod(method string) Method {
	tPortMethods := map[string]Method{
		"open": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("port.open expects 0 arguments")
			}

			// Check if the port is open locally
			if !t.isOpen() {
				return &tBool{value: false}, fmt.Errorf("port is in use")
			}

			return &tBool{value: true}, nil
		},
		"live": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("port.live expects 0 arguments")
			}

			// Check if the port is open locally
			if !t.isLive() {
				return &tBool{value: false}, fmt.Errorf("port is not live")
			}

			return &tBool{value: true}, nil
		},
		"toInt": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("port.toInt expects 0 arguments")
			}

			// Convert to string
			return &tInt{value: t.value}, nil
		},
	}

	// Check if method doesn't exist
	if _, ok := tPortMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("host does not have method %s", method)
		}
	}

	return tPortMethods[method]
}

func isValidPort(port int) bool {
	return port > 0 && port <= 65535
}

func (t tPort) isOpen() bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(int(t.value)))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

func (t tPort) isLive() bool {
	timeout := 10 * time.Second
	conn, err := net.DialTimeout("tcp", ":"+strconv.Itoa(int(t.value)), timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
