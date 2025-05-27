package docente

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type ListaIndex struct {
	TrilhaID        pgtype.UUID `json:"trilha_id"`
	ListaID         pgtype.UUID `json:"lista_id"`
	PosicoesATrocar int16       `json:"posicoes_a_trocar"`
}

func AddLista(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req generated.Listum
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

	lastIndex, err := q.GetListaCountByTrilhaId(context.Background(), req.TrilhaID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	lastIndex = lastIndex + 1
	_, err = q.CreateLista(context.Background(), generated.CreateListaParams{
		Name:            "Lista " + fmt.Sprint(lastIndex),
		TrilhaID:        req.TrilhaID,
		OrderIndex:      lastIndex,
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

func RemoveLista(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var uuidLista pgtype.UUID
	uuidLista.Scan(r.PathValue("id"))

	_, err := q.DeleteLista(context.Background(), uuidLista)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateLista(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req generated.Listum
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Requisição com payload inválido",
		})
		return
	}

	lista, _ := q.GetListaById(context.Background(), req.ID)
	lista.Name = req.Name

	_, err = q.UpdateLista(context.Background(), generated.UpdateListaParams{
		ID:              lista.ID,
		CreatedByUserID: lista.CreatedByUserID,
		TrilhaID:        lista.TrilhaID,
		OrderIndex:      lista.OrderIndex,
		Name:            lista.Name,
		Description:     lista.Description,
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateListaIndex(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req ListaIndex
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Requisição com payload inválido",
		})
		return
	}

	lista, _ := q.GetListaById(context.Background(), req.ListaID)

	//	troca pela posicao do outro lista
	posicaoAtrocarOutroLista := lista.OrderIndex + req.PosicoesATrocar

	outroLista, err := q.GetListaByOrderIndexAndTrilhaId(context.Background(), generated.GetListaByOrderIndexAndTrilhaIdParams{
		OrderIndex: posicaoAtrocarOutroLista,
		TrilhaID:   req.TrilhaID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	outroLista.OrderIndex = lista.OrderIndex
	q.UpdateLista(context.Background(), generated.UpdateListaParams{
		ID:              outroLista.ID,
		CreatedByUserID: outroLista.CreatedByUserID,
		TrilhaID:        outroLista.TrilhaID,
		OrderIndex:      outroLista.OrderIndex,
		Name:            outroLista.Name,
		Description:     outroLista.Description,
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	//  altera para a nova posicao do lista
	lista.OrderIndex = posicaoAtrocarOutroLista
	q.UpdateLista(context.Background(), generated.UpdateListaParams{
		ID:              lista.ID,
		CreatedByUserID: lista.CreatedByUserID,
		TrilhaID:        lista.TrilhaID,
		OrderIndex:      lista.OrderIndex,
		Name:            lista.Name,
		Description:     lista.Description,
		UpdateAt:        pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	w.WriteHeader(http.StatusOK)
}
