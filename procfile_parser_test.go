package main

import "testing"

func TestParseProcfile(t *testing.T) {
	procfile, err := ParseProcfile("./__support__/sample_procfile")

	if err != nil {
		t.Fail()
	}

	if len(procfile) != 3 {
		t.Fail()
	}

	command, ok := procfile["name"]

	if !ok || command != "command" {
		t.Fail()
	}

	command, ok = procfile["name2"]

	if !ok || command != "command" {
		t.Fail()
	}

	command, ok = procfile["name3"]

	if !ok || command != "command:with:colon" {
		t.Fail()
	}
}
