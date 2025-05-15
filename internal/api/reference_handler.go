package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
)

type ReferenceHandler struct {
	logger *log.Logger
}

func NewReferenceHandler() *ReferenceHandler {
	return &ReferenceHandler{}
}

func (h *ReferenceHandler) HandleGetReference(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./internal/doc/openapi.json",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "REST API GO",
		},
		DarkMode: true,
		Theme:    scalar.ThemeKepler,
	})

	if err != nil {
		fmt.Printf("%v", err)
	}

	fmt.Fprintln(w, htmlContent)
}
