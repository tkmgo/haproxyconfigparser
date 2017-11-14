package haproxyconfigparser

type EventType string

type EmitEvent func(eventType EventType, line string, parser Parser)

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

// Registers an event listener to the event store.
//
// The group is any of "global", "backend", "frontend" or "*" where the asterisk means any.
//
// The name is the name is the first item on the line (i.e "bind", "acl", "use_backend" etc)
func RegisterEvent(group string, name string, emitter EmitEvent) {
	if _, ok := events[group]; !ok {
		events[group] = map[string][]EmitEvent{}
	}
	events[group][name] = append(events[group][name], emitter)
	numberOfEvents++
}

func emitEvent(group string, name, line string, eventType EventType, parser Parser) {
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
