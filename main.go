package main

import (
	"bufio"
	"fmt"
	"os"
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


func main() {
	// var test int
	// var err error
	// if len(os.Args) < 2 {
	// 	test = 2
	// } else {
	// 	test, err = strconv.Atoi(os.Args[1])
	// 	if err != nil {
	// 		panic("wrong arguman")
	// 	}
	// }

	// file, err := os.Open(fmt.Sprintf("tests/13-tls-termination/%d.in", test))
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var lb = NewCertManager()

	for sc.Scan() {
		if sc.Err() != nil {
			panic("error in scaning")
		}
		line := sc.Text()
		if line == "" {
			continue
		}
		args := strings.Split(line, " ")
		switch args[0] {
		case "CERT":
			lb.Cert(args[1], args[2])
			fmt.Println("OK")
		case "LOOKUP":
			ans := lb.LookUp(args[1])
			fmt.Println(ans)
		case "LIST":
			ans := lb.List()
			fmt.Println(ans)
		default:
			fmt.Println("Wrong argument:", args[0])
		}
	}
}
