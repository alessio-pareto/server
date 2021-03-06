package server

import (
	"fmt"
	"time"

	"github.com/felixge/httpsnoop"
)

func (srv *Server) WriteLogStart(t time.Time) {
	fmt.Fprint(srv.LogFile, "\n     /\\ /\\ /\\                                            /\\ /\\ /\\")
	fmt.Fprint(srv.LogFile, "\n     <> <> <> - [" + t.Format("02/Jan/2006:15:04:05") + "] - SERVER ONLINE - <> <> <>")
	fmt.Fprint(srv.LogFile, "\n     \\/ \\/ \\/                                            \\/ \\/ \\/\n\n")
}

func (srv *Server) WriteLogClosure(t time.Time) {
	fmt.Fprint(srv.LogFile, "\n     /\\ /\\ /\\                                             /\\ /\\ /\\")
	fmt.Fprint(srv.LogFile, "\n     <> <> <> - [" + t.Format("02/Jan/2006:15:04:05") + "] - SERVER OFFLINE - <> <> <>")
	fmt.Fprint(srv.LogFile, "\n     \\/ \\/ \\/                                             \\/ \\/ \\/\n\n")
}

func (route *Route) logInfo(metrics httpsnoop.Metrics) {
	lock := "\U0000274C"
	if route.Secure {
		lock = "\U0001F512"
	}

	fmt.Fprintf(route.Srv.LogFile, "   Info: %-16s - [%s] - %-4s %-65s %s %d %10.3f MB - (%6d ms) \u279C %s (%s) via %s\n",
		route.RemoteAddress,
		time.Now().Format("02/Jan/2006:15:04:05"),
		route.R.Method,
		route.logRequestURI,
		lock,
		metrics.Code,
		(float64(metrics.Written)/1000000.),
		time.Since(route.ConnectionTime).Milliseconds(),
		route.Website.Name,
		route.Domain.Name,
		route.Host,
	)
}

func (route *Route) logWarning(metrics httpsnoop.Metrics) {
	lock := "\U0000274C"
	if route.Secure {
		lock = "\U0001F512"
	}

	fmt.Fprintf(route.Srv.LogFile, "Warning: %-16s - [%s] - %-4s %-65s %s %d %10.3f MB - (%6d ms) \u279C %s (%s) via %s \u279C %s\n",
		route.RemoteAddress,
		time.Now().Format("02/Jan/2006:15:04:05"),
		route.R.Method,
		route.logRequestURI,
		lock,
		metrics.Code,
		(float64(metrics.Written)/1000000.),
		time.Since(route.ConnectionTime).Milliseconds(),
		route.Website.Name,
		route.Domain.Name,
		route.Host,
		route.logErrMessage,
	)
}

func (route *Route) logError(metrics httpsnoop.Metrics) {
	lock := "\U0000274C"
	if route.Secure {
		lock = "\U0001F512"
	}

	fmt.Fprintf(route.Srv.LogFile, "Error: %-16s - [%s] - %-4s %-65s %s %d %10.3f MB - (%6d ms) \u279C %s (%s) via %s \u279C %s\n",
		route.RemoteAddress,
		time.Now().Format("02/Jan/2006:15:04:05"),
		route.R.Method,
		route.logRequestURI,
		lock,
		metrics.Code,
		(float64(metrics.Written)/1000000.),
		time.Since(route.ConnectionTime).Milliseconds(),
		route.Website.Name,
		route.Domain.Name,
		route.Host,
		route.logErrMessage,
	)
}

func (route *Route) serveError() {
	if route.W.written || route.errTemplate == nil {
		return
	}

	err := route.errTemplate.Execute(route.W, struct{ Code int; Message string }{ Code: route.W.code, Message: route.errMessage })
	if err != nil {
		fmt.Fprintf(route.Srv.LogFile, "Error serving template file: %v\n", err)
		return
	}
}
