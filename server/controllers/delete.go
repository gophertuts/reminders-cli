package controllers

import (
	"log"
	"net/http"

	"github.com/gophertuts/reminders-cli/server/services"
)

type deleter interface {
	Delete(ids []int) services.IDsResponse
}

func deleteReminders(service deleter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := parseIDsParam(r.Context())
		idsRes := service.Delete(ids)
		if len(idsRes.NotFoundIDs) > 0 {
			log.Printf("could not delete ids: %v\n", idsRes.NotFoundIDs)
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
