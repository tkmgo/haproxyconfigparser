package haproxyconfigparser

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestRegisterEvent(t *testing.T) {
	backend := 0
	RegisterEvent("backend", "*", func(eventType EventType, line string, parser Parser) {
		backend++
	})

	frontend := 0
	RegisterEvent("frontend", "*", func(eventType EventType, line string, parser Parser) {
		frontend++
	})

	global := 0
	RegisterEvent("global", "*", func(eventType EventType, line string, parser Parser) {
		global++
	})

	_, err := ParseFromFile("testdata/haproxy.cfg")
	if err != nil {
		t.Errorf("Failed to parse data: %s", err)
	}
	assert.Equal(t, 6, global, "Expected different number of events for global.")
	assert.Equal(t, 16, backend, "Expected different number of events for backend.")
	assert.Equal(t, 21, frontend, "Expected different number of events for frontend.")
}

func ExampleRegisterEvent() {
	RegisterEvent("backend", "*", func(eventType EventType, line string, parser Parser) {
		if eventType == START_SECTION {
			fmt.Printf("EventType=%s, line=%s\n", eventType, line)
		}
	})

	_, err := ParseFromFile("testdata/haproxy.cfg")
	if err != nil {
		fmt.Errorf("Failed to parse data: %s", err)
	}
	//Output:
	//EventType=start_section, line=backend profileEditingService_20000
	//EventType=start_section, line=backend accountCreationService_10000
}
