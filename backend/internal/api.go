package internal

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/alfuveam/tcc/backend/config"
	"github.com/alfuveam/tcc/backend/controller"
	"github.com/alfuveam/tcc/backend/models"
)

type APIServer struct {
	addr string
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{
		addr: addr,
	}
}

func (s *APIServer) Run(db *sql.DB) error {
	router := http.NewServeMux()
	router.HandleFunc("POST /api/v1/login", func(w http.ResponseWriter, r *http.Request) {
		controller.OnLogin(w, r, db)
	})

	router.HandleFunc("POST /api/v1/registeraccount", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterUser(w, r, db)
	})

	loggedRouter := http.NewServeMux()
	loggedRouter.HandleFunc("GET /api/v1/dashboard", func(w http.ResponseWriter, r *http.Request) {
		controller.OnLoadDashBoard(w, r, db)
	})

	adminRouter := http.NewServeMux()
	// adminRouter.HandleFunc("POST /api/v1/func", func(w http.ResponseWriter, r *http.Request) {
	// 	controller.func(w, r, db)
	// })

	defaultMiddlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		CorsMiddleware,
	)

	authMiddlewareChain := MiddlewareChain(
		RequestLoggerMiddleware,
		CorsMiddleware,
		RequireAuthMiddleware,
	)

	router.Handle("/", authMiddlewareChain(loggedRouter))

	adminMiddlewareChain := MiddlewareChain(
		AdminRequereMiddleware,
	)

	loggedRouter.Handle("/", adminMiddlewareChain(adminRouter))

	server := http.Server{
		Addr:    s.addr,
		Handler: defaultMiddlewareChain(router),
		// Handler: middlewareChain(router),
		// Handler: RequireAuthMiddleware(RequestLoggerMiddleware(router)),
	}

	log.Printf("Server has started %s", s.addr)

	return server.ListenAndServe()
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middleware ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next.ServeHTTP
	}
}

func RequestLoggerMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method %s, path: %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func RequireAuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		//	Pull out the token
		encodedToken := strings.TrimPrefix(authorization, "Bearer ")

		user, _, err := controller.ValidateJWT(encodedToken)

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), config.MySigningKey, user)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	}
}

func CorsMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Se for uma requisição OPTIONS, responda com status 200 e termine a execução
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		//	only return json
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// Chame o próximo handler
		next.ServeHTTP(w, r)
	})
}

func AdminRequereMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(config.MySigningKey).(models.User)

		if !ok {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if user.UserType != 3 {
			log.Println("AdminRequereMiddleware error in elevated type: %v - %v", user.Id, user.UserType)
			http.Error(w, "Invalid user elevated type", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}
