package ui

import (
	"context"
	"fmt"

	"yiu-ops/internal/app"
)

type Service struct {
	appCtx *app.Context
}

func NewService(appCtx *app.Context) *Service {
	return &Service{appCtx: appCtx}
}

func (s *Service) Start(ctx context.Context, port int) error {
	if s.appCtx != nil && s.appCtx.Logger != nil {
		s.appCtx.Logger.Info(fmt.Sprintf("在端口 %d 上启动服务器", port))
	}
	return nil
}
