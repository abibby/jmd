package run

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type StringStack []string

func (p *StringStack) Push(str string) {
	*p = append(*p, str)
}
func (p *StringStack) Pop() string {
	val := (*p)[p.Len()-1]
	(*p) = (*p)[:p.Len()-1]
	return val
}
func (p *StringStack) Len() int {
	return len(*p)
}
func (p *StringStack) Join(sep string) string {
	return strings.Join(*p, sep)
}
func ExtractFile(file, section string) (string, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	file, err = filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return Extract(file, string(src), section)
}

func Extract(file, src, section string) (string, error) {
	out := ""
	lines(src, func(i int, line, subSection string) {
		if !strings.HasPrefix(subSection, section+".") && subSection != section {
			return
		}
		out += line + "\n"
	})

	if out == "" {
		return "", fmt.Errorf("No section %s", section)
	}

	parts := strings.SplitN(out, "\n", 2)

	if len(parts) < 2 {
		return "", nil
	}

	out, err := Run(file, parts[1])
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}
func lines(src string, cb func(i int, line, section string)) {
	currentPath := StringStack{}
	re := regexp.MustCompile("(#+)[ \t]*(.*)[ \t]*$")

	for i, line := range strings.Split(src, "\n") {
		parts := re.FindStringSubmatch(line)

		if len(parts) != 0 {
			for len(parts[1]) <= currentPath.Len() {
				currentPath.Pop()
			}
			currentPath.Push(parts[2])
		}
		cb(i, line, currentPath.Join("."))

	}
}
