package main

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

func ParseProcfile(path string) (procs map[string]string, parseError error) {
	procs = make(map[string]string)
	reg := regexp.MustCompile(`(?i)(?P<name>[a-z0-9]+)\:{1}\s?(?P<command>.+)`)
	file, err := os.Open(path)

	if err != nil {
		parseError = err
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := reg.FindStringSubmatch(scanner.Text())

		if len(parts) != 3 {
			parseError = errors.New("Input file is invalid")
			return
		}

		name := parts[1]
		command := strings.Trim(parts[2], "\t ")

		procs[name] = command
	}

	return
}
