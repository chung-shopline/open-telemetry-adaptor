# Open Telemetry Adapters

This repository contains carrier adapters for working with Open Telemetry traces.

## carrier

A module implements the `propagation.TextMapCarrier` interface, which basically wraps `propagation.MapCarrier` with more convenient handling like nil map checks.

## pulsar

A module implements the `propagation.TextMapCarrier` interface for Pulsar producer message and Pulsar consumer message.

Usage in producer:

```go
// Producer
import pulsar_otel "github.com/shoplineapp/open-telemetry-adapters/pulsar"

// ctx contains the span context
ctx, span := otel.Tracer("my-app").Start(ctx, "pulsar-producer")
msg := pulsar.ProducerMessage{
  // ...
}
// propagate the span context to the message using the carrier
otel.GetTextMapPropagator().Inject(ctx, pulsar_otel.NewProducerMessageCarrier(&msg))
```

Usage in consumer:

```go
// Producer
import pulsar_otel "github.com/shoplineapp/open-telemetry-adapters/pulsar"

func myConsumer(ctx context.Context, msg pulsar.ConsumerMessage) error {
  // extract tracing info from message into ctx
  ctx = otel.GetTextMapPropagator().Extract(ctx, pulsar_otel.NewConsumerMessageCarrier(msg))
  // start a span which propagates the span context from the message
  ctx, span := otel.Tracer("my-tracer").Start(ctx, "pulsar-consumer")
}
```
