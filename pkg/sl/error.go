package sl

import "log/slog"

// Err wraps errors as slog attribute
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
