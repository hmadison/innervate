package main

import (
	"net/url"
	"strings"
	"net"
)

type ApplicationtSet interface {
	ScanForApplications(dir string, startingPort *int, scannedSet ApplicationtSet) (ApplicationtSet, error)
	PortFor(host string, url *url.URL, scannedSet ApplicationtSet) (int, error)
	HasAppWithDomain(domain string) (bool)
	Applications() (map[string]Application)
}

func HostWithoutPortOrTld(input string) (string, error) {
	host := input

	if strings.Contains(input, ":") {
		splitHost, _, err := net.SplitHostPort(input)
		host = splitHost

		if err != nil {
			return "", err
		}
	}

	return strings.TrimSuffix(host, "."+tld), nil
}
