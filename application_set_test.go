package main

import "testing"

func TestHostWithoutPortOrTld(t *testing.T) {
	result, err := HostWithoutPortOrTld("test.localhost:4567")

	if result != "test" || err != nil {
		t.Fail()
	}

	result2, err := HostWithoutPortOrTld("sub.domain.test.localhost:4567")
	
	if result2 != "sub.domain.test" || err != nil {
		t.Fail()
	}

	result3, err := HostWithoutPortOrTld("test.localhost")
	
	if result3 != "test" || err != nil {
		t.Fail()
	}
}
