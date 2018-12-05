package main

import (
	"database/sql"
	"final-project-zco/servers/gateway/handlers"
	"final-project-zco/servers/gateway/models/users"
	"final-project-zco/servers/gateway/sessions"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
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

	dbaddr := os.Getenv("DBADDR")
	if len(dbaddr) == 0 {
		dbaddr = "127.0.0.1:3306"
	}
	taskaddr := os.Getenv("TASKADDR")
	dsn := fmt.Sprintf("root:%s@tcp(%s)/userDB", mysqlPassWord, dbaddr)
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

	n := handlers.NewNotifier()
	ctx := &handlers.HandlerContext{
		SigningKey: sessionkey,
		Session:    redisStore,
		User:       store,
		Family:     store,
		Socket:     handlers.NewWebSocketsHandler(n),
		Request:    make(map[int64][]*users.User, 0),
	}
	rabbit := os.Getenv("RABBITADDR")

	conn, err := amqp.Dial("amqp://" + rabbit + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	request, err := ch.QueueDeclare(
		"taskQueue", // name matches what we used in our nodejs services
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")
	msgs, err := ch.Consume(
		request.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	fmt.Printf("Debug: after consume: %v", msgs)
	failOnError(err, "Failed to register a consumer")
	// go n.Start(msgs, request.Name, ctx)
	//
	// request, err := ch.QueueDeclare(
	// 	"authQueue", // name matches what we used in our go auth
	// 	true,        // durable
	// 	false,       // delete when unused
	// 	false,       // exclusive
	// 	false,       // no-wait
	// 	nil,         // arguments
	// )
	// failOnError(err, "Failed to declare a queue")

	// err = ch.Qos(
	// 	1,     // prefetch count
	// 	0,     // prefetch size
	// 	false, // global
	// )
	// failOnError(err, "Failed to set QoS")
	// msgs, err := ch.Consume(
	// 	request.Name, // queue
	// 	"",           // consumer
	// 	false,        // auto-ack
	// 	false,        // exclusive
	// 	false,        // no-local
	// 	false,        // no-wait
	// 	nil,          // args
	// )
	// go n.Start()
	mux := http.NewServeMux()
	mux.Handle("/tasks/", ctx.NewServiceProxy(taskaddr))

	mux.HandleFunc("/users", ctx.UsersHandler)
	mux.HandleFunc("/create", ctx.CreateHandler)
	mux.HandleFunc("/join", ctx.JoinHandler)
	mux.HandleFunc("/receive", ctx.ReceiveHandler)
	mux.HandleFunc("/accept", ctx.AcceptRequest)
	mux.HandleFunc("/users/", ctx.SpecificUserHandler)
	mux.HandleFunc("/sessions", ctx.SessionHandler)
	mux.HandleFunc("/sessions/", ctx.SpecificSessionHandler)
	mux.HandleFunc("/delete", ctx.DeleteHandler)
	mux.HandleFunc("/memberlist/", ctx.DisplayHandler)

	// mux.HandleFunc("/ws", ctx.WebSocketsHandler)
	wrappedMux := handlers.NewCors(mux)
	log.Printf("server is listening at %s...", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCertPath, tlsKeyPath, wrappedMux))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// func processMessages(ctx *handlers.MyHandler, msgs <-chan amqp.Delivery) {
// 	for msg := range msgs {
// 		// newMessage := msg.Body
// 		handler := &ctx.Chann
// 		err := json.Unmarshal(msg.Body, handler)
// 		if err != nil {
// 			log.Printf("error unmarshal %s", err)
// 			return
// 		}
// 		log.Printf("received message: %s", string(msg.Body))
// 		// var sharedChannel chan []byte
// 		// sharedChannel <- newMessage
// 		msg.Ack(false)
// 	}
// }
