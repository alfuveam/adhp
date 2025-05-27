package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/controller"
	"github.com/alfuveam/adhp/backend/controller/discente"
	"github.com/alfuveam/adhp/backend/controller/docente"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
)

type APIServer struct {
	addr    string
	db      *sql.DB
	queries *generated.Queries
}

func NewAPIServer(addr string, queries *generated.Queries) *APIServer {
	return &APIServer{
		addr:    addr,
		queries: queries,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	router.HandleFunc("POST /v1/login", func(w http.ResponseWriter, r *http.Request) {
		controller.OnLogin(w, r, s.queries)
	})

	router.HandleFunc("POST /v1/registeraccount", func(w http.ResponseWriter, r *http.Request) {
		controller.RegisterUser(w, r, s.queries)
	})

	loggedRouter := http.NewServeMux()
	loggedRouter.HandleFunc("POST /v1/logout", func(w http.ResponseWriter, r *http.Request) {
		controller.OnLogout(w, r, s.db)
	})

	loggedRouter.HandleFunc("GET /v1/dashboard_discente", func(w http.ResponseWriter, r *http.Request) {
		discente.DashboardDiscente(w, r, s.queries)
	})

	loggedRouter.HandleFunc("GET /v1/exercicios_habilitados_by_lista/{id}", func(w http.ResponseWriter, r *http.Request) {
		discente.ExerciciosHabilitadosByLista(w, r, s.queries)
	})

	loggedRouter.HandleFunc("GET /v1/feedback_by_exercicio_id/{id}/{tipo_feedback}", func(w http.ResponseWriter, r *http.Request) {
		discente.GetFeedbackByExercicioBaseId(w, r, s.queries)
	})

	loggedRouter.HandleFunc("POST /v1/discente_submit_exercicio", func(w http.ResponseWriter, r *http.Request) {
		discente.OnDiscenteSubmitExercicio(w, r, s.queries)
	})

	// loggedRouter.HandleFunc("PUT /v1/discente_submit_tempo_repeticao/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	discente.OnDiscenteSubmitTempoRepeticao(w, r, s.queries)
	// })

	loggedRouter.HandleFunc("GET /v1/exercicios_repeticao_by_user", func(w http.ResponseWriter, r *http.Request) {
		discente.ExerciciosRepeticaoByUser(w, r, s.queries)
	})

	loggedRouter.HandleFunc("GET /v1/discente_get_exerc_rep_espacada/{id}", func(w http.ResponseWriter, r *http.Request) {
		discente.GetExercicioRepeticaoEspacada(w, r, s.queries)
	})

	loggedRouter.HandleFunc("POST /v1/discente_submit_exerc_rep_espacada", func(w http.ResponseWriter, r *http.Request) {
		discente.OnDiscenteSubmitExercicioRepeticaoEspacada(w, r, s.queries)
	})

	loggedRouter.HandleFunc("POST /v1/metricas_exercicio", func(w http.ResponseWriter, r *http.Request) {
		discente.MetricasInicioExercicio(w, r, s.queries)
	})

	loggedRouter.HandleFunc("POST /v1/metricas_repeticao_espacada", func(w http.ResponseWriter, r *http.Request) {
		discente.MetricasInicioRepeticao(w, r, s.queries)
	})

	//	Docente
	docenteRouter := http.NewServeMux()
	docenteRouter.HandleFunc("POST /v1/add_trilha", func(w http.ResponseWriter, r *http.Request) {
		docente.AddTrilha(w, r, s.queries)
	})

	docenteRouter.HandleFunc("GET /v1/trilhas_lista_exercicios", func(w http.ResponseWriter, r *http.Request) {
		docente.GetTrilhasListasExercicios(w, r, s.queries)
	})

	docenteRouter.HandleFunc("PUT /v1/update_trilha", func(w http.ResponseWriter, r *http.Request) {
		docente.UpdateTrilha(w, r, s.queries)
	})

	docenteRouter.HandleFunc("GET /v1/get_trilha/{id}", func(w http.ResponseWriter, r *http.Request) {
		docente.GetTrilhaById(w, r, s.queries)
	})

	docenteRouter.HandleFunc("DELETE /v1/remover_trilha/{id}", func(w http.ResponseWriter, r *http.Request) {
		docente.RemoveTrilha(w, r, s.queries)
	})

	docenteRouter.HandleFunc("POST /v1/add_lista", func(w http.ResponseWriter, r *http.Request) {
		docente.AddLista(w, r, s.queries)
	})

	docenteRouter.HandleFunc("PUT /v1/update_lista", func(w http.ResponseWriter, r *http.Request) {
		docente.UpdateLista(w, r, s.queries)
	})

	docenteRouter.HandleFunc("DELETE /v1/remove_lista/{id}", func(w http.ResponseWriter, r *http.Request) {
		docente.RemoveLista(w, r, s.queries)
	})

	docenteRouter.HandleFunc("PUT /v1/update_lista_index", func(w http.ResponseWriter, r *http.Request) {
		docente.UpdateListaIndex(w, r, s.queries)
	})

	docenteRouter.HandleFunc("POST /v1/adicionar_exericicio", func(w http.ResponseWriter, r *http.Request) {
		docente.AddExercicio(w, r, s.queries)
	})

	docenteRouter.HandleFunc("DELETE /v1/remover_exercicio/{id}", func(w http.ResponseWriter, r *http.Request) {
		docente.RemoverExercicio(w, r, s.queries)
	})

	docenteRouter.HandleFunc("POST /v1/atualizar_exericicio", func(w http.ResponseWriter, r *http.Request) {
		docente.AtualizarExercicio(w, r, s.queries)
	})

	docenteRouter.HandleFunc("GET /v1/exercicio/{id}", func(w http.ResponseWriter, r *http.Request) {
		docente.GetExercicio(w, r, s.queries)
	})

	docenteRouter.HandleFunc("PUT /v1/update_exercicio_index", func(w http.ResponseWriter, r *http.Request) {
		docente.UpdateExercicioIndex(w, r, s.queries)
	})

	docenteRouter.HandleFunc("DELETE /v1/remove_feedback/{id}", func(w http.ResponseWriter, r *http.Request) {
		docente.RemoveFeedback(w, r, s.queries)
	})

	// adminRouter := http.NewServeMux()
	// adminRouter.HandleFunc("POST /v1/func", func(w http.ResponseWriter, r *http.Request) {
	// 	// controller.func(w, r, db)
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

	docenteMiddlewareChain := MiddlewareChain(
		DocenteRequereMiddleware,
	)

	loggedRouter.Handle("/", docenteMiddlewareChain(docenteRouter))

	// adminMiddlewareChain := MiddlewareChain(
	// 	AdminRequereMiddleware,
	// )

	// docenteRouter.Handle("/", adminMiddlewareChain(adminRouter))

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
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
		user, ok := r.Context().Value(config.MySigningKey).(models.UserJwt)

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Usuário com ID inválido",
			})
			return
		}

		if user.UserType != 3 {
			log.Println("AdminRequereMiddleware error in elevated type: %v - %v", user.Id, user.UserType)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid user elevated type",
			})
			return
		}

		next.ServeHTTP(w, r)
	}
}

func DocenteRequereMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(config.MySigningKey).(models.UserJwt)

		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Usuário com ID inválido",
			})
			return
		}

		// ??
		if (user.UserType < 2) == true {
			log.Println("DocenteRequereMiddleware error in elevated type: %v - %v", user.Id, user.UserType)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid user elevated type",
			})
			return
		}

		next.ServeHTTP(w, r)
	}
}
