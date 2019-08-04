// Package consul provide gozix bundle implementation.
package consul

import (
	"fmt"
	"net"

	"github.com/gozix/viper/v2"
	"github.com/hashicorp/consul/api"
	"github.com/sarulabs/di/v2"
)

type (
	// Bundle implements the glue.Bundle interface.
	Bundle struct{}

	// Client is type alias of api.Client
	Client = api.Client

	// KVPair type alias of api.KVPair.
	KVPair = api.KVPair

	// QueryMeta type alias of api.QueryMeta.
	QueryMeta = api.QueryMeta
)

// BundleName is default definition name.
const BundleName = "consul"

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

// Key implements the glue.Bundle interface.
func (b *Bundle) Name() string {
	return BundleName
}

// Build implements the glue.Bundle interface.
func (b *Bundle) Build(builder *di.Builder) error {
	return builder.Add(di.Def{
		Name: BundleName,
		Build: func(ctn di.Container) (_ interface{}, err error) {
			var cfg *viper.Viper
			if err = ctn.Fill(viper.BundleName, &cfg); err != nil {
				return nil, err
			}

			var c = api.DefaultConfig()
			c.Address = net.JoinHostPort(
				cfg.GetString(fmt.Sprintf("%s.host", BundleName)),
				cfg.GetString(fmt.Sprintf("%s.port", BundleName)),
			)

			var key = fmt.Sprintf("%s.datacenter", BundleName)
			if cfg.IsSet(key) {
				c.Datacenter = cfg.GetString(key)
			}

			key = fmt.Sprintf("%s.scheme", BundleName)
			if cfg.IsSet(key) {
				c.Scheme = cfg.GetString(key)
			}

			key = fmt.Sprintf("%s.token", BundleName)
			if cfg.IsSet(key) {
				c.Token = cfg.GetString(key)
			}

			key = fmt.Sprintf("%s.wait_time", BundleName)
			if cfg.IsSet(key) {
				c.WaitTime = cfg.GetDuration(key)
			}

			return api.NewClient(c)
		},
	})
}

// DependsOn implements the glue.DependsOn interface.
func (b *Bundle) DependsOn() []string {
	return []string{viper.BundleName}
}
