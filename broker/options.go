package broker

import "context"

////////////////////////////////////////////////////////////
// Options
////////////////////////////////////////////////////////////

type Options struct {
	Context context.Context
}

type Option func(*Options)

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func NewOptions(opts ...Option) Options {
	opt := Options{
		Context: context.Background(),
	}

	opt.Apply(opts...)

	return opt
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func OptionsContextWithValue(k, v interface{}) Option {
	return func(o *Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, k, v)
	}
}
