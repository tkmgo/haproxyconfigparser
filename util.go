package haproxyconfigparser

import (
	"fmt"
	"strconv"
	"strings"
)

func SeparateHostAndPort(address string) (string, int, error) {
	hostAndPort := strings.Split(address, ":")
	if len(hostAndPort) != 2 {
		return "", 0, fmt.Errorf("'%s' is not in host:port style", address)
	}
	port, err := strconv.Atoi(hostAndPort[1])
	if err != nil {
		return "", 0, err
	}
	return hostAndPort[0], port, nil
}

func Uncomment(line string) (string, bool) {
	buf := make([]rune, 0)
	enabled := true
	for _, c := range strings.TrimSpace(line) {
		if c == '#' {
			if len(buf) == 0 { // case: '#this is line'
				enabled = false
			} else {
				break
			}
		} else {
			buf = append(buf, c)
		}
	}
	return strings.TrimSpace(string(buf)), enabled
}

func SeparateConfigLine(line string) ([]string, bool) {
	items := make([]string, 0)
	uncomment, enable := Uncomment(strings.Replace(line, "\t", " ", -1))
	for _, n := range strings.Split(uncomment, " ") {
		if nn := strings.TrimSpace(n); nn != "" {
			items = append(items, nn)
		}
	}
	return items, enable
}

func ParseSockAddr(path string) (string, string) {
	addrs := strings.Split(path, "@") // separate ipv4@0.0.0.0:9999
	addr := addrs[len(addrs)-1]
	if strings.HasPrefix(addr, "/") {
		return addr, "unix"
	} else {
		return addr, "tcp"
	}
}
