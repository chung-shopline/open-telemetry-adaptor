package pulsar

import (
	"context"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// FakePulsarMessageForCarrierOnly is used for testing the ConsumerCarrier.
// We need a pulsar.Message, but the interface contains a lot of methods that we don't need for testing.
// Therefore we use an "anti-pattern" here: embed a nil pulsar.Message and only implement the methods we need for the carrier.
// As it's basically a nil implementation of pulsar.Message, IT SHOULD NOT BE USED IN OTHER CASES (including other test cases which needs a pulsar.Message).
type FakePulsarMessageForCarrierOnly struct {
	pulsar.Message
	msg *pulsar.ProducerMessage
}

func (m FakePulsarMessageForCarrierOnly) Properties() map[string]string {
	return m.msg.Properties
}

func TestPulsarCarrier(t *testing.T) {
	traceState, err := trace.ParseTraceState("a=b,c=d")
	if err != nil {
		panic(err)
	}
	traceID := trace.TraceID{0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11}
	spanID := trace.SpanID{0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22}
	srcCtx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceState: traceState,
	}))

	producerMsg := pulsar.ProducerMessage{}
	propagation.TraceContext{}.Inject(srcCtx, NewProducerMessageCarrier(&producerMsg))

	consumerMsg := FakePulsarMessageForCarrierOnly{msg: &producerMsg}
	dstCtx := propagation.TraceContext{}.Extract(context.Background(), ConsumerMessageCarrier{msg: &consumerMsg})
	dstSpanContext := trace.SpanContextFromContext(dstCtx)
	assert.Equal(t, traceID, dstSpanContext.TraceID())
	assert.Equal(t, spanID, dstSpanContext.SpanID())
	assert.Equal(t, traceState, dstSpanContext.TraceState())
}
