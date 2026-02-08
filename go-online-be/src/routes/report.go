package routes

import (
	"encoding/json"
	"net/http"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
	"github.com/go-chi/chi/v5"
)

func ReportRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.IsLoggedIn)
	r.Post("/", postReportHandler)
	r.Group(func(r chi.Router) {
		r.Use(middleware.IsModerator)
		r.Get("/", getReportsHandler)
		r.Delete("/{reportId}", deleteReportHandler)
	})
	return r
}

func getReportsHandler(w http.ResponseWriter, r *http.Request) {
	reports, err := gogame.LoadAllReportsFromDB()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, reports)
}

type createReportReq struct {
	GameID       string `json:"game_id"`
	FromUsername string `json:"username"`
}

func postReportHandler(w http.ResponseWriter, r *http.Request) {
	var req createReportReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := gogame.SaveReportToDB(req.GameID, req.FromUsername); err != nil {
		http.Error(w, "could not save report", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "report saved"})
}

func deleteReportHandler(w http.ResponseWriter, r *http.Request) {
	reportId := chi.URLParam(r, "reportId")
	if err := gogame.DeleteReportFromDB(reportId); err != nil {
		http.Error(w, "could not delete report", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "report deleted"})
}
