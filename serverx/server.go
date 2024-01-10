package serverx

import (
	"log/slog"
	"net/http"
	"time"
)

func Run(svr *http.Server) {
	if svr.Handler == nil {
		panic("RequiredHttpHandler")
	}
	if svr.Addr == "" {
		svr.Addr = ":8888"
	}
	if svr.ReadTimeout == 0 {
		svr.ReadTimeout = 60 * time.Second
	}
	if svr.WriteTimeout == 0 {
		svr.WriteTimeout = 60 * time.Second
	}

	slog.Info("server is listening " + svr.Addr)
	svr.ListenAndServe()
}
