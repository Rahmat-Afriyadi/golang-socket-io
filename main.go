// package main

// import (
// 	"log"
// 	"net/http"

// 	socketio "github.com/googollee/go-socket.io"
// 	// "github.com/gorilla/mux"
// 	// "github.com/rs/cors"
// )

// func SayHelloWorld(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte("Hello, World!"))
// }

// // Easier to get running with CORS. Thanks for help @Vindexus and @erkie

// type CrossOriginServer struct{}

// func (s *CrossOriginServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
// 	log.Println(req.Header.Get("Origin"))
// 	log.Println(req.Method)

// 	if origin := req.Header.Get("Origin"); origin != "" {
// 		rw.Header().Set("Access-Control-Allow-Origin", origin)
// 		rw.Header().Set("Access-Control-Allow-Credentials", "true")
// 		rw.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE, OPTIONS")
// 		rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
// 	}
// 	if req.Method == "OPTIONS" {
// 		return
// 	}

// 	server := socketio.NewServer(nil)

// 	server.OnConnect("/", func(so socketio.Conn) error {
// 		so.SetContext("")
// 		log.Println("connected:", so.ID())
// 		return nil
// 	})

// 	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
// 		log.Println("notice:", msg)
// 		s.Emit("reply", "have "+msg)
// 	})

// 	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
// 		s.SetContext(msg)
// 		return "recv " + msg
// 	})

// 	server.OnEvent("/", "bye", func(s socketio.Conn) string {
// 		last := s.Context().(string)
// 		s.Emit("bye", last)
// 		s.Close()
// 		return last
// 	})

// 	server.OnError("/", func(s socketio.Conn, e error) {
// 		log.Println("meet error:", e)
// 	})

// 	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
// 		log.Println("closed", reason)
// 	})

// 	go func() {
// 		if err := server.Serve(); err != nil {
// 			log.Fatalf("socketio listen error: %s\n", err)
// 		}
// 	}()
// 	defer server.Close()

// 	mux := http.NewServeMux()
// 	mux.Handle("/socket.io", server)
// 	// mux.HandleFunc("/", SayHelloWorld)

// 	mux.ServeHTTP(rw, req)
// }

// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

// 		w.Header().Set("Content-Type", "application/json")
// 		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
// 		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true")
// 		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

// 		next.ServeHTTP(w, r)
// 	})
// }

// func main() {

// server := socketio.NewServer(nil)
// server := socketio.NewServer(&engineio.Options{
// 	Transports: []transport.Transport{
// 		&polling.Transport{
// 			CheckOrigin: allowOriginFunc,
// 		},
// 		&websocket.Transport{
// 			CheckOrigin: allowOriginFunc,
// 		},
// 	},
// })

// 	server.OnConnect("/", func(s socketio.Conn) error {
// 		s.SetContext("")
// 		log.Println("connected:", s.ID())
// 		return nil
// 	})

// 	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
// 		log.Println("notice:", msg)
// 		s.Emit("reply", "have "+msg)
// 	})

// 	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
// 		s.SetContext(msg)
// 		return "recv " + msg
// 	})

// 	server.OnEvent("/", "bye", func(s socketio.Conn) string {
// 		last := s.Context().(string)
// 		s.Emit("bye", last)
// 		s.Close()
// 		return last
// 	})

// 	server.OnError("/", func(s socketio.Conn, e error) {
// 		log.Println("meet error:", e)
// 	})

// 	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
// 		log.Println("closed", reason)
// 	})

// 	go func() {
// 		if err := server.Serve(); err != nil {
// 			log.Fatalf("socketio listen error: %s\n", err)
// 		}
// 	}()
// 	defer server.Close()

// 	// c := cors.New(cors.Options{
// 	// 	AllowedOrigins:   []string{"https://syafiq-ina.vercel.app", "http://localhost:3000"},
// 	// 	AllowCredentials: true,
// 	// })

// 	// router := mux.NewRouter()

// 	// router.HandleFunc("/socket.io/", server)
// 	// router.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./assets/")))

// 	// handler := c.Handler(router)

// 	http.Handle("/socket.io/", corsMiddleware(server))
// 	http.Handle("/", corsMiddleware(http.FileServer(http.Dir("./asset"))))

// 	log.Println("Serving at localhost:8000...")
// 	log.Fatal(http.ListenAndServe(":8000", nil))

// 	// log.Println("listening on :8080")
// 	// http.ListenAndServe(":8080", &CrossOriginServer{})
// }

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// "github.com/graarh/golang-socketio/transport"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.New()

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		fmt.Println("conected")
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		fmt.Println(msg)
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		log.Println("closed", msg)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	router.Use(GinMiddleware("http://localhost:8080"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("./asset"))

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal("failed run app: ", err)
	}
}
