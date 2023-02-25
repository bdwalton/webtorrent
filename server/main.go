// package server implements the webtorrent http server.
package server

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/ini.v1"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello, world!</h1>")
}

func ListenAndServe(ctx context.Context, cfg *ini.File) error {
	http.HandleFunc("/", rootHandler)
	return http.ListenAndServe(":"+cfg.Section("server").Key("port").String(), nil)
}
