package load_balancer

import (
	"sort"
	"strings"
)

type CertManager struct {
	certs map[string]string
}

func NewCertManager() *CertManager {
	return &CertManager{certs: make(map[string]string)}
}

func (c *CertManager) Cert(pattern, name string) {
	c.certs[pattern] = name
}

func (c *CertManager) LookUp(host string) string {
	if name, ok := c.certs[host]; ok {
		return name
	}

	if i := strings.IndexByte(host, '.'); i != -1 {
		if name, ok := c.certs["*."+host[i+1:]]; ok {
			return name
		}
	}

	return "default"
}

func (c *CertManager) List() string {
	keys := make([]string, 0, len(c.certs))
	for k := range c.certs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	lines := make([]string, len(keys))
	for i, k := range keys {
		lines[i] = k + " -> " + c.certs[k]
	}

	return strings.Join(lines, "\n")
}
