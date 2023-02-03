package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"sub-service/data"
	"sync"
	"testing"
	"time"
)

var testApp Config

func Testmain(m *testing.M) {

	gob.Register(data.User{})

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	testApp = Config{
		Session:       session,
		DB:            nil,
		InfoLog:       log.New(os.Stdout, "Info", log.Ldate|log.Ltime),
		ErrorLog:      log.New(os.Stdout, "Error", log.Ldate|log.Ltime|log.Lshortfile),
		Wait:          &sync.WaitGroup{},
		ErroChan:      make(chan error),
		ErrorChanDone: make(chan bool),
	}

	errorChan := make(chan error)
	mailerchan := make(chan Message, 100)
	mailerDoneChan := make(chan bool)

	testApp.Mailer = Mail{
		Wait:       testApp.Wait,
		ErrorChan:  errorChan,
		MailerChan: mailerchan,
		DoneChan:   mailerDoneChan,
	}

	go func() {
		select {
		case <-testApp.Mailer.MailerChan:
		case <-testApp.Mailer.ErrorChan:
		case <-testApp.Mailer.DoneChan:
			return

		}
	}()

	go func() {
		for {
			select {
			case err := <-testApp.ErroChan:
				testApp.ErrorLog.Println(err)
			case <-testApp.ErrorChanDone:
				return

			}
		}
	}()

	os.Exit(m.Run())
}
