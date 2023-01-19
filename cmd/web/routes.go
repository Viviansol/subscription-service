package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *Config) routes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)
	mux.Get("/", app.HomePage)
	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.LoginPage)
	mux.Get("/register", app.RegisterPage)
	mux.Post("/regist", app.PostRegisterPage)
	mux.Get("/activate", app.ActiveAccount)

	return mux

}
