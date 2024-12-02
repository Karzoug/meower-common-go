package otlp

import (
	"strings"

	"go.opentelemetry.io/otel/sdk/trace"
)

type endpointExcluder struct {
	httpEndpoints map[string]struct{}
	grpcMethods   map[string]string // rpc.service -> rpc.method
	probability   float64
}

func newEndpointExcluder(httpEndpoints map[string]struct{}, grpcMethods map[string]string, probability float64) endpointExcluder {
	return endpointExcluder{
		httpEndpoints: httpEndpoints,
		grpcMethods:   grpcMethods,
		probability:   probability,
	}
}

// ShouldSample implements the sampler interface. It prevents the specified
// endpoints from being added to the trace.
func (ee endpointExcluder) ShouldSample(parameters trace.SamplingParameters) trace.SamplingResult {
	if len(ee.httpEndpoints) != 0 {
		for i := range parameters.Attributes {
			if parameters.Attributes[i].Key == "http.target" {
				if _, exists := ee.httpEndpoints[parameters.Attributes[i].Value.AsString()]; exists {
					return trace.SamplingResult{Decision: trace.Drop}
				}
			}
		}
	}

	if len(ee.grpcMethods) != 0 {
		for i := range parameters.Attributes {
			if parameters.Attributes[i].Key == "rpc.service" {
				if method, exists := ee.grpcMethods[parameters.Attributes[i].Value.AsString()]; exists {
					for j := range parameters.Attributes {
						if parameters.Attributes[j].Key == "rpc.method" {
							if parameters.Attributes[j].Value.AsString() == method {
								return trace.SamplingResult{Decision: trace.Drop}
							}
						}
					}
				}
			}
			// drop all requests to get server reflection info
			if parameters.Attributes[i].Key == "rpc.service" &&
				strings.HasPrefix(parameters.Attributes[i].Value.AsString(), "grpc.reflection") {
				return trace.SamplingResult{Decision: trace.Drop}
			}
		}
	}

	return trace.TraceIDRatioBased(ee.probability).ShouldSample(parameters)
}

// Description implements the sampler interface.
func (endpointExcluder) Description() string {
	return "customSampler"
}
