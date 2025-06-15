package channel

import "time"

type Options struct {
	SingleChannelBufferSize uint
	PublishToChannelDelay   time.Duration
}

var DefautOptions = Options{
	SingleChannelBufferSize: 1024,
	PublishToChannelDelay:   1 * time.Second,
}

type OptionsSetter func(*Options)

func WithSingleChannelBufferSize(value uint) OptionsSetter {
	return func(o *Options) {
		o.SingleChannelBufferSize = value
	}
}

func WithPublishToChannelDelay(value time.Duration) OptionsSetter {
	return func(o *Options) {
		o.PublishToChannelDelay = value
	}
}
