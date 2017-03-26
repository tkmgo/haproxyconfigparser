package haproxyconfigparser

func Parse(config []string) (*Services, error) {
	services := NewServices()
	var (
		parser Parser
		err    error
	)
	for _, line := range config {
		items, enable := SeparateConfigLine(line)
		if len(items) == 0 {
			emitEvent("", "", line, EMPTY_LINE, nil)
			continue
		}
		maybeSection := false
		switch items[0] {
		case "global":
			maybeApply(services, parser)
			parser = NewGlobalParser()
		case "frontend":
			maybeApply(services, parser)
			parser = NewFrontendParser(items[1])
		case "backend":
			maybeApply(services, parser)
			parser = NewBackendParser(items[1])
		case "defaults", "listen":
			maybeApply(services, parser)
			parser = NewNilParser(items[0]) // TODO has not implemented yet
		default:
			maybeSection = true
			if parser != nil {
				if err = parser.Parse(items[0], items[1:], enable); err != nil {
					return services, err
				}
				ev := IN_SECTION
				if !enable {
					ev = COMMENT_OUT
				}
				emitEvent(parser.Name(), items[0], line, ev, parser)
			} else {
				emitEvent("*", "*", line, UNKNOWN, nil) // case: unknown section
			}
		}
		if !maybeSection {
			emitEvent(items[0], "", line, START_SECTION, parser)
		}
	}
	maybeApply(services, parser)

	for _, f := range services.Frontends {
		if err := backendReferenceByAcl(f, services.Backends); err != nil {
			return services, err
		}
	}

	return services, nil
}


func maybeApply(s *Services, parser Parser) {
	if parser != nil {
		parser.Install(s)
	}
}
