package postgresql

type Config struct {
	URI string `env:"URI,notEmpty"`
}
