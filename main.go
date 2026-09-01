package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)



func main() {
	// file, err := os.Open("tests/03-least-connections/4.in")
	// if err != nil {
	// 	panic(err)
	// }
	// defer file.Close()
	// sc := bufio.NewScanner(file)
	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	var lb *LeastConnSel
	for sc.Scan() {
		if sc.Err() != nil{
			panic("error in scaning")
		}
		line := sc.Text()
		if line == "" {
			continue
		}
		args := strings.Split(line, " ")
		switch args[0] {
		case "POOL":
			lb = NewLeastConnSel(args[1:])
			fmt.Println("OK")
		case "DONE":
			lb.Done(args[1])
			fmt.Println("OK")
		case "PICK":
			ans := lb.Pick()
			fmt.Println(ans)
		case "STATUS":
			lb.Status()
		default:
			fmt.Println("wrong input:", args[0])
		}
	}
}
