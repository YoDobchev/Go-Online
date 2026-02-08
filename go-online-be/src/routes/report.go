package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	gogame "github.com/YoDobchev/Go-Online/src/game/go"
	"github.com/YoDobchev/Go-Online/src/middleware"
	"github.com/go-chi/chi/v5"
)

func ReportRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.IsLoggedIn)
	r.Get("/", getReportsHandler)
	r.Post("/", postReportHandler)
	r.Delete("/{reportId}", deleteReportHandler)
	return r
}

func getReportsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromCtx(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" && user.Role != "moderator" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

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

	fmt.Println(req.GameID)

	if err := gogame.SaveReportToDB(req.GameID, req.FromUsername); err != nil {
		http.Error(w, "could not save report", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "report saved"})
}

func deleteReportHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromCtx(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" && user.Role != "moderator" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	reportId := chi.URLParam(r, "reportId")
	if err := gogame.DeleteReportFromDB(reportId); err != nil {
		http.Error(w, "could not delete report", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "report deleted"})
}
