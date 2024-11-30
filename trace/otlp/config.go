package otlp

type Config struct {
	ServiceName    string              `env:"-"`
	ServiceVersion string              `env:"-"`
	ExcludedRoutes map[string]struct{} `env:"-"`
	Probability    float64             `env:"PROBABILITY" envDefault:"0.05"`
}
