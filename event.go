package haproxyconfigparser

type EventType string

type EmitEvent func(EventType, string, Parser)

const (
	EMPTY_LINE    EventType = "empty"
	COMMENT_OUT   EventType = "comment_out"
	START_SECTION EventType = "start_section"
	IN_SECTION    EventType = "in_section"
	NORMAL        EventType = "normal"
	UNKNOWN       EventType = "unknown"
)

var (
	events         = map[string]map[string][]EmitEvent{}
	numberOfEvents = 0
)

func RegisterEvent(group, name string, emitter EmitEvent) {
	if _, ok := events[group]; !ok {
		events[group] = map[string][]EmitEvent{}
	}
	events[group][name] = append(events[group][name], emitter)
	numberOfEvents++
}

func emitEvent(group, name, line string, eventType EventType, parser Parser) {
	if numberOfEvents < 1 {
		return
	}
	names, ok := events[group]
	if ok {
		emitters, ok := names[name]
		if ok {
			for _, emit := range emitters {
				emit(eventType, line, parser)
			}
		} else if name != "*" {
			emitEvent(group, "*", line, eventType, parser)
		}
	} else if group != "*" {
		emitEvent("*", name, line, eventType, parser)
	}
}
