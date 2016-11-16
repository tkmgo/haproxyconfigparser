package haproxyconfigparser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ParseFromStdin() (*Services, error) {
	var fp *os.File
	var err error
	fp = os.Stdin
	reader := bufio.NewReaderSize(fp, 4096)
	dest := make([]string, 0)
	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		dest = append(dest, line)
	}
	if len(dest) == 0 {
		return NewServices(), fmt.Errorf("No config lines")
	}
	return Parse(dest)
}

func ParseFromFile(path string) (*Services, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return NewServices(), err
	}
	dest := strings.Split(string(content), "\n")
	if len(dest) == 0 {
		return NewServices(), fmt.Errorf("No config lines")
	}
	return Parse(dest)
}
