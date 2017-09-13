package model

import (
	"regexp"
	"strings"
)

var (
	endpointURLKeysPattern = regexp.MustCompile(`/\{([a-zA-Z\-_0-9]+)\}`)
	hostPattern            = regexp.MustCompile(`(https?://)?([a-zA-Z0-9\._\-]+)(:[0-9]{2,6})?/?`)
)

// URIParser defines the interface for all the URI manipulation required by KrakenD
type Gateway_URIParser interface {
	Gateway_CleanHosts([]string) []string
	Gateway_CleanHost(string) string
	Gateway_CleanPath(string) string
	Gateway_GetEndpointPath(string, []string) string
}

// NewURIParser creates a new URIParser using the package variable RoutingPattern
func Gateway_NewURIParser() Gateway_URIParser {
	return Gateway_URI(Gateway_RoutingPattern)
}

// URI implements the URIParser interface
type Gateway_URI int

// CleanHosts applies the CleanHost method to every member of the received array of hosts
func (u Gateway_URI) Gateway_CleanHosts(hosts []string) []string {
	cleaned := []string{}
	for i := range hosts {
		cleaned = append(cleaned, u.Gateway_CleanHost(hosts[i]))
	}
	return cleaned
}

// CleanHost sanitizes the received host
func (Gateway_URI) Gateway_CleanHost(host string) string {
	matches := hostPattern.FindAllStringSubmatch(host, -1)
	if len(matches) != 1 {
		panic(errInvalidHost)
	}
	keys := matches[0][1:]
	if keys[0] == "" {
		keys[0] = "http://"
	}
	return strings.Join(keys, "")
}

// CleanPath trims all the extra slashes from the received URI path
func (Gateway_URI) Gateway_CleanPath(path string) string {
	return "/" + strings.TrimPrefix(path, "/")
}

// GetEndpointPath applies the proper replacement in the received path to generate valid route patterns
func (u Gateway_URI) Gateway_GetEndpointPath(path string, params []string) string {
	result := path
	if u == Gateway_ColonRouterPatternBuilder {
		for p := range params {
			result = strings.Replace(result, "/{"+params[p]+"}", "/:"+params[p], -1)
		}
	}
	return result
}