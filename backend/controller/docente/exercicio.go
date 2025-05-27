package docente

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/alfuveam/adhp/backend/types"
	"github.com/jackc/pgx/v5/pgtype"
)

type FeedbackHandler struct {
	Descricao string `json:"descricao"`
}

type ExercicioAdd struct {
	Titulo          string            `json:"titulo"`
	ListaId         string            `json:"lista_id"`
	OrderIndex      int               `json:"order_index"`
	CodigoBase      string            `json:"codigo_base"`
	CodigoTeste     string            `json:"codigo_teste"`
	FeedbackHandler []FeedbackHandler `json:"feedbacks"`
}

type FeedbackManipulador struct {
	Id        string `json:"id"`
	Descricao string `json:"descricao"`
}

type ExercicioManipulador struct {
	Id              pgtype.UUID           `json:"id"`
	Titulo          string                `json:"titulo"`
	ListaId         string                `json:"lista_id"`
	OrderIndex      int                   `json:"order_index"`
	CodigoBase      string                `json:"codigo_base"`
	CodigoTeste     string                `json:"codigo_teste"`
	FeedbackHandler []FeedbackManipulador `json:"feedbacks"`
}

type UnitTestHandler struct {
	SourceFromUser  string `json:"source_from_user"`
	SourceUnitTeste string `json:"source_unit_teste"`
	Lista           string `json:"lista"`
	Exercicio       string `json:"exercicio"`
	Usuario         string `json:"usuario"`
}

type ResponseData struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

type UnitTestResponse struct {
	OutPutFromUser  ResponseData `json:"out_put_from_user"`
	OutPutUnitTeste ResponseData `json:"out_put_unit_teste"`
}

type ExercicioIndex struct {
	ListaID         pgtype.UUID `json:"lista_id"`
	ExercicioID     pgtype.UUID `json:"exercicio_id"`
	PosicoesATrocar int16       `json:"posicoes_a_trocar"`
}

type CodeHandlerMetricas struct {
	CreatedByUserId  pgtype.UUID `json:"created_by_user_id"`
	TrilhaId         pgtype.UUID `json:"trilha_id"`
	ListaId          pgtype.UUID `json:"lista_id"`
	ExerciciosBaseId pgtype.UUID `json:"exercicios_base_id"`
	LinguagemTipo    int16       `json:"linguagem_tipo"`
}

type DocenteController interface {
	OnRunTestCode(unitTHandler UnitTestHandler, codeHandler CodeHandlerMetricas, q *generated.Queries) UnitTestResponse
}

func OnRunTestCode(unitTHandler UnitTestHandler, codeHandler CodeHandlerMetricas, q *generated.Queries) UnitTestResponse {
	q.CreateCodeHandlerMetricas(context.Background(), generated.CreateCodeHandlerMetricasParams{
		CreatedByUserID:  codeHandler.CreatedByUserId,
		TrilhaID:         codeHandler.TrilhaId,
		ListaID:          codeHandler.ListaId,
		ExerciciosBaseID: codeHandler.ExerciciosBaseId,
		HorarioAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
		Tipo:             types.MetricasInicio,
		Linguagem:        codeHandler.LinguagemTipo,
	})

	lan := "golang"
	if codeHandler.LinguagemTipo == 2 {
		lan = "python"
	}

	reqBodyCodeHandler, _ := json.Marshal(unitTHandler)
	reqCodeHandler, _ := http.NewRequest("POST", config.CodeHandlerApiUrl+"/api/v1/run-test-"+lan, bytes.NewBuffer(reqBodyCodeHandler))
	reqCodeHandler.Header.Add("accept", "application/json")
	reqCodeHandler.Header.Add("Content-Type", "application/json")
	reqCodeHandler.Header.Add("X-Auth-Token", config.CodeHandlerKey)

	resCodeHandler, _ := http.DefaultClient.Do(reqCodeHandler)

	body, _ := io.ReadAll(resCodeHandler.Body)
	defer resCodeHandler.Body.Close()

	var testResponse UnitTestResponse
	err := json.Unmarshal(body, &testResponse)
	if err != nil {
		q.CreateCodeHandlerMetricas(context.Background(), generated.CreateCodeHandlerMetricasParams{
			CreatedByUserID:  codeHandler.CreatedByUserId,
			TrilhaID:         codeHandler.TrilhaId,
			ListaID:          codeHandler.ListaId,
			ExerciciosBaseID: codeHandler.ExerciciosBaseId,
			HorarioAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
			Tipo:             types.MetricasTentou,
			Linguagem:        codeHandler.LinguagemTipo,
		})

		return UnitTestResponse{}
	}

	if resCodeHandler.StatusCode != 200 {
		q.CreateCodeHandlerMetricas(context.Background(), generated.CreateCodeHandlerMetricasParams{
			CreatedByUserID:  codeHandler.CreatedByUserId,
			TrilhaID:         codeHandler.TrilhaId,
			ListaID:          codeHandler.ListaId,
			ExerciciosBaseID: codeHandler.ExerciciosBaseId,
			HorarioAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
			Tipo:             types.MetricasTentou,
			Linguagem:        codeHandler.LinguagemTipo,
		})

		return testResponse
	}

	q.CreateCodeHandlerMetricas(context.Background(), generated.CreateCodeHandlerMetricasParams{
		CreatedByUserID:  codeHandler.CreatedByUserId,
		TrilhaID:         codeHandler.TrilhaId,
		ListaID:          codeHandler.ListaId,
		ExerciciosBaseID: codeHandler.ExerciciosBaseId,
		HorarioAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
		Tipo:             types.MetricasRodou,
		Linguagem:        codeHandler.LinguagemTipo,
	})

	return testResponse
}

func AddExercicio(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req ExercicioAdd
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

	unitTestHandler := UnitTestHandler{
		SourceFromUser:  req.CodigoBase,
		SourceUnitTeste: req.CodigoTeste,
		Lista:           req.ListaId,
		Exercicio:       "exercicio_novo", //	warning
		Usuario:         user.Id,
	}

	var uuidLista pgtype.UUID
	uuidLista.Scan(req.ListaId)

	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	lista, _ := q.GetListaById(context.Background(), uuidLista)
	trilha, _ := q.GetTrilhaById(context.Background(), lista.TrilhaID)

	codeHandler := CodeHandlerMetricas{
		CreatedByUserId:  uuidUser,
		TrilhaId:         trilha.ID,
		ListaId:          uuidLista,
		ExerciciosBaseId: pgtype.UUID{},
		LinguagemTipo:    trilha.TipoDaLinguagem,
	}

	testResponse := OnRunTestCode(unitTestHandler, codeHandler, q)

	if !testResponse.OutPutFromUser.Success && !testResponse.OutPutUnitTeste.Success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(testResponse)
		return
	}

	orderIndex, _ := q.GetExerciciosBaseCountByListaId(context.Background(), uuidLista)
	orderIndex = orderIndex + 1

	exercicio, err := q.CreateExerciciosBase(context.Background(), generated.CreateExerciciosBaseParams{
		CreatedByUserID: uuidUser,
		ListaID:         uuidLista,
		OrderIndex:      orderIndex,
		Titulo:          req.Titulo,
		CodigoBase:      req.CodigoBase,
		CodigoTeste:     req.CodigoTeste,
		CreateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "resposta com payload inválido",
		})
		return
	}

	for _, feedback := range req.FeedbackHandler {
		_, err = q.CreateFeedback(context.Background(), generated.CreateFeedbackParams{
			CreatedByUserID:  uuidUser,
			ExerciciosBaseID: exercicio.ID,
			Descricao:        feedback.Descricao,
			CreateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "resposta com payload inválido",
			})
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testResponse)
}

func RemoverExercicio(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidExercicioId pgtype.UUID
	uuidExercicioId.Scan(r.PathValue("id"))

	_, err := q.DeleteExerciciosBase(context.Background(), uuidExercicioId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func AtualizarExercicio(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req ExercicioManipulador
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

	exercicio, _ := q.GetExerciciosBaseById(context.Background(), req.Id)

	unitTestHandler := UnitTestHandler{
		SourceFromUser:  req.CodigoBase,
		SourceUnitTeste: req.CodigoTeste,
		Lista:           exercicio.ListaID.String(),
		Exercicio:       exercicio.ID.String(),
		Usuario:         user.Id,
	}

	lista, _ := q.GetListaById(context.Background(), exercicio.ListaID)
	trilha, _ := q.GetTrilhaById(context.Background(), lista.TrilhaID)

	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	codeHandler := CodeHandlerMetricas{
		CreatedByUserId:  uuidUser,
		TrilhaId:         trilha.ID,
		ListaId:          lista.ID,
		ExerciciosBaseId: exercicio.ID,
		LinguagemTipo:    trilha.TipoDaLinguagem,
	}

	testResponse := OnRunTestCode(unitTestHandler, codeHandler, q)

	if !testResponse.OutPutFromUser.Success && !testResponse.OutPutUnitTeste.Success {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(testResponse)
		return
	}

	_, err = q.UpdateExerciciosBase(context.Background(), generated.UpdateExerciciosBaseParams{
		ID:              exercicio.ID,
		CreatedByUserID: exercicio.CreatedByUserID,
		ListaID:         exercicio.ListaID,
		OrderIndex:      exercicio.OrderIndex,
		Titulo:          req.Titulo,
		CodigoBase:      req.CodigoBase,
		CodigoTeste:     req.CodigoTeste,
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "resposta com payload inválido",
		})
		return
	}

	for _, feedback := range req.FeedbackHandler {
		if feedback.Id == "0" {
			_, err = q.CreateFeedback(context.Background(), generated.CreateFeedbackParams{
				CreatedByUserID:  uuidUser,
				ExerciciosBaseID: exercicio.ID,
				Descricao:        feedback.Descricao,
				CreateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
			})
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "resposta com payload inválido",
				})
				return
			}
			continue
		}
		var uuidFeedback pgtype.UUID
		uuidFeedback.Scan(feedback.Id)

		_, err = q.UpdateFeedbackDescricaoById(context.Background(), generated.UpdateFeedbackDescricaoByIdParams{
			ID:        uuidFeedback,
			Descricao: feedback.Descricao,
			UpdateAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
		})

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "resposta com payload inválido",
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testResponse)
}

func GetExercicio(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidExercicioId pgtype.UUID
	uuidExercicioId.Scan(r.PathValue("id"))

	exercicio, err := q.GetExerciciosBaseById(context.Background(), uuidExercicioId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	feedback, _ := q.GetFeedbackByExerciciosBaseId(context.Background(), uuidExercicioId)

	var feedbackManipulador []FeedbackManipulador
	for _, item := range feedback {
		feedbackAtualizarR := FeedbackManipulador{
			Id:        item.ID.String(),
			Descricao: item.Descricao,
		}

		feedbackManipulador = append(feedbackManipulador, feedbackAtualizarR)
	}

	exercicioManipulador := ExercicioManipulador{
		Id:              exercicio.ID,
		Titulo:          exercicio.Titulo,
		CodigoBase:      exercicio.CodigoBase,
		CodigoTeste:     exercicio.CodigoTeste,
		FeedbackHandler: feedbackManipulador,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exercicioManipulador)
}

func UpdateExercicioIndex(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req ExercicioIndex
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Requisição com payload inválido",
		})
		return
	}

	exercicio, _ := q.GetExerciciosBaseById(context.Background(), req.ExercicioID)

	//	troca pela posicao do outro exercicio
	posicaoAtrocarOutroExercicio := exercicio.OrderIndex + req.PosicoesATrocar

	outroExercicio, err := q.GetExercicioByOrderIndexAndListaId(context.Background(), generated.GetExercicioByOrderIndexAndListaIdParams{
		OrderIndex: posicaoAtrocarOutroExercicio,
		ListaID:    req.ListaID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	outroExercicio.OrderIndex = exercicio.OrderIndex
	q.UpdateExerciciosBase(context.Background(), generated.UpdateExerciciosBaseParams{
		ID:              outroExercicio.ID,
		CreatedByUserID: outroExercicio.CreatedByUserID,
		ListaID:         outroExercicio.ListaID,
		OrderIndex:      outroExercicio.OrderIndex,
		Titulo:          outroExercicio.Titulo,
		CodigoBase:      outroExercicio.CodigoBase,
		CodigoTeste:     outroExercicio.CodigoTeste,
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	//  altera para a nova posicao do exercicio
	exercicio.OrderIndex = posicaoAtrocarOutroExercicio
	q.UpdateExerciciosBase(context.Background(), generated.UpdateExerciciosBaseParams{
		ID:              exercicio.ID,
		CreatedByUserID: exercicio.CreatedByUserID,
		ListaID:         exercicio.ListaID,
		OrderIndex:      exercicio.OrderIndex,
		Titulo:          exercicio.Titulo,
		CodigoBase:      exercicio.CodigoBase,
		CodigoTeste:     exercicio.CodigoTeste,
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	w.WriteHeader(http.StatusOK)
}
