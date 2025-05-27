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

	var listasDashboard []models.ListaDashboardDiscente
	var exerciciosDashboard []models.ExercicioDashboardDiscente

	for indexExerc, exercicio := range exercicios {
		exercicioFoundDiscente := false

		for _, exercicioDiscente := range exerciciosDiscente {
			if exercicio.ID == exercicioDiscente.ExerciciosBaseID {
				exerciciosDashboard = append(exerciciosDashboard, models.ExercicioDashboardDiscente{
					Key: string(rune(indexExerc)),
					Data: models.ExercicioDataDiscente{
						Name:        exercicio.Titulo,
						Id:          exercicio.ID.String(),
						CodigoRodou: exercicioDiscente.CodigoRodou,
						Habilitado:  exercicioDiscente.Habilitado,
						CodigoBase:  exercicioDiscente.CodigoBase,
					},
				})
				exercicioFoundDiscente = true
			}
		}

		if !exercicioFoundDiscente {
			exerciciosDashboard = append(exerciciosDashboard, models.ExercicioDashboardDiscente{
				Key: string(rune(indexExerc)),
				Data: models.ExercicioDataDiscente{
					Name:        exercicio.Titulo,
					Id:          exercicio.ID.String(),
					CodigoRodou: false,
					Habilitado:  false,
					CodigoBase:  "",
				},
			})
		}
	}

	listasDashboard = append(listasDashboard, models.ListaDashboardDiscente{
		Name:       "-",
		Id:         uuidListaId.String(),
		Exercicios: exerciciosDashboard,
	})

	json.NewEncoder(w).Encode(listasDashboard)
}
