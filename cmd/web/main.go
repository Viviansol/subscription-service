package main

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sub-service/data"
	"sync"
	"syscall"
	"time"
)

const webPort = "80"

func main() {

	db := initDB()
	db.Ping()

	session := initSession()

	infolog := log.New(os.Stdout, "Info", log.Ldate|log.Ltime)
	errlog := log.New(os.Stdout, "Error", log.Ldate|log.Ltime|log.Lshortfile)

	wg := sync.WaitGroup{}

	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infolog,
		ErrorLog: errlog,
		Wait:     &wg,
		Models:   data.New(db),
	}
	go app.listenForShutDown()
	app.serve()

}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	app.InfoLog.Println("Starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func initDB() *sql.DB {

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cant connect to DB")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0
	dsn := os.Getenv("DSN")
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not ready yet...")
			log.Println(err)

		} else {
			log.Println("connected to database")
			return connection
		}

		if counts > 10 {
			return nil
		}
		log.Println("Backing off for 1 second")
		time.Sleep(1 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSession() *scs.SessionManager {

	gob.Register(data.User{})

	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session

}

func initRedis() *redis.Pool {

	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS"))
		},
	}
	return redisPool

}

func (app *Config) listenForShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.shutDown()
	os.Exit(0)
}

func (app *Config) shutDown() {
	app.InfoLog.Println("would run cleanup tasks...")
	app.Wait.Wait()
	app.InfoLog.Println("closing channels and shutting down application...")

}
