package haproxyconfigparser

// based on http://www.haproxy.org/download/1.4/doc/configuration.txt

import (
	"fmt"
	"strconv"
)

func Parse(config []string) (*Services, error) {
	services := NewServices()
	var (
		parser Parser
		err    error
	)
	for _, line := range config {
		items, enable := SeparateConfigLine(line)
		if len(items) == 0 {
			continue
		}
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
			parser = nil
		default:
			if parser != nil {
				if err = parser.Parse(items[0], items[1:], enable); err != nil {
					return services, err
				}
			}
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

type Parser interface {
	Parse(node string, options []string, enable bool) error
	Install(s *Services)
}

/*
The following keywords are supported in the "global" section :

 * Process management and security
   - chroot
   - daemon
   - gid
   - group
   - log
   - log-send-hostname
   - nbproc
   - pidfile
   - uid
   - ulimit-n
   - user
   - stats
   - node
   - description

 * Performance tuning
   - maxconn
   - maxpipes
   - noepoll
   - nokqueue
   - nopoll
   - nosepoll
   - nosplice
   - spread-checks
   - tune.bufsize
   - tune.chksize
   - tune.maxaccept
   - tune.maxpollevents
   - tune.maxrewrite
   - tune.rcvbuf.client
   - tune.rcvbuf.server
   - tune.sndbuf.client
   - tune.sndbuf.server

 * Debugging
   - debug
   - quiet
*/
type GlobalParser struct {
	Global Global
}

func (self *GlobalParser) Parse(node string, options []string, enable bool) error {
	if !enable {
		return nil
	}
	switch node {
	case "stats":
		size := len(options)
		for i, v := range options {
			if v == "socket" {
				if size > i+1 {
					addr, typ := ParseSockAddr(options[i+1])
					self.Global.Stats = append(self.Global.Stats, &Socket{
						Type: typ,
						Addr: addr,
					})
				} else {
					return fmt.Errorf("Can not get socket path from options.")
				}
			}
		}
	case "daemon":
		self.Global.Daemon = true
	case "user":
		self.Global.User = options[0]
	case "group":
		self.Global.Group = options[0]
	case "maxconn":
		maxconn, err := strconv.Atoi(options[0])
		if err != nil {
			return err
		}
		self.Global.Maxconn = maxconn
	}
	return nil
}

func (self *GlobalParser) Install(s *Services) {
	s.Global = self.Global
}

func NewGlobalParser() *GlobalParser {
	return &GlobalParser{}
}

type FrontendParser struct {
	Frontend Frontend
}

func (self *FrontendParser) Parse(node string, options []string, enable bool) error {
	if !enable {
		return nil
	}
	// TODO support default_backend
	switch node {
	case "bind": // TODO '0.0.0.0:443 ssl crt /path/to/server.pem' style
		host, port, err := SeparateHostAndPort(options[0])
		if err != nil {
			return err
		}
		self.Frontend.Bind.Host = host
		self.Frontend.Bind.Port = port
	case "mode":
		self.Frontend.Mode = options[0]
	case "maxconn":
		maxconn, err := strconv.Atoi(options[0])
		if err != nil {
			return err
		}
		self.Frontend.Maxconn = maxconn
	case "acl":
		acl := &Acl{
			Name:      options[0],
			Type:      options[1],
			Condition: options[2:],
		}
		self.Frontend.Acls = append(self.Frontend.Acls, acl)
	case "use_backend":
		if len(options) < 3 {
			return fmt.Errorf("[ACL] No conditions are defined in use_backend '%s'", options)
		}
		bq, err := CreateUseBackendClauses(options[1], options[2:])
		if err != nil {
			return err
		}
		self.Frontend.UseBackends = append(self.Frontend.UseBackends, &UseBackend{
			Name:      options[0],
			Condition: bq,
		})
	}
	return nil
}

func (self *FrontendParser) Install(s *Services) {
	s.Frontends = append(s.Frontends, self.Frontend)
}

func NewFrontendParser(name string) *FrontendParser {
	return &FrontendParser{
		Frontend: Frontend{
			Name: name,
		},
	}
}

type BackendParser struct {
	Backend Backend
}

func (self *BackendParser) Parse(node string, options []string, enable bool) error {
	switch node {
	case "option":
		if enable {
			self.Backend.Options = append(self.Backend.Options, options)
		}
	case "server":
		host, port, err := SeparateHostAndPort(options[1])
		if err != nil {
			return err
		}
		s := &Server{
			Label:   options[0],
			Host:    host,
			Port:    port,
			Enabled: enable,
		}
		if len(options) >= 3 {
			s.Options = options[2:]
		}
		self.Backend.Servers = append(self.Backend.Servers, s)
	case "mode":
		if enable {
			self.Backend.Mode = options[0]
		}
	case "balance":
		if enable {
			self.Backend.Balance = options[0]
		}
	case "http-request":
		if enable {
			self.Backend.HttpRequest = options
		}
	}
	return nil
}

func (self *BackendParser) Install(s *Services) {
	s.Backends = append(s.Backends, self.Backend)
}

func NewBackendParser(name string) *BackendParser {
	return &BackendParser{
		Backend: Backend{
			Name: name,
		},
	}
}

func maybeApply(s *Services, parser Parser) {
	if parser != nil {
		parser.Install(s)
	}
}
