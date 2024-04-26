package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"aitunews.kz/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions" // New import    fkmm flqo zpha zddd
	"golang.org/x/time/rate"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	services      *mysql.ServiceModel
	appointments  *mysql.AppointmentModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {

	// Sender's email credentials
	from := "daniyar.0586@gmail.com"
	password := "fkmm flqo zpha zddd"

	// Receiver's email address
	to := []string{"daniarm146@gmail.com"}

	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message
	msg := []byte("To: daniarm146@gmail.com\r\n" +
		"Subject: Test Email\r\n" +
		"\r\n" +
		"This is a test email sent from Golang.")

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Email sent successfully!")

	//Rate Limiting ------------------------------------------------
	// Create a limiter that allows 1 event per second.
	limiter := rate.NewLimiter(rate.Limit(1), 1)

	// Create a context for the limiter.
	ctx := context.Background()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for {
			// Wait for permission from the limiter.
			if err := limiter.Wait(ctx); err != nil {
				fmt.Println("Rate limit exceeded")
			}
			// Perform some task
			//fmt.Println("Task executed at:", time.Now().Format(time.StampMilli))
		}
	}()

	// Output "hello"
	fmt.Println("hello")

	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()

	// Create a multiwriter that writes to both the file and os.Stdout (terminal)
	multiWriter := io.MultiWriter(file, os.Stdout)

	log.SetOutput(multiWriter)

	log.Println("\n\n---------------------------------\n")

	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.SetOutput(multiWriter)
	errorLog.SetOutput(multiWriter)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	// Initialize a mysql.UserModel instance and add it to the application
	// dependencies.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		services:      &mysql.ServiceModel{DB: db},
		appointments:  &mysql.AppointmentModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)

	//--------------------------------------
	select {}
}

// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
