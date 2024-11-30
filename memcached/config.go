package memcached

type Config struct {
	Addresses []string `env:"ADDRESSES,notEmpty" envSeparator:","`
}
