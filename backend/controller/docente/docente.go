package docente

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/jackc/pgx/v5/pgtype"
)

func AddTrilha(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req generated.Trilha
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Requisição com payload inválido",
		})
		return
	}

	user, ok := r.Context().Value(config.MySigningKey).(models.UserJwt)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuário com ID inválido",
		})
		return
	}

	//	convert string to uuid []byte
	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	_, err = q.CreateTrilha(context.Background(), generated.CreateTrilhaParams{
		Name:            req.Name,
		CreatedByUserID: uuidUser,
		CreateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateTrilha(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req generated.UpdateTrilhaParams
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Requisição com payload inválido",
		})
		return
	}

	trilha, _ := q.GetTrilhaById(context.Background(), req.ID)

	trilha.Name = req.Name
	trilha.TipoDaLinguagem = req.TipoDaLinguagem

	q.UpdateTrilha(context.Background(), generated.UpdateTrilhaParams{
		ID:              trilha.ID,
		CreatedByUserID: trilha.CreatedByUserID,
		Name:            trilha.Name,
		Description:     trilha.Description,
		TipoDaLinguagem: trilha.TipoDaLinguagem,
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	w.WriteHeader(http.StatusOK)
}

func RemoveTrilha(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidTrilha pgtype.UUID
	uuidTrilha.Scan(r.PathValue("id"))

	_, err := q.DeleteTrilha(context.Background(), uuidTrilha)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetTrilhasListasExercicios(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	trilhas, err := q.GetAllTrilhas(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	var trilhasDashboard []models.TrilhaDashboard

	for _, trilha := range trilhas {
		listas, err := q.GetListaByTrilhaId(context.Background(), trilha.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		var listasDashboard []models.ListaDashboard

		for _, lista := range listas {
			exercicios, err := q.GetExerciciosBaseByListaId(context.Background(), lista.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			var exerciciosDashboard []models.ExercicioDashboard

			for indexExerc, exercicio := range exercicios {
				exerciciosDashboard = append(exerciciosDashboard, models.ExercicioDashboard{
					Key: string(rune(indexExerc)),
					Data: models.ExercicioData{
						Name: exercicio.Titulo,
						Id:   exercicio.ID.String(),
					},
				})
			}

			listasDashboard = append(listasDashboard, models.ListaDashboard{
				Name:       lista.Name,
				Id:         lista.ID.String(),
				Exercicios: exerciciosDashboard,
			})
		}

		trilhasDashboard = append(trilhasDashboard, models.TrilhaDashboard{
			Name:            trilha.Name,
			Id:              trilha.ID.String(),
			TipoDaLinguagem: trilha.TipoDaLinguagem,

			Listas: listasDashboard,
		})
	}
	json.NewEncoder(w).Encode(trilhasDashboard)
}

func GetTrilhaById(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidTrilha pgtype.UUID
	uuidTrilha.Scan(r.PathValue("id"))

	trilha, err := q.GetTrilhaById(context.Background(), uuidTrilha)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(trilha)
}
