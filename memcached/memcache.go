package memcached

import (
	"context"
	"fmt"
	"net"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/rs/zerolog"
)

type Client struct {
	*memcache.Client
}

// NewStaticClient creates memcached client with static list of addresses.
func NewStaticClient(ctx context.Context, cfg Config, logger zerolog.Logger) (Client, error) {
	naddr := make([]net.Addr, 0)
	for _, addr := range cfg.Addresses {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return Client{}, fmt.Errorf("invalid address %q: %w", addr, err)
		}
		ips, err := net.LookupIP(host)
		if err != nil {
			return Client{}, fmt.Errorf("unexpected error while looking up memcached IPs: %w", err)
		}

		for i := range ips {
			tcpaddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(ips[i].String(), port))
			if err != nil {
				return Client{}, fmt.Errorf("unexpected error while resolving memcached address: %w", err)
			}
			naddr = append(naddr, newStaticAddr(tcpaddr))
		}
	}

	c := memcache.NewFromSelector(&statiicServerList{addrs: naddr})

	logger = logger.With().
		Str("component", "memcached client").
		Logger()
	if err := c.Ping(); err != nil {
		// problem with memcached may be temporarily
		logger.Warn().
			Strs("addresses", cfg.Addresses).
			Err(err).
			Msg("ping: failed")
	} else {
		logger.Info().
			Msg("ping: ok")
	}

	return Client{c}, nil
}

func (c *Client) Close(_ context.Context) error {
	return c.Client.Close()
}
