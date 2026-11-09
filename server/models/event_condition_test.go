package models

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestNewServerEventCarriesSourceAndRequestTrace(t *testing.T) {
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled, Remote: true,
	}))

	event := NewServerEvent(ctx).Build()
	if event.Source == nil || *event.Source != EventSourceServer {
		t.Fatalf("source = %v, want %q", event.Source, EventSourceServer)
	}
	if event.TraceID == nil || *event.TraceID != traceID.String() || event.ParentID == nil || *event.ParentID != spanID.String() {
		t.Fatalf("trace context = %v/%v, want the request span's", event.TraceID, event.ParentID)
	}

	// Without a span, and with no context at all, no trace context is minted.
	for _, c := range []context.Context{context.Background(), nil} {
		if event := NewServerEvent(c).Build(); event.TraceID != nil || event.ParentID != nil {
			t.Fatalf("NewServerEvent(%v) minted trace context %v/%v", c, event.TraceID, event.ParentID)
		}
	}
}
