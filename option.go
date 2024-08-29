package telegram

import "time"

type option struct {
	proxyUrl string
	poller   time.Duration
}

type Option func(*option)

func newOptions(opts ...Option) option {
	var options option
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

func WithProxy(u string) Option {
	return func(opts *option) {
		opts.proxyUrl = u
	}
}

func WithPoller(t time.Duration) Option {
	return func(opts *option) {
		opts.poller = t
	}
}
