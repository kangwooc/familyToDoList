package main

import (
	"database/sql"
	"fmt"
	"final-project-zco/servers/gateway/handlers"
	"homework-juan3674-1532739/servers/gateway/models/users"
	"homework-juan3674-1532739/servers/gateway/sessions"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

//main is the main entry point for the server
func main() {
	addr := os.Getenv("ADDR")
	if len(addr) == 0 {
		addr = ":443"
	}

	tlsKeyPath := os.Getenv("TLSKEY")
	if len(tlsKeyPath) == 0 {
		os.Stdout.WriteString("tlskey is not found")
		os.Exit(1)
	}
	tlsCertPath := os.Getenv("TLSCERT")
	if len(tlsCertPath) == 0 {
		os.Stdout.WriteString("tlscert is not found")
		os.Exit(1)
	}
	sessionkey := os.Getenv("SESSIONKEY")
	if len(sessionkey) == 0 {
		sessionkey = "default"
	}
	redisAddr := os.Getenv("REDISADDR")
	if len(redisAddr) == 0 {
		redisAddr = "127.0.0.1:6379"
	}
	
	mysqlPassWord := os.Getenv("MYSQL_ROOT_PASSWORD")
	if len(mysqlPassWord) == 0 {
		os.Stdout.WriteString("mysqlPassWord is not found")
		os.Exit(1)
	}
	dsn := fmt.Sprintf("root:%s@tcp(mysqlserver:3306)/userdb", mysqlPassWord)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	redisStore := sessions.NewRedisStore(client, time.Hour)
	store := users.NewMySQLStore(db)
	ctx := &handlers.HandlerContext{
		SigningKey: sessionkey,
		Session:    redisStore,
		User:       store,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", ctx.UsersHandler)
	mux.HandleFunc("/v1/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/v1/sessions", ctx.SessionHandler)
	mux.HandleFunc("/v1/sessions/", ctx.SpecificSessionHandler)
	wrappedMux := handlers.NewCors(mux)

	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}