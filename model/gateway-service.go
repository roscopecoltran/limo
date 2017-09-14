package model

import (
	"errors"
	"fmt"
	// "log"
	"regexp"
	"strings"
	"time"
	"github.com/devopsfaith/krakend/encoding"
	"github.com/sirupsen/logrus"
)

const (
	// BracketsRouterPatternBuilder uses brackets as route params delimiter
	Gateway_BracketsRouterPatternBuilder = iota
	// ColonRouterPatternBuilder use a colon as route param delimiter
	Gateway_ColonRouterPatternBuilder
)

// RoutingPattern to use during route conversion. By default, use the colon router pattern
var Gateway_RoutingPattern = Gateway_ColonRouterPatternBuilder

// ServiceConfig defines the krakend service
type Gateway_ServiceConfig struct {
	// set of endpoint definitions
	Endpoints []*Gateway_EndpointConfig `mapstructure:"endpoints"`
	// defafult timeout
	Timeout time.Duration `mapstructure:"timeout"`
	// default TTL for GET
	CacheTTL time.Duration `mapstructure:"cache_ttl"`
	// default set of hosts
	Host []string `mapstructure:"host"`
	// port to bind the krakend service
	Port int `mapstructure:"port"`
	// version code of the configuration
	Version int `mapstructure:"version"`

	// run krakend in debug mode
	Debug     bool
	uriParser Gateway_URIParser
}

// EndpointConfig defines the configuration of a single endpoint to be exposed
// by the krakend service
type Gateway_EndpointConfig struct {
	// url pattern to be registered and exposed to the world
	Endpoint string `mapstructure:"endpoint"`
	// HTTP method of the endpoint (GET, POST, PUT, etc)
	Method string `mapstructure:"method"`
	// set of definitions of the backends to be linked to this endpoint
	Backend []*Gateway_Backend `mapstructure:"backend"`
	// number of concurrent calls this endpoint must send to the backends
	ConcurrentCalls int `mapstructure:"concurrent_calls"`
	// timeout of this endpoint
	Timeout time.Duration `mapstructure:"timeout"`
	// duration of the cache header
	CacheTTL time.Duration `mapstructure:"cache_ttl"`
	// list of query string params to be extracted from the URI
	QueryString []string `mapstructure:"querystring_params"`
	// Endpoint Extra configuration for customized behaviour
	ExtraConfig Gateway_ExtraConfig `mapstructure:"extra_config"`
}

// Backend defines how krakend should connect to the backend service (the API resource to consume)
// and how it should process the received response
type Gateway_Backend struct {
	// the name of the group the response should be moved to. If empty, the response is
	// not changed
	Group string `mapstructure:"group"`
	// HTTP method of the request to send to the backend
	Method string `mapstructure:"method"`
	// Set of hosts of the API
	Host []string `mapstructure:"host"`
	// False if the hostname should be sanitized
	HostSanitizationDisabled bool `mapstructure:"disable_host_sanitize"`
	// URL pattern to use to locate the resource to be consumed
	URLPattern string `mapstructure:"url_pattern"`
	// set of response fields to remove. If empty, the filter id not used
	Blacklist []string `mapstructure:"blacklist"`
	// set of response fields to allow. If empty, the filter id not used
	Whitelist []string `mapstructure:"whitelist"`
	// map of response fields to be renamed and their new names
	Mapping map[string]string `mapstructure:"mapping"`
	// the encoding format
	Encoding string `mapstructure:"encoding"`
	// the response to process is a collection
	IsCollection bool `mapstructure:"is_collection"`
	// name of the field to extract to the root. If empty, the formater will do nothing
	Target string `mapstructure:"target"`

	// list of keys to be replaced in the URLPattern
	URLKeys []string
	// number of concurrent calls this endpoint must send to the API
	ConcurrentCalls int
	// timeout of this backend
	Timeout time.Duration
	// decoder to use in order to parse the received response from the API
	Decoder encoding.Decoder
	// Backend Extra configuration for customized behaviours
	ExtraConfig Gateway_ExtraConfig `mapstructure:"extra_config"`
}

type Gateway_ExtraConfig map[string]interface{}

// ConfigGetter is a function for parsing ExtraConfig into a previously know type
type Gateway_ConfigGetter func(Gateway_ExtraConfig) interface{}

// DefaultConfigGetter is the Default implementation for ConfigGetter, it just returns the ExtraConfig map.
func Gateway_DefaultConfigGetter(extra Gateway_ExtraConfig) interface{} { return extra }

const defaultNamespace = "github.com/roscopecoltran/sniperkit-limo/model"

// ConfigGetters map than match namespaces and ConfigGetter so the components knows which type to expect returned by the
// ConfigGetter ie: if we look for the defaultNamespace in the map, we will get the DefaultConfigGetter implementation
// which will return a ExtraConfig when called
var Gateway_ConfigGetters = map[string]Gateway_ConfigGetter{defaultNamespace: Gateway_DefaultConfigGetter}

var (
	simpleURLKeysPattern = regexp.MustCompile(`\{([a-zA-Z\-_0-9]+)\}`)
	debugPattern         = "^[^/]|/__debug(/.*)?$"
	errInvalidHost       = errors.New("invalid host")
	defaultPort          = 8080
)

// Init initializes the configuration struct and its defined endpoints and backends.
// Init also sanitizes the values, applies the default ones whenever necessary and
// normalizes all the things.
func (s *Gateway_ServiceConfig) Init() error {
	s.uriParser = Gateway_NewURIParser()
	if s.Version != 1 {
		return fmt.Errorf("Unsupported version: %d\n", s.Version)
	}
	if s.Port == 0 {
		s.Port = defaultPort
	}
	s.Host = s.uriParser.Gateway_CleanHosts(s.Host)
	for i, e := range s.Endpoints {
		e.Endpoint = s.uriParser.Gateway_CleanPath(e.Endpoint)

		if err := e.validate(); err != nil {
			return err
		}

		inputParams := s.extractPlaceHoldersFromURLTemplate(e.Endpoint, endpointURLKeysPattern)
		inputSet := map[string]interface{}{}
		for ip := range inputParams {
			inputSet[inputParams[ip]] = nil
		}

		e.Endpoint = s.uriParser.Gateway_GetEndpointPath(e.Endpoint, inputParams)

		s.initEndpointDefaults(i)

		for j, b := range e.Backend {

			s.initBackendDefaults(i, j)

			b.Method = strings.ToTitle(b.Method)

			if err := s.initBackendURLMappings(i, j, inputSet); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Gateway_ServiceConfig) extractPlaceHoldersFromURLTemplate(subject string, pattern *regexp.Regexp) []string {
	matches := pattern.FindAllStringSubmatch(subject, -1)
	keys := make([]string, len(matches))
	for k, v := range matches {
		keys[k] = v[1]
	}
	return keys
}

func (s *Gateway_ServiceConfig) initEndpointDefaults(e int) {
	endpoint := s.Endpoints[e]
	if endpoint.Method == "" {
		endpoint.Method = "GET"
	} else {
		endpoint.Method = strings.ToTitle(endpoint.Method)
	}
	if s.CacheTTL != 0 && endpoint.CacheTTL == 0 {
		endpoint.CacheTTL = s.CacheTTL
	}
	if s.Timeout != 0 && endpoint.Timeout == 0 {
		endpoint.Timeout = s.Timeout
	}
	if endpoint.ConcurrentCalls == 0 {
		endpoint.ConcurrentCalls = 1
	}
}

func (s *Gateway_ServiceConfig) initBackendDefaults(e, b int) {
	endpoint := s.Endpoints[e]
	backend := endpoint.Backend[b]
	if len(backend.Host) == 0 {
		backend.Host = s.Host
	} else if !backend.HostSanitizationDisabled {
		backend.Host = s.uriParser.Gateway_CleanHosts(backend.Host)
	}
	if backend.Method == "" {
		backend.Method = endpoint.Method
	}
	backend.Timeout = endpoint.Timeout
	backend.ConcurrentCalls = endpoint.ConcurrentCalls
	switch strings.ToLower(backend.Encoding) {
	case encoding.XML:
		backend.Decoder = encoding.NewXMLDecoder(backend.IsCollection)
	case encoding.RSS:
		backend.Decoder = encoding.NewRSSDecoder()
	default:
		backend.Decoder = encoding.NewJSONDecoder(backend.IsCollection)
	}
}

func (s *Gateway_ServiceConfig) initBackendURLMappings(e, b int, inputParams map[string]interface{}) error {
	backend := s.Endpoints[e].Backend[b]

	backend.URLPattern = s.uriParser.Gateway_CleanPath(backend.URLPattern)

	outputParams := s.extractPlaceHoldersFromURLTemplate(backend.URLPattern, simpleURLKeysPattern)

	outputSet := map[string]interface{}{}
	for op := range outputParams {
		outputSet[outputParams[op]] = nil
	}

	if len(outputSet) > len(inputParams) {
		return fmt.Errorf("Too many output params! input: %v, output: %v\n", outputSet, outputParams)
	}

	tmp := backend.URLPattern
	backend.URLKeys = make([]string, len(outputParams))
	for o := range outputParams {
		if _, ok := inputParams[outputParams[o]]; !ok {
			return fmt.Errorf("Undefined output param [%s]! input: %v, output: %v\n", outputParams[o], inputParams, outputParams)
		}
		tmp = strings.Replace(tmp, "{"+outputParams[o]+"}", "{{."+strings.Title(outputParams[o])+"}}", -1)
		backend.URLKeys = append(backend.URLKeys, strings.Title(outputParams[o]))
	}
	backend.URLPattern = tmp
	return nil
}

func (e *Gateway_EndpointConfig) validate() error {
	matched, err := regexp.MatchString(debugPattern, e.Endpoint)
	if err != nil {
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "gateway-service.go", 
							"method_name": "validate", 
							"engine": "sniperkit-gateway", 
							}).Warnf("ERROR: parsing the endpoint url [%s]: %s. Ignoring\n", e.Endpoint, err.Error())
		// log.Printf("ERROR: parsing the endpoint url [%s]: %s. Ignoring\n", e.Endpoint, err.Error())
		return err
	}
	if matched {
		err := fmt.Errorf("ERROR: the endpoint url path [%s] is not a valid one!!! Ignoring\n", e.Endpoint)
		log.WithError(err).WithFields(
			logrus.Fields{	"file": "gateway-service.go", 
							"method_name": "validate", 
							"engine": "sniperkit-gateway", 
							}).Warnf("ERROR: the endpoint url path [%s] is not a valid one!!! Ignoring\n", e.Endpoint)
		return err
	}

	if len(e.Backend) == 0 {
		return fmt.Errorf("WARNING: the [%s] endpoint has 0 backends defined! Ignoring\n", e.Endpoint)
	}
	return nil
}