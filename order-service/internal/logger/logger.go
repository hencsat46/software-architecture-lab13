package logger

import "go.uber.org/zap"

type Logger struct{ *zap.Logger }

func New(level string) (*Logger, error) {
	c := zap.NewProductionConfig()
	if err := c.Level.UnmarshalText([]byte(level)); err != nil {
		return nil, err
	}
	l, err := c.Build()
	return &Logger{l}, err
}

func (l *Logger) WithContext(layer, scope string) *Logger {
	return &Logger{l.With(zap.String("layer", layer), zap.String("scope", scope))}
}

func (l *Logger) WithMethod(method string) *Logger {
	return &Logger{l.With(zap.String("method", method))}
}
