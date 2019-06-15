package apple

import (
	"net/http"
	"net/http/httptest"

	"github.com/musicmash/musicmash/internal/db"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	db.DbMgr = db.NewFakeDatabaseMgr()
}

func teardown() {
	server.Close()
	_ = db.DbMgr.DropAllTables()
	_ = db.DbMgr.Close()
}
