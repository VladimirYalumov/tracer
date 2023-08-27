package tracer

import (
	"fmt"
)

var defaultOptions = []Option{
	CollectorHost("localhost"),
	CollectorPort(4317),
}

type Option interface {
	apply(o *Options)
}

type (
	Noop          bool
	CollectorHost string
	CollectorPort uint16
)

type Options struct {
	noop bool
	host CollectorHost
	port CollectorPort
}

func BuildOptions(opts []Option) Options {
	options := Options{}

	opts = append(defaultOptions, opts...)
	for _, opt := range opts {
		opt.apply(&options)
	}

	return options
}

func (o Options) GetTarget() string {
	return fmt.Sprintf("%s:%d", o.host, o.port)
}

func (o Options) IsNoop() bool {
	return o.noop
}

func (n Noop) apply(o *Options) {
	o.noop = bool(n)
}

func (h CollectorHost) apply(o *Options) {
	o.host = h
}

func (p CollectorPort) apply(o *Options) {
	o.port = p
}
