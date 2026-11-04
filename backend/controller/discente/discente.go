package discente

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func DashboardDiscente(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	user, ok := r.Context().Value(config.MySigningKey).(models.UserJwt)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuário com ID inválido",
		})
		return
	}

	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	trilhas, err := q.GetAllTrilhas(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao buscar trilhas",
		})
		return
	}

	var trilhasDashboard []models.TrilhaDashboardDiscente

	for _, trilha := range trilhas {
		listas, err := q.GetListaByTrilhaId(context.Background(), trilha.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erro ao buscar listas",
			})
			return
		}

		var listasDashboard []models.ListaDashboardDiscente

		for _, lista := range listas {
			exercicios, err := q.GetExerciciosBaseByListaId(context.Background(), lista.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Erro ao buscar exercicios",
				})
				return
			}

			exerciciosDiscente, err := q.GetExerciciosDiscenteByListaId(context.Background(), generated.GetExerciciosDiscenteByListaIdParams{
				ListaID:         lista.ID,
				CreatedByUserID: uuidUser,
			})

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"error": "Erro ao buscar o progresso"})
				return
			}
			exerciciosDashboard := studentExercises(exercicios, exerciciosDiscente, false)

			listasDashboard = append(listasDashboard, models.ListaDashboardDiscente{
				Name:       lista.Name,
				Id:         lista.ID.String(),
				Exercicios: exerciciosDashboard,
			})
		}

		trilhasDashboard = append(trilhasDashboard, models.TrilhaDashboardDiscente{
			Name:            trilha.Name,
			Id:              trilha.ID.String(),
			TipoDaLinguagem: trilha.TipoDaLinguagem,

			Listas: listasDashboard,
		})
	}

	userDB, err := q.GetUserById(context.Background(), uuidUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao buscar usuário",
		})
		return
	}

	idTempoRepeticao := 0

	if userDB.RepeticaoEspacadaMinutos == (60) { // 1 Hora
		idTempoRepeticao = 1
	} else if userDB.RepeticaoEspacadaMinutos == (60 * 9) { // 9 Horas
		idTempoRepeticao = 2
	} else if userDB.RepeticaoEspacadaMinutos == (60 * 24) { // 1 Dia
		idTempoRepeticao = 3
	} else if userDB.RepeticaoEspacadaMinutos == (60 * 24 * 6) { // 6 Dias
		idTempoRepeticao = 4
	} else if userDB.RepeticaoEspacadaMinutos == (60 * 24 * 31) { // 31 Dias
		idTempoRepeticao = 5
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"trilhas":            trilhasDashboard,
		"id_tempo_repeticao": idTempoRepeticao,
	})
}
