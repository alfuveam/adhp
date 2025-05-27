package discente

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/controller/docente"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/alfuveam/adhp/backend/types"
	"github.com/jackc/pgx/v5/pgtype"
)

type ExerciciosRepeticao struct {
	ID                 string `json:"id"`
	DescricaoExercicio string `json:"descricao_exercicio"`
	ProximaRepeticao   string `json:"proxima_repeticao"`
	UltimaRepeticao    string `json:"ultima_repeticao"`
	IDExercicio        string `json:"id_exercicio"`
}

type SubmitExercicioRepeticao struct {
	RepeticaoEspacadaId pgtype.UUID `json:"repeticao_espacada_id"`
	ExercicioID         pgtype.UUID `json:"exercicio_id"`
	CodigoBase          string      `json:"codigo_base"`
}

type LoadRepeticaoEspacada struct {
	ExercicioID     string `json:"id"`
	ExercicioTitulo string `json:"name"`
}

func OnDiscenteSubmitTempoRepeticao(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	idTempoRepeticao, err := strconv.Atoi(string(r.PathValue("id")))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao selecionar o tempo de repetição",
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

	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	valorEmMinutos := 0

	if idTempoRepeticao == 1 { // 1 Hora
		valorEmMinutos = 60
	} else if idTempoRepeticao == 2 { // 9 Horas
		valorEmMinutos = 60 * 9
	} else if idTempoRepeticao == 3 { // 1 Dia
		valorEmMinutos = 60 * 24
	} else if idTempoRepeticao == 4 { // 6 Dias
		valorEmMinutos = 60 * 24 * 6
	} else if idTempoRepeticao == 5 { // 31 Dias
		valorEmMinutos = 60 * 24 * 31
	}

	userDB, _ := q.GetUserById(context.Background(), uuidUser)

	_, err = q.UpdateUserRepeticao(context.Background(), generated.UpdateUserRepeticaoParams{
		ID:                       userDB.ID,
		RepeticaoEspacadaMinutos: int64(valorEmMinutos),
	})

	w.WriteHeader(http.StatusOK)
}

func ExerciciosRepeticaoByUser(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
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

	repeticaoExercicios, err := q.GetExercRepByUserId(context.Background(), uuidUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao buscar exercícios de repetição",
		})
		return
	}

	var exerciciosRepeticao []ExerciciosRepeticao

	for _, exercicioRep := range repeticaoExercicios {
		//	Isso pode ocasionar problemas de performance
		exercicio, _ := q.GetExerciciosBaseById(context.Background(), exercicioRep.ExerciciosBaseID)

		exerciciosRepeticao = append(exerciciosRepeticao, ExerciciosRepeticao{
			ID:                 exercicioRep.ID.String(),
			DescricaoExercicio: exercicio.Titulo,
			UltimaRepeticao:    exercicioRep.UpdateAt.Time.Format(time.RFC822),
			ProximaRepeticao:   exercicioRep.ProximaRepeticao.Time.Format(time.RFC822),
			IDExercicio:        exercicioRep.ExerciciosBaseID.String(),
		})
	}

	json.NewEncoder(w).Encode(exerciciosRepeticao)
}

func GetExercicioRepeticaoEspacada(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidExercicioRep pgtype.UUID
	uuidExercicioRep.Scan(r.PathValue("id"))

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

	exercicioRep, err := q.GetExercRepByIdAndUserId(context.Background(), generated.GetExercRepByIdAndUserIdParams{
		ID:              uuidExercicioRep,
		CreatedByUserID: uuidUser,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Exercício de repetição não encontrado",
		})
		return
	}

	exercicio, err := q.GetExerciciosBaseById(context.Background(), exercicioRep.ExerciciosBaseID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Exercício base nao encontrado",
		})
		return
	}

	lista, err := q.GetListaById(context.Background(), exercicio.ListaID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Lista nao encontrada",
		})
		return
	}

	trilha, err := q.GetTrilhaById(context.Background(), lista.TrilhaID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Trilha nao encontrada",
		})
		return
	}

	exercicioRepeticao := LoadRepeticaoEspacada{
		ExercicioID:     exercicio.ID.String(),
		ExercicioTitulo: exercicio.Titulo,
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"exercicio":         exercicioRepeticao,
		"tipo_da_linguagem": trilha.TipoDaLinguagem,
	})
}

func OnDiscenteSubmitExercicioRepeticaoEspacada(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req SubmitExercicioRepeticao
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

	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	exercicio, err := q.GetExerciciosBaseById(context.Background(), req.ExercicioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Exercício base nao encontrado",
		})
		return
	}

	unitTHandler := docente.UnitTestHandler{
		SourceFromUser:  req.CodigoBase,
		SourceUnitTeste: exercicio.CodigoTeste,
		Lista:           exercicio.ListaID.String(),
		Exercicio:       exercicio.ID.String(),
		Usuario:         user.Id,
	}

	lista, err := q.GetListaById(context.Background(), exercicio.ListaID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Lista nao encontrada",
		})
		return
	}

	trilha, err := q.GetTrilhaById(context.Background(), lista.TrilhaID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Trilha nao encontrada",
		})
		return
	}

	codeHandler := docente.CodeHandlerMetricas{
		CreatedByUserId:  uuidUser,
		TrilhaId:         trilha.ID,
		ListaId:          lista.ID,
		ExerciciosBaseId: exercicio.ID,
		LinguagemTipo:    trilha.TipoDaLinguagem,
	}

	unitTestRep := docente.OnRunTestCode(unitTHandler, codeHandler, q)

	userDB, err := q.GetUserById(context.Background(), uuidUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Usuário nao encontrado",
		})
		return
	}

	if !unitTestRep.OutPutFromUser.Success && !unitTestRep.OutPutUnitTeste.Success {
		err = HandlerMetricasRepeticao(userDB, exercicio.ID, types.MetricasTentou, q)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(unitTestRep)
		return
	}

	//	insere data de repetição
	dataRepeticao := pgtype.Timestamp{Time: time.Now().Add(time.Minute * time.Duration(userDB.RepeticaoEspacadaMinutos)), Valid: true}

	q.UpdateExercRepRepeticao(context.Background(), generated.UpdateExercRepRepeticaoParams{
		CreatedByUserID:  uuidUser,
		ExerciciosBaseID: exercicio.ID,
		ProximaRepeticao: dataRepeticao,
		UpdateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	err = HandlerMetricasRepeticao(userDB, exercicio.ID, types.MetricasRodou, q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(unitTestRep)
}
