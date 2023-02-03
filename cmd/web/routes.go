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
	mux.Post("/register", app.PostRegisterPage)
	mux.Get("/activate", app.ActiveAccount)
	mux.Mount("/members", app.authRouter())

	return mux

}

func (app *Config) authRouter() http.Handler {
	mux := chi.NewRouter()
	mux.Use(app.Auth)
	mux.Get("/plans", app.ChooseSubscription)
	mux.Get("/subscribe", app.SubscribeToPlan)
	return mux
}
