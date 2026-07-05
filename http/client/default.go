package client

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

const (
	DefaultMaxIdleConns          = 100
	DefaultDialTimeout           = 5 * time.Second
	DefaultKeepAliveDuration     = 30 * time.Second
	DefaultIdleConnTimeout       = 90 * time.Second
	DefaultTLSHandshakeTimeout   = 4 * time.Second
	DefaultExpectContinueTimeout = 1 * time.Second
)

// transport returns a new http.Transport with idle connections and keepalives
// disabled.
func transport(opts Options) *http.Transport {
	transport := pooledTransport(opts)
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// pooledTransport returns a new http.Transport with similar default values to
// http.DefaultTransport. Do not use this for transient transports as it can
// leak file descriptors over time. Only use this for transports that will be
// re-used for the same host(s).
func pooledTransport(opts Options) *http.Transport {
	return &http.Transport{
		MaxIdleConns: opts.MaxIdleConns,
		Proxy:        http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   opts.DialTimeout,
			KeepAlive: opts.KeepAliveDuration,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:       opts.IdleConnTimeout,
		TLSHandshakeTimeout:   opts.TLSHandshakeTimeout,
		ExpectContinueTimeout: opts.ExpectContinueTimeout,
		ForceAttemptHTTP2:     true,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
}

type Options struct {
	// MaxIdleConns is the maximum number of idle connections a pooled
	// http.Client can have.
	MaxIdleConns int
	// DialTimeout is the timeout duration to use when attempting to connect to
	// the server.
	DialTimeout time.Duration
	// DisableKeepAlives disables the use of keep-alive sessions for the
	// http.Client.
	DisableKeepAlives bool
	// KeepAliveDuration is the duration to keep keel-alive sessions open.
	KeepAliveDuration time.Duration
	// IdleConnTimeout is the timeout duration for when a connection should be
	// considered idle.
	IdleConnTimeout time.Duration
	// TLSHandshakeTimeout is the timeout duration for the TLS handshake.
	TLSHandshakeTimeout time.Duration
	// ExpectContinueTimeout is the timeout duration to apply to the
	// ExpectContinue operation.
	ExpectContinueTimeout time.Duration
}

// SetDefaults sets any empty fields to their default values.
func (o *Options) SetDefaults() {
	if o.MaxIdleConns == 0 {
		o.MaxIdleConns = DefaultMaxIdleConns
	}
	if o.DialTimeout == 0 {
		o.DialTimeout = DefaultDialTimeout
	}
	if o.KeepAliveDuration == 0 {
		o.KeepAliveDuration = DefaultKeepAliveDuration
	}
	if o.TLSHandshakeTimeout == 0 {
		o.TLSHandshakeTimeout = DefaultTLSHandshakeTimeout
	}
	if o.IdleConnTimeout == 0 {
		o.IdleConnTimeout = DefaultIdleConnTimeout
	}
	if o.ExpectContinueTimeout == 0 {
		o.ExpectContinueTimeout = DefaultExpectContinueTimeout
	}
}

type Option func(*Options)

// WithDialTimeout sets the returned client's dial timeout.
func WithDialTimeout(v time.Duration) Option {
	return func(o *Options) {
		o.DialTimeout = v
	}
}

// New returns a new http.Client with similar default values to http.Client,
// but with a non-pooled Transport, idle connections disabled, and keepalives
// disabled.
func New(opts ...Option) *http.Client {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	o.SetDefaults()
	return &http.Client{
		Transport: transport(o),
	}
}

// NewPooled returns a new http.Client with similar default values to
// http.Client, but with a pooled Transport that has keep-alives and idle
// connections enabled. Do not use this function for transient clients as it
// can leak file descriptors over time. Only use this for clients that will be
// re-used for the same host(s).
func NewPooled(opts ...Option) *http.Client {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	o.SetDefaults()
	return &http.Client{
		Transport: pooledTransport(o),
	}
}
