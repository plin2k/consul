// Copyright 2018 Sergey Novichkov. All rights reserved.
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package consul

import (
	"fmt"
	"net"

	"github.com/gozix/di"
	"github.com/gozix/glue/v3"
	gzViper "github.com/gozix/viper/v3"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// Bundle implements the glue.Bundle interface.
type Bundle struct{}

// BundleName is default definition name.
const BundleName = "consul"

// Bundle implements the glue.Bundle interface.
var _ glue.Bundle = (*Bundle)(nil)

// NewBundle create bundle instance.
func NewBundle() *Bundle {
	return new(Bundle)
}

func (b *Bundle) Name() string {
	return BundleName
}

func (b *Bundle) Build(builder di.Builder) error {
	return builder.Provide(b.provideClient)
}

func (b *Bundle) DependsOn() []string {
	return []string{
		gzViper.BundleName,
	}
}

func (b *Bundle) provideClient(cfg *viper.Viper) (*api.Client, error) {
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
}
