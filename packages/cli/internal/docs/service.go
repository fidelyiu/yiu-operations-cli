package docs

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"
)

type Service struct {
	buildFS fs.FS
}

func NewService(buildFS fs.FS) *Service {
	return &Service{buildFS: buildFS}
}

func (s *Service) Serve(ctx context.Context, host string, port int) error {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: http.FileServer(http.FS(s.buildFS)),
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}
