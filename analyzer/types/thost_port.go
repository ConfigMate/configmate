package types

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

type tHostPort struct {
	host string
	port int
}

func hostPortFactory(value interface{}) (IType, error) {
	if value, ok := value.(string); ok {
		host, port, err := isValidHostPort(value)
		if err != nil {
			return nil, err
		}

		return &tHostPort{host: host, port: port}, nil
	}

	return nil, fmt.Errorf("value is not an host:port string")
}

func (t tHostPort) TypeName() string {
	return "host_port"
}

func (t tHostPort) Value() interface{} {
	return t.host + ":" + strconv.Itoa(t.port)
}

func (t tHostPort) Methods() []string {
	return []string{
		"live",
		"getHost",
		"getPort",
		"toString",
	}
}

func (t tHostPort) MethodDescription(method string) string {
	tHostPortMethodsDescriptions := map[string]string{
		"live":     "host_port.live() bool : Checks that the host:port is live",
		"getHost":  "host_port.getHost() host : Gets the host",
		"getPort":  "host_port.getPort() port : Gets the port",
		"toString": "host_port.toString() string : Converts the value to a string",
	}

	return tHostPortMethodsDescriptions[method]
}

func (t tHostPort) GetMethod(method string) Method {
	tHostPortMethods := map[string]Method{
		"live": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("host_port.live expects 0 arguments")
			}

			// Check if the port is open locally
			if !t.isLive() {
				return &tBool{value: false}, fmt.Errorf("host:port is not live")
			}

			return &tBool{value: true}, nil
		},
		"getHost": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("host_port.getHost expects 0 arguments")
			}

			return &tHost{value: t.host}, nil
		},
		"getPort": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("host_port.getPort expects 0 arguments")
			}

			return &tPort{value: t.port}, nil
		},
		"toString": func(args []IType) (IType, error) {
			// Check that the correct number of arguments were passed
			if len(args) != 0 {
				return nil, fmt.Errorf("host_port.toString expects 0 arguments")
			}

			// Convert to string
			return &tString{value: t.host + ":" + strconv.Itoa(t.port)}, nil
		},
	}

	// Check if method doesn't exist
	if _, ok := tHostPortMethods[method]; !ok {
		return func(args []IType) (IType, error) {
			return nil, fmt.Errorf("host_port does not have method %s", method)
		}
	}

	return tHostPortMethods[method]
}

func isValidHostPort(hostPort string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", 0, err // Not a valid host:port combination
	}

	// Convert port to integer
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, err // Not a valid port
	}

	// Additional checks can be added here if needed
	return host, port, nil
}

func (t tHostPort) isLive() bool {
	timeout := 10 * time.Second
	address := net.JoinHostPort(t.host, strconv.Itoa(t.port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
