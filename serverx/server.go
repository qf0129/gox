package serverx

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/qf0129/gox/constx"
)

func Run(server *http.Server) {
	if server.Handler == nil {
		panic("RequiredHttpHandler")
	}
	if server.Addr == "" {
		server.Addr = constx.DefaultListenAddr
	}
	if server.ReadTimeout == 0 {
		server.ReadTimeout = constx.DefaultReadTimeout * time.Second
	}
	if server.WriteTimeout == 0 {
		server.WriteTimeout = constx.DefaultWriteTimeout * time.Second
	}
	if err := server.ListenAndServe(); err != nil {
		slog.Error(err.Error())
	} else {
		slog.Info("### server is listening " + server.Addr)
	}
}
