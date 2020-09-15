package app

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var webTemplate *template.Template = template.Must(template.ParseFiles("/var/run/ko/layout.html"))

type FooApp struct {
	DB     *sql.DB
	Bucket string
	Logger *zap.Logger
}

type PetData struct {
	Name    string    `json:"name,omitempty"`
	Owner   string    `json:"owner,omitempty"`
	Species string    `json:"species,omitempty"`
	Sex     string    `json:"sex,omitempty"`
	Birth   time.Time `json:"birth,omitempty"`
	Link    string    `json:"link,omitempty"`
}

func (app *FooApp) Handler(w http.ResponseWriter, req *http.Request) {
	app.Logger.Info("Receiver request", zap.String("path", req.URL.Path), zap.String("method", req.Method))

	rows, err := app.DB.Query(`SELECT name, owner, species, sex, birth FROM pet;`)
	if err != nil {
		app.Logger.Error("Failed to query pet table", zap.Error(err))
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	ds := make([]PetData, 0)
	for rows.Next() {
		d := PetData{}
		if err := rows.Scan(&d.Name, &d.Owner, &d.Species, &d.Sex, &d.Birth); err != nil {
			app.Logger.Error("Failed to parse row", zap.Error(err))
		} else {
			d.Link = fmt.Sprintf("https://storage.googleapis.com/%s/%s.jpg", app.Bucket, d.Name)
			ds = append(ds, d)
		}
	}

	w.WriteHeader(http.StatusOK)
	webTemplate.Execute(w, ds)
	return
}
