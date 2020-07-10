package redis

import (
	"context"
	goredis "github.com/go-redis/redis/v7"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var cxtKey contextKey

type contextKey struct{}

type contextVal struct {
	starter time.Time
}

// implement go redis hook
type hook struct {
	schema string
}

func NewHook(schema string) *hook {
	return &hook{schema: schema}
}

// set ctx val for stats redis rt
func (h *hook) BeforeProcess(ctx context.Context, cmd goredis.Cmder) (context.Context, error) {
	if ctx == nil {
		ctx = context.TODO()
	}

	return context.WithValue(ctx, cxtKey, &contextVal{starter: time.Now()}), nil
}

func (h *hook) AfterProcess(ctx context.Context, cmd goredis.Cmder) error {
	var (
		status = GORedisStatus
	)

	rCtxVal, ok := ctx.Value(cxtKey).(*contextVal)
	if !ok || rCtxVal == nil {
		return nil
	}

	if err := cmd.Err(); err != nil {
		if err == goredis.Nil { // redis miss
			status = GoRedisMissStatus
		} else {
			status = GoRedisErrorStatus
		}
	}

	GoRedisError.With(prometheus.Labels{
		"schema": h.schema,
		"cmd":    cmd.Name(),
		"status": status,
	}).Inc()

	GoRedisDuration.With(prometheus.Labels{
		"schema": h.schema,
		"cmd":    cmd.Name(),
		"cost":   "redis_duration",
	}).Observe(time.Since(rCtxVal.starter).Seconds())

	GoRedisTotal.With(prometheus.Labels{"schema": h.schema, "cmd": "cmd_total"}).Inc()

	return nil
}

func (h *hook) BeforeProcessPipeline(ctx context.Context, cmd []goredis.Cmder) (context.Context, error) {
	if ctx == nil {
		ctx = context.TODO()
	}

	return context.WithValue(ctx, cxtKey, &contextVal{starter: time.Now()}), nil
}

func (h *hook) AfterProcessPipeline(ctx context.Context, cmd []goredis.Cmder) error {
	var (
		status = GORedisStatus
	)

	rCtxVal, ok := ctx.Value(cxtKey).(*contextVal)
	if !ok || rCtxVal == nil {
		return nil
	}

	if err := cmd[0].Err(); err != nil {
		if err == goredis.Nil { // redis miss
			status = GoRedisMissStatus
		} else {
			status = GoRedisErrorStatus
		}
	}

	GoRedisError.With(prometheus.Labels{
		"schema": h.schema,
		"cmd":    "pipeline",
		"status": status,
	}).Inc()

	GoRedisDuration.With(prometheus.Labels{
		"schema": h.schema,
		"cmd":    "pipeline",
		"cost":   "redis_duration",
	}).Observe(time.Since(rCtxVal.starter).Seconds())

	GoRedisTotal.With(prometheus.Labels{"schema": h.schema, "cmd": "cmd_total"}).Inc()

	return nil
}
