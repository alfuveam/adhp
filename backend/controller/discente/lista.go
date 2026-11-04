package discente

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func ExerciciosHabilitadosByLista(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidListaId pgtype.UUID
	uuidListaId.Scan(r.PathValue("id"))

	user, ok := r.Context().Value(config.MySigningKey).(models.UserJwt)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuário com ID inválido",
		})
		return
	}

	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	exercicios, err := q.GetExerciciosBaseByListaId(context.Background(), uuidListaId)
	if err != nil {
		log.Println("GetExerciciosBaseByListaId: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao buscar exercícios da lista",
		})
		return
	}

	exerciciosDiscente, err := q.GetExerciciosDiscenteByListaId(context.Background(), generated.GetExerciciosDiscenteByListaIdParams{
		ListaID:         uuidListaId,
		CreatedByUserID: uuidUser,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao buscar o progresso"})
		return
	}
	var listasDashboard []models.ListaDashboardDiscente
	exerciciosDashboard := studentExercises(exercicios, exerciciosDiscente, true)

	listasDashboard = append(listasDashboard, models.ListaDashboardDiscente{
		Name:       "-",
		Id:         uuidListaId.String(),
		Exercicios: exerciciosDashboard,
	})

	json.NewEncoder(w).Encode(listasDashboard)
}
