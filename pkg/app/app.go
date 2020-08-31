package app

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type FooApp struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func (app *FooApp) Handler(w http.ResponseWriter, req *http.Request) {
	app.Logger.Info("Receiver request", zap.String("path", req.URL.Path), zap.String("method", req.Method))
	if req.URL.Path == "/db" {
		rows, err := app.DB.Query(`show tables;`)
		if err != nil {
			app.Logger.Error("Failed to list tables", zap.Error(err))
			http.Error(w, "Failed to call DB", http.StatusInternalServerError)
			return
		}

		ts := make([]string, 0)
		for rows.Next() {
			var table string
			rows.Scan(&table)
			ts = append(ts, table)
		}
		b, err := json.Marshal(ts)
		if err != nil {
			http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
