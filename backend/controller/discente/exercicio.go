package discente

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

type DiscenteSubmitExercicio struct {
	ExercicioId string `json:"exercicio_id"`
	CodigoBase  string `json:"codigo_base"`
}

func GetFeedbackByExercicioBaseId(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidExercicioBase pgtype.UUID
	uuidExercicioBase.Scan(r.PathValue("id"))

	tipoFeedbackk, err := strconv.ParseInt(r.PathValue("tipo_feedback"), 10, 16)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao converter o tipo de feedback",
		})
		return
	}
	tipoFeedback := int16(tipoFeedbackk)
	fmt.Println("tipoFeedback: ", tipoFeedback)

	feedback, err := q.GetRandomFeedbackByExerciciosBaseId(context.Background(), uuidExercicioBase)
	if err != nil {
		log.Println("GetFeedbackById: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao buscar feedback",
		})
		return
	}

	//	métricas
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

	exercicio, err := q.GetExerciciosBaseById(context.Background(), uuidExercicioBase)
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

	_, err = q.CreateDiscenteMetricasFeedback(context.Background(), generated.CreateDiscenteMetricasFeedbackParams{
		CreatedByUserID:  userUUID,
		TrilhaID:         trilha.ID,
		ListaID:          exercicio.ListaID,
		ExerciciosBaseID: exercicio.ID,
		HorarioAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
		TipoExercicio:    tipoFeedback,
	})

	if err != nil {
		log.Println("CreateDiscenteMetricasFeedback: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao criar métrica",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"feedback": feedback.Descricao,
	})
}

func OnDiscenteSubmitExercicio(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req DiscenteSubmitExercicio
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

	var uuidExercicioBase pgtype.UUID
	uuidExercicioBase.Scan(req.ExercicioId)

	exercicio, err := q.GetExerciciosBaseById(context.Background(), uuidExercicioBase)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Não achou o exercicio base",
		})
		return
	}

	unitTestHandler := docente.UnitTestHandler{
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
			"error": "Erro ao buscar a lista",
		})
		return
	}
	trilha, err := q.GetTrilhaById(context.Background(), lista.TrilhaID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Erro ao buscar a trilha",
		})
		return
	}

	var uuidUser pgtype.UUID
	uuidUser.Scan(user.Id)

	codeHandler := docente.CodeHandlerMetricas{
		CreatedByUserId:  uuidUser,
		TrilhaId:         trilha.ID,
		ListaId:          lista.ID,
		ExerciciosBaseId: exercicio.ID,
		LinguagemTipo:    trilha.TipoDaLinguagem,
	}

	unitTestCodeDiscente := docente.OnRunTestCode(unitTestHandler, codeHandler, q)

	if unitTestCodeDiscente.OutPutFromUser.Success && unitTestCodeDiscente.OutPutUnitTeste.Success {
		_, err := q.UpsertExerciciosDiscente(context.Background(), generated.UpsertExerciciosDiscenteParams{
			ExerciciosBaseID: exercicio.ID,
			CreatedByUserID:  uuidUser,
			ListaID:          exercicio.ListaID,
			CodigoBase:       req.CodigoBase,
			CodigoRodou:      true,
			CreateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true}, //	no update não atualiza
			UpdateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
			Habilitado:       true, // nesse caso sempre verdadeiro, porque não a como responder um exercício sem ter respondido o anterior
		})

		if err != nil {
			log.Println("UpsertExerciciosDiscente: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erro ao atualizar exercício",
			})
			return
		}

		err = HandlerMetricasExercicio(uuidUser, exercicio.ID, types.MetricasRodou, q)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}

		proximoExercicio, _ := q.GetExercicioByOrderIndexAndListaId(context.Background(), generated.GetExercicioByOrderIndexAndListaIdParams{
			OrderIndex: exercicio.OrderIndex + 1,
			ListaID:    exercicio.ListaID,
		})

		temExercicio, _ := q.CheckIfHasExerciciosDiscenteByExerciciosBaseId(context.Background(), generated.CheckIfHasExerciciosDiscenteByExerciciosBaseIdParams{
			ExerciciosBaseID: proximoExercicio.ID,
			CreatedByUserID:  uuidUser,
		})

		if !temExercicio {
			q.CreateExerciciosDiscente(context.Background(), generated.CreateExerciciosDiscenteParams{
				ExerciciosBaseID: proximoExercicio.ID,
				CreatedByUserID:  uuidUser,
				ListaID:          exercicio.ListaID,
				CodigoBase:       "",
				CodigoRodou:      false,
				Habilitado:       true,
				CreateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
				UpdateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
			})
		}

		userDB, err := q.GetUserById(context.Background(), uuidUser)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erro ao buscar o usuário",
			})
			return
		}

		//	insere data de repetição
		dataRepeticao := pgtype.Timestamp{Time: time.Now().Add(time.Minute * time.Duration(userDB.RepeticaoEspacadaMinutos)), Valid: true}

		//	Nesse contexto não devo somar a repetição. nome da variavel no banco de dados: repeticao
		_, err = q.UpsertExercRepNaoSomaRep(context.Background(), generated.UpsertExercRepNaoSomaRepParams{
			CreatedByUserID:  uuidUser,
			ExerciciosBaseID: exercicio.ID,
			ProximaRepeticao: dataRepeticao,
			CreateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
			UpdateAt:         pgtype.Timestamp{Time: time.Now(), Valid: true},
		})

		if err != nil {
			log.Println("UpsertExercRepNaoSomaRep: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Erro ao inserir data de repetição do exercício ",
			})
			return
		}
	} else {
		err := HandlerMetricasExercicio(uuidUser, exercicio.ID, types.MetricasTentou, q)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(unitTestCodeDiscente)
}
