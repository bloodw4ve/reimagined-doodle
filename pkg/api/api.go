package api

import (
	"encoding/json"
	"mod36a41/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// server program interface
type API struct {
	db     storage.DBinterface
	router *mux.Router
}

// API construct
func New(db storage.DBinterface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// router
func (api *API) Router() *mux.Router {
	return api.router
}

// API handlers
func (api *API) endpoints() {
	api.router.HandleFunc("/news/{n}", api.lastNewsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (api *API) lastNewsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	news, err := api.db.LastNews(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(news)
}
