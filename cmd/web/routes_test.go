package main

import (
	"fmt"
	"testing"

	"github.com/SIBUK/Bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing. Test passed
	default:
		t.Error(fmt.Sprintf("Type is no *chi.Mux, type is %T", v))
	}
}