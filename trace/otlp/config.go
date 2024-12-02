package otlp

type Config struct {
	ServiceName         string  `env:"-"`
	ServiceVersion      string  `env:"-"`
	Probability         float64 `env:"PROBABILITY" envDefault:"0.05"`
	ExcludedHTTPRoutes  map[string]struct{}
	ExcludedGrpcMethods map[string]string // rpc.service -> rpc.method
}
