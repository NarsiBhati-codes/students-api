package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/NarsiBhati-codes/students-api/internal/types"
	response "github.com/NarsiBhati-codes/students-api/internal/utils"
	"github.com/go-playground/validator/v10"
)


func New() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

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

		response.WriteJson(res, http.StatusCreated, map[string]string {"success" : "ok"})
	}
}