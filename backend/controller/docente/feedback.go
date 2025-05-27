package docente

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/alfuveam/adhp/backend/generated"
	"github.com/jackc/pgx/v5/pgtype"
)

func RemoveFeedback(w http.ResponseWriter, r *http.Request, queries *generated.Queries) {
	var uuidFeedback pgtype.UUID
	uuidFeedback.Scan(r.PathValue("id"))

	_, err := queries.DeleteFeedback(context.Background(), uuidFeedback)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
}
