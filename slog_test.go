package slog_test

import (
	"testing"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/stretchr/testify/assert"
)

func ExampleNew() {
	slog.Info("info log message")
	slog.Infof("info log %s", "message")
}

func TestInfof(t *testing.T) {
	slog.AddHandler(handler.NewConsoleHandler(slog.AllLevels))

	h2 := handler.NewConsoleHandler(slog.AllLevels)
	h2.SetFormatter(slog.NewJSONFormatter(slog.StringMap{
		"level": "levelName",
		"message": "msg",
		"data": "params",
	}))

	slog.AddHandler(h2)
	slog.AddProcessor(slog.AddHostname())

	slog.Infof("info %s", "message")
}

func TestLevelName(t *testing.T) {
	for level, wantName := range slog.LevelNames {
		realName := slog.LevelName(level)
		assert.Equal(t, wantName, realName)
	}

	assert.Equal(t, "unknown", slog.LevelName(20))
}

func TestName2Level(t *testing.T) {
	for wantLevel, name := range slog.LevelNames {
		level, err := slog.Name2Level(name)
		assert.NoError(t, err)
		assert.Equal(t, wantLevel, level)
	}

	level, err := slog.Name2Level("")
	assert.NoError(t, err)
	assert.Equal(t, slog.InfoLevel, level)

	level, err = slog.Name2Level("unknown")
	assert.Error(t, err)
	assert.Equal(t, slog.Level(0), level)
}
