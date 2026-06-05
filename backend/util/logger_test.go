package util

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	l := NewLogger()
	assert.NotNil(t, l)
	assert.NotNil(t, l.log)
}

func TestLoggerStart(t *testing.T) {
	l := NewLogger()
	ctx := context.Background()
	before := time.Now()
	start := l.Start(ctx, "TestMethod")
	after := time.Now()

	assert.True(t, start.After(before) || start.Equal(before))
	assert.True(t, start.Before(after) || start.Equal(after))
}

func TestLoggerFinish(t *testing.T) {
	l := NewLogger()
	ctx := context.Background()
	start := l.Start(ctx, "TestMethod")

	// no error path
	l.Finish(ctx, "TestMethod", start, nil)

	// error path
	l.Finish(ctx, "TestMethod", start, errors.New("something went wrong"))
}
