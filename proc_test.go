package main

import "testing"

func TestStartRestartStop(t *testing.T) {
	proc := Proc{Dir: "./__support__", Command: "bash wait.sh"}
	proc.Start()

	if !proc.IsRunning() {
		t.Fail()
	}

	proc.Restart()

	if !proc.IsRunning() {
		t.Fail()
	}

	proc.Stop()

	if proc.IsRunning() {
		t.Fail()
	}
}
