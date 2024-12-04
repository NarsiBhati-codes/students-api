package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/NarsiBhati-codes/students-api/internal/storage"
	"github.com/NarsiBhati-codes/students-api/internal/types"
	response "github.com/NarsiBhati-codes/students-api/internal/utils"
	"github.com/go-playground/validator/v10"
)


func New(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		slog.Info("creating a student")

		var student types.Student 

		err := json.NewDecoder(req.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(
				res,
				http.StatusBadRequest,
				response.GeneralError(fmt.Errorf("empty body")),
			)
			return
		}

		if err != nil {
			response.WriteJson(res, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(res, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		slog.Info("user created successfully", slog.String("userId",fmt.Sprint(lastId)))


		if err != nil {
			response.WriteJson(res, http.StatusInternalServerError,err)
			return
		}

		response.WriteJson(res, http.StatusCreated, map[string]int64{"id": lastId})
	}
}