package carrier

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Carrier is basically the same as propagation.MapCarrier
// But we reimplement it for easier interop (e.g. serialization and deserialization without caring map being nil)
// Usage: embed into a struct for passing message, e.g.:
//
//	type MyEventMessage struct {
//		Payload []byte `json:"payload"`
//		OtherData string `json:"other_data"`
//		TracingInfo Carrier `json:"tracing_info"`
//	}
//
// Then do normal propagation in the standard way:
//
// otel.GetTextMapPropagator().Inject(ctx, &myEventMessage.Carrier)
// ctx := otel.GetTextMapPropagator().Extract(ctx, &myEventMessage.Carrier)
//
// Or by helper methods:
//
// myEventMessage.Carrier.InjectContext(ctx)
// ctx := myEventMessage.Carrier.PropagateIntoContext(ctx)
type Carrier struct {
	Carrier propagation.MapCarrier
}

var _ propagation.TextMapCarrier = (*Carrier)(nil)
var _ json.Marshaler = (*Carrier)(nil)
var _ json.Unmarshaler = (*Carrier)(nil)

// ensureCarrierNotNil is a helper method which initializes the MapCarrier
// It should be called before any access on the MapCarrier to prevent nil pointer dereference
func (d *Carrier) ensureCarrierNotNil() {
	if d.Carrier == nil {
		d.Carrier = make(propagation.MapCarrier)
	}
}

// Get implements propagation.TextMapCarrier by delegating the call to the underlying MapCarrier
func (d *Carrier) Get(s string) string {
	d.ensureCarrierNotNil()
	return d.Carrier.Get(s)
}

// Set implements propagation.TextMapCarrier by delegating the call to the underlying MapCarrier
func (d *Carrier) Set(key, value string) {
	d.ensureCarrierNotNil()
	d.Carrier.Set(key, value)
}

// Keys implements propagation.TextMapCarrier by delegating the call to the underlying MapCarrier
func (d *Carrier) Keys() []string {
	d.ensureCarrierNotNil()
	return d.Carrier.Keys()
}

// MarshalJSON implements json.Marshaler by delegating the call to the underlying MapCarrier
func (d Carrier) MarshalJSON() ([]byte, error) {
	d.ensureCarrierNotNil()
	return json.Marshal(d.Carrier)
}

// UnmarshalJSON implements json.Unmarshaler by delegating the call to the underlying MapCarrier
func (d *Carrier) UnmarshalJSON(bz []byte) error {
	return json.Unmarshal(bz, &d.Carrier)
}

// getPropagator returns a "sane" default implementation of propagation.TextMapPropagator.
// The default value of otel.GetTextMapPropagator() without setup will be a no-op propagator.
// We want to make sure that tracing works, so we prepend our propagator by a composite propagator.
// If some use cases don't want propagation, ClearContext can be used.
func getPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, otel.GetTextMapPropagator())
}

// PropagateIntoContext propagates the tracing info into a new context based on the given context
func (d *Carrier) PropagateIntoContext(ctx context.Context) context.Context {
	// extract the tracing info from the carrier d into ctx
	return getPropagator().Extract(ctx, d)
}

// NewCarrierFromContext initializes a DistributedTracingInfo using the tracing info from a context
func NewCarrierFromContext(ctx context.Context) Carrier {
	var d Carrier
	d.InjectContext(ctx)
	return d
}

// InjectContext propagates the tracing info from the provided context into the DistributedTracingInfo
func (d *Carrier) InjectContext(ctx context.Context) *Carrier {
	getPropagator().Inject(ctx, d)
	return d
}

// ClearContext clears the tracing info in the DistributedTracingInfo.
// It can be used in cases which don't want propagation, e.g. sending to 3rd party services.
func (d *Carrier) ClearContext() *Carrier {
	d.Carrier = nil
	d.ensureCarrierNotNil()
	return d
}

// GetTraceParent returns the traceparent header value, which provides a convenient way to get a
// representation of the current trace context.
// Note that this assumes that the previous TextMapPropagator is / contains TraceContext,
// if it's not the case, traceparent will be empty even if there is a trace context in the carrier.
func (d *Carrier) GetTraceParent() string {
	return d.Get("traceparent")
}
