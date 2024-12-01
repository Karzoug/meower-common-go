package mongo

type Config struct {
	URI string `env:"URI,notEmpty"`
}
