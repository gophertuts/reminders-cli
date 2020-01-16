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
		checkHTTPMethod(w, r, http.MethodDelete)
		var body struct {
			IDs []int `json:"ids"`
		}
		jsonDecode(r.Body, &body)
		ids := service.Delete(body.IDs)
		if len(ids.NotFoundIDs) > 0 {
			log.Printf("could not delete ids: %v\n", ids.NotFoundIDs)
		}
		type response struct {
			NotFoundIDs []int `json:"not_found_ids"`
			DeletedIDs  []int `json:"deleted_ids"`
		}
		jsonEncode(w, response{
			NotFoundIDs: ids.NotFoundIDs,
			DeletedIDs:  ids.DeletedIDs,
		})
	})
}
