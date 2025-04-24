package logger

import (
	"context"
	"encoding/json"
	"os"

	"github.com/conan194351/BTL-KTPM/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	logger *zap.Logger
}

func (l *ZapLogger) GetLogger() *zap.Logger {
	return l.logger
}

func (l *ZapLogger) Info(msg string, ctx context.Context) {

	l.logger.Info(msg, getMetadata(ctx)...)
}

func (l *ZapLogger) Debug(msg string, ctx context.Context) {
	l.logger.Debug(msg, getMetadata(ctx)...)
}

func (l *ZapLogger) Warn(msg string, ctx context.Context) {
	l.logger.Warn(msg, getMetadata(ctx)...)
}

func (l *ZapLogger) Fatal(err error, msg string, ctx context.Context) {
	l.logger.Fatal(msg, append(getMetadata(ctx), zap.Error(err))...)
}

func (l *ZapLogger) Panic(err error, msg string, ctx context.Context) {
	l.logger.Panic(msg, append(getMetadata(ctx), zap.Error(err))...)
}

func (l *ZapLogger) Error(err error, msg string, ctx context.Context) {
	l.logger.Error(msg, append(getMetadata(ctx), zap.Error(err))...)
}

func (l *ZapLogger) CtxWithMetadata(ctx context.Context, value interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, LOG_METADATA, value)

}

func CtxWithMetadata(ctx context.Context, value interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, LOG_METADATA, value)

}

func getMetadata(ctx context.Context) []zap.Field {

	if ctx == nil {
		return []zap.Field{}
	}

	if ctx.Value(LOG_METADATA) == nil {
		return []zap.Field{}
	}

	var zapFields []zap.Field

	metadataBytes, err := json.Marshal(ctx.Value(LOG_METADATA))
	if err != nil {
		return zapFields
	}

	metadata := make(map[string]interface{})
	err = json.Unmarshal(metadataBytes, &metadata)

	if err != nil {
		return zapFields
	}

	for key, value := range metadata {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	return zapFields
}

func NewZapLogger(
	name string,
	haveCaller bool,
) *ZapLogger {
	atom := zap.NewAtomicLevel()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

	zapFields := []zap.Field{
		zap.String("service", name),
	}

	opts := []zap.Option{
		zap.Fields(zapFields...),
	}
	if haveCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zap.ErrorLevel))
	}

	fileWriter := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.GetConfig().App.LogPath,
			LocalTime:  true,
			MaxSize:    100,
			MaxBackups: 10,
		}),
		atom)

	stdoutWriter := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(os.Stdout),
		atom)

	logger := zap.New(zapcore.NewTee(fileWriter, stdoutWriter), opts...)

	defer logger.Sync()

	atom.SetLevel(zap.DebugLevel)

	return &ZapLogger{
		logger: logger,
	}
}
