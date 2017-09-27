package haproxyconfigparser

import (
	"fmt"
	"strings"
)

func CreateUseBackendClauses(judge string, source []string) (*UseBackendClauses, error) {
	// https://www.haproxy.com/doc/aloha/7.0/haproxy/conditions.html
	dest := &UseBackendClauses{}

	if judge == "if" {
		dest.ReverseJudge = false
	} else if judge == "unless" {
		dest.ReverseJudge = true
	} else {
		return dest, fmt.Errorf("expected if|unless, but '%s'", judge)
	}

	// TODO suport [!], such as !if_xxx
	buf := make([]string, 0)

	for _, n := range source {
		if n == "or" || n == "OR" || n == "||" {
			if len(buf) > 0 {
				dest.Any = append(dest.Any, buf)
			}
		} else if n == "and" || n == "AND" {
			continue
		} else {
			buf = append(buf, n)
		}
	}
	if len(buf) > 0 {
		dest.Any = append(dest.Any, buf)
	}

	return dest, nil
}

func backendReferenceByAcl(frontend Frontend, backends []Backend) error {
	for _, ub := range frontend.UseBackends {
		b, err := findBackendByName(ub.Name, backends)
		if err != nil {
			return err
		}
		ub.Backend = b

		for _, a := range ub.Condition.Any {
			// TODO support or/and conditions
			for _, s := range a {
				//TODO Handle ! correctly instead of ignoring it
				if strings.HasPrefix(s, "!") {
					s=strings.TrimLeft(s,"!")
				}
				acl, err := findAclByName(s, &frontend)
				if err != nil {
					return err
				}
				ub.Acls = append(ub.Acls, acl)
			}
		}
	}
	return nil
}

func findAclByName(name string, frontend *Frontend) (*Acl, error) {
	for _, acl := range frontend.Acls {
		if acl.Name == name {
			return acl, nil
		}
	}
	return &Acl{}, fmt.Errorf("ACL '%s' not found", name)
}

func findBackendByName(name string, backends []Backend) (*Backend, error) {
	for _, b := range backends {
		if b.Name == name {
			return &b, nil
		}
	}
	return nil, fmt.Errorf("Backend '%s' not found", name)
}
