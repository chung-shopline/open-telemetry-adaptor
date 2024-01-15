package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"go.opentelemetry.io/otel/propagation"
)

var _ propagation.TextMapCarrier = (*ProducerMessageCarrier)(nil)

// ProducerMessageCarrier implements propagation.TextMapCarrier for a Pulsar producer message.
// It is used for propagating tracing context into a Pulsar producer message.
// Consumer should use ConsumerMessageCarrier to extract the propagated trace context from the Pulsar consumer message.
type ProducerMessageCarrier struct {
	msg *pulsar.ProducerMessage
}

func NewProducerMessageCarrier(msg *pulsar.ProducerMessage) ProducerMessageCarrier {
	return ProducerMessageCarrier{msg: msg}
}
func (c ProducerMessageCarrier) Get(key string) string {
	var properties = c.msg.Properties
	if properties == nil {
		properties = make(map[string]string)
	}
	return properties[key]
}

func (c ProducerMessageCarrier) Set(key, val string) {
	var properties = c.msg.Properties
	if properties == nil {
		properties = make(map[string]string)
		c.msg.Properties = properties
	}
	properties[key] = val
}

func (c ProducerMessageCarrier) Keys() []string {
	var properties = c.msg.Properties
	if len(properties) == 0 {
		return make([]string, 0)
	}
	out := make([]string, len(properties))
	for key := range properties {
		out = append(out, properties[key])
	}
	return out
}
