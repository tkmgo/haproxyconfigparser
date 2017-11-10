package haproxyconfigparser

import (
	"log"
	"testing"
)

var services *Services

func BenchmarkParseFromFile(b *testing.B) {
	var srvs *Services
	for n := 0; n < b.N; n++ {
		services, err := ParseFromFile("testdata/haproxy.cfg")
		handleError(err)
		srvs = services
	}
	services = srvs
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
