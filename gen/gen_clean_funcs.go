//go:generate go run gen_clean_funcs.go

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.OpenFile("../clean_funcs.gen.go", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString("// GENERATED FILE. DO NOT EDIT.\n\n")
	f.WriteString("package panik\n\n")
	f.WriteString("var cleanFuncs []string = []string {\n")
	f.WriteString(fmt.Sprintf("\t\"%s\",\n", "panic"))
	f.WriteString(fmt.Sprintf("\t\"%s\",\n", "runtime/debug.Stack"))
	for _, funcName := range getNames() {
		f.WriteString(fmt.Sprintf("\t\"github.com/setlog/panik.%s\",\n", funcName))
	}
	f.WriteString("}\n")
}

func getNames() (lines []string) {
	f, err := os.Open("../panik.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "func") {
			funcName := strings.TrimSpace(line[4:strings.Index(line, "(")])
			if funcName != "" {
				lines = append(lines, funcName)
			}
		}
	}
	return lines
}
