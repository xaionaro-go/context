package context

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Fields map[string]interface{}

type TraceID string

func (traceID TraceID) String() string {
	return string(traceID)
}

func NewTraceID() TraceID {
	return TraceID(uuid.New().String())
}

type Context interface {
	context.Context

	TraceID() TraceID
	WithTraceID(TraceID) Context

	Logger() Logger
	WithLogger(logger Logger) Context

	Metrics() Metrics
	WithMetrics(Metrics) Context

	Tracer() Tracer
	WithTracer(Tracer) Context

	WithTag(key string, value interface{}) Context
	WithTags(Fields) Context

	WithField(key string, value interface{}) Context
	WithFields(Fields) Context

	Finish()
}

type Tracer interface {
	WithTag(key string, value interface{}) Tracer
	WithTags(Fields) Tracer
	Finish()
}

type Metrics interface {
	WithTag(key string, value interface{}) Metrics
	WithTags(Fields) Metrics
}

type SimpleContext struct {
	context.Context
	TraceIDValue    TraceID
	LoggerInstance  Logger
	MetricsInstance Metrics
	TracerInstance  Tracer
}

func Background() Context {
	return NewContext(context.Background(), "", &dummyLogger{}, nil, nil)
}

func NewContext(stdCtx context.Context, traceID TraceID, logger MinimalLogger, metrics Metrics, tracer Tracer) *SimpleContext {
	if stdCtx == nil {
		stdCtx = Background()
	}
	if traceID == "" {
		traceID = NewTraceID()
	}

	ctx := &SimpleContext{
		TraceIDValue:    traceID,
		Context:         stdCtx,
		MetricsInstance: metrics,
		TracerInstance:  tracer,
	}

	if logger, ok := logger.(Logger); ok {
		ctx.LoggerInstance = logger
	} else {
		ctx.LoggerInstance = WrapLogger(logger)
	}

	return ctx
}

type CancelFunc = context.CancelFunc

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	stdCtx, cancelFunc := context.WithTimeout(parent, timeout)
	return NewContext(stdCtx, parent.TraceID(), parent.Logger(), parent.Metrics(), parent.Tracer()), cancelFunc
}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	stdCtx, cancelFunc := context.WithDeadline(parent, d)
	return NewContext(stdCtx, parent.TraceID(), parent.Logger(), parent.Metrics(), parent.Tracer()), cancelFunc
}

func WithValue(parent Context, key, val interface{}) Context {
	stdCtx := context.WithValue(parent, key, val)
	return NewContext(stdCtx, parent.TraceID(), parent.Logger(), parent.Metrics(), parent.Tracer())
}

func WithCancel(parent Context) (Context, CancelFunc) {
	stdCtx, cancelFunc := context.WithCancel(parent)
	return NewContext(stdCtx, parent.TraceID(), parent.Logger(), parent.Metrics(), parent.Tracer()), cancelFunc
}

func (ctx SimpleContext) TraceID() TraceID {
	return ctx.TraceIDValue
}

func (ctx SimpleContext) WithTraceID(traceID TraceID) Context {
	ctx.TraceIDValue = traceID
	return &ctx
}

func (ctx *SimpleContext) Logger() Logger {
	return ctx.LoggerInstance
}

func (ctx SimpleContext) WithLogger(logger Logger) Context {
	ctx.LoggerInstance = logger
	return &ctx
}

func (ctx *SimpleContext) Metrics() Metrics {
	return ctx.MetricsInstance
}

func (ctx SimpleContext) WithMetrics(metrics Metrics) Context {
	ctx.MetricsInstance = metrics
	return &ctx
}

func (ctx *SimpleContext) Tracer() Tracer {
	return ctx.TracerInstance
}

func (ctx SimpleContext) WithTracer(tracer Tracer) Context {
	ctx.TracerInstance = tracer
	return &ctx
}

func (ctx *SimpleContext) Finish() {
	ctx.TracerInstance.Finish()
}

func (ctx SimpleContext) WithTag(key string, value interface{}) Context {
	if ctx.MetricsInstance != nil {
		ctx.MetricsInstance = ctx.MetricsInstance.WithTag(key, value)
	}
	if ctx.TracerInstance != nil {
		ctx.TracerInstance = ctx.TracerInstance.WithTag(key, value)
	}
	ctx.LoggerInstance = ctx.LoggerInstance.WithField(key, value)
	return &ctx
}

func (ctx SimpleContext) WithTags(fields Fields) Context {
	if ctx.MetricsInstance != nil {
		ctx.MetricsInstance = ctx.MetricsInstance.WithTags(fields)
	}
	if ctx.TracerInstance != nil {
		ctx.TracerInstance = ctx.TracerInstance.WithTags(fields)
	}
	ctx.LoggerInstance = ctx.LoggerInstance.WithFields(fields)
	return &ctx
}

func (ctx SimpleContext) WithField(key string, value interface{}) Context {
	ctx.LoggerInstance = ctx.LoggerInstance.WithField(key, value)
	return &ctx
}

func (ctx SimpleContext) WithFields(fields Fields) Context {
	ctx.LoggerInstance = ctx.LoggerInstance.WithFields(fields)
	return &ctx
}

func (ctx SimpleContext) Value(key interface{}) interface{} {
	return ctx.Context.Value(key)
}
