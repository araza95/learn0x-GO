package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	types "github.com/araza95/learn0x-GO/internal/models"
	"github.com/araza95/learn0x-GO/internal/utils/response"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var students types.Student

		err := json.NewDecoder(r.Body).Decode(&students)

		slog.Info("Creating a student.")

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Validate the request
		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}
