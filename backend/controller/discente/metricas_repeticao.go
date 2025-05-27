package discente

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type RepeticaoMetricas struct {
	RepeticaoID pgtype.UUID `json:"exercicio_id"`
	TipoMetrica int16       `json:"tipo_metrica"`
}

func HandlerMetricasRepeticao(user generated.User, exercicioID pgtype.UUID, tipoMetrica int16, q *generated.Queries) error {
	exercicio, err := q.GetExerciciosBaseById(context.Background(), exercicioID)
	if err != nil {
		return errors.New("Exercício base nao encontrado")
	}

	lista, err := q.GetListaById(context.Background(), exercicio.ListaID)
	if err != nil {
		return errors.New("Lista nao encontrada")
	}

	trilha, err := q.GetTrilhaById(context.Background(), lista.TrilhaID)
	if err != nil {
		return errors.New("Trilha nao encontrada")
	}

	_, err = q.CreateDiscenteMetricasRepeticao(context.Background(), generated.CreateDiscenteMetricasRepeticaoParams{
		CreatedByUserID:          user.ID,
		TrilhaID:                 trilha.ID,
		ListaID:                  exercicio.ListaID,
		ExerciciosBaseID:         exercicio.ID,
		HorarioAt:                pgtype.Timestamp{Time: time.Now(), Valid: true},
		Tipo:                     tipoMetrica,
		RepeticaoEspacadaMinutos: user.RepeticaoEspacadaMinutos,
	})

	if err != nil {
		return errors.New("Erro ao criar metrica")
	}

	return nil
}

func MetricasInicioRepeticao(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var repeticaoMetricas RepeticaoMetricas
	err := json.NewDecoder(r.Body).Decode(&repeticaoMetricas)
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

	var userUUID pgtype.UUID
	userUUID.Scan(user.Id)

	userDB, err := q.GetUserById(context.Background(), userUUID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuário nao encontrado",
		})
		return
	}

	err = HandlerMetricasRepeticao(userDB, repeticaoMetricas.RepeticaoID, repeticaoMetricas.TipoMetrica, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
}
