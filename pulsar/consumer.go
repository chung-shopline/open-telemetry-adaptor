package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*ConsumerMessageCarrier)(nil)

// ConsumerMessageCarrier implements propagation.TextMapCarrier for a Pulsar consumer message.
// It is used for extracting tracing context from a Pulsar consumer message.
type ConsumerMessageCarrier struct {
	msg pulsar.Message
}

func NewConsumerMessageCarrier(msg pulsar.Message) ConsumerMessageCarrier {
	return ConsumerMessageCarrier{msg: msg}
}

func (c ConsumerMessageCarrier) Get(key string) string {
	var properties = c.msg.Properties()
	if properties == nil {
		properties = make(map[string]string)
	}
	return properties[key]
}

func (c ConsumerMessageCarrier) Set(key, val string) {
	var properties = c.msg.Properties()
	if properties == nil {
		properties = make(map[string]string)
	}
	properties[key] = val
}

func (c ConsumerMessageCarrier) Keys() []string {
	var properties = c.msg.Properties()
	if len(properties) == 0 {
		return make([]string, 0)
	}
	out := make([]string, len(properties))
	for key := range properties {
		out = append(out, properties[key])
	}
	return out
}
