package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	types "github.com/araza95/learn0x-GO/internal/models"
	"github.com/araza95/learn0x-GO/internal/storage"
	"github.com/araza95/learn0x-GO/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
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
		if err := validator.New().Struct(students); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateStudent(students.Name, students.Email, students.Age)

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
		}

		slog.Info("user created successfully")

		response.WriteJson(w, http.StatusCreated, map[string]int64{"ID": id})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var id = (r.PathValue("id"))

		fmt.Printf("this is query %s", id)

		studentId, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			response.WriteJson(w, http.StatusNotAcceptable, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById((studentId))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("user found")

		response.WriteJson(w, http.StatusOK, student)

	}

}
