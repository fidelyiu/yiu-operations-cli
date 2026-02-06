package app

import "log/slog"

type Context struct {
	Logger *slog.Logger
}

func New() *Context {
	return &Context{Logger: slog.Default()}
}
