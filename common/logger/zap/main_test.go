package zap_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/lruggieri/fxnow/common/logger"
	loggerzap "github.com/lruggieri/fxnow/common/logger/zap"
)

var sinks = map[string]*testLogSink{}

func TestMain(m *testing.M) {
	_ = zap.RegisterSink("test", func(u *url.URL) (zap.Sink, error) {
		return getSink(u.Hostname()), nil
	})

	os.Exit(m.Run())
}

func createLogger(sinkName string) *loggerzap.ImplLogger {
	sinks[sinkName] = &testLogSink{}

	return loggerzap.New(loggerzap.Config{
		Development:    false,
		Level:          logger.LevelDebug,
		OutputPaths:    []string{fmt.Sprintf("test://%s", sinkName)},
		SkipLineEnding: true,
	})
}

func getSink(name string) *testLogSink {
	if sink, ok := sinks[name]; ok {
		return sink
	}

	panic(errors.Errorf("unknown sink name: %s", name))
}

type testLogSink struct {
	line []map[string]interface{}
}

func (t *testLogSink) Write(p []byte) (int, error) {
	var data map[string]interface{}

	err := json.Unmarshal(p, &data)
	if err != nil {
		return 0, err
	}

	t.line = append(t.line, data)

	return len(p), nil
}

func (*testLogSink) Sync() error {
	return nil
}

func (*testLogSink) Close() error {
	return nil
}
