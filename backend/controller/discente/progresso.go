package discente

import (
	"sort"
	"strconv"

	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
)

// Keep the dashboard and the exercise runner consistent about availability and completion.
func studentExercises(base []generated.ExerciciosBase, saved []generated.ExerciciosDiscente, includeCode bool) []models.ExercicioDashboardDiscente {
	ordered := append([]generated.ExerciciosBase(nil), base...)
	sort.SliceStable(ordered, func(i, j int) bool { return ordered[i].OrderIndex < ordered[j].OrderIndex })
	result := make([]models.ExercicioDashboardDiscente, 0, len(ordered))
	for _, exercise := range ordered {
		data := models.ExercicioDataDiscente{Id: exercise.ID.String(), Name: exercise.Titulo, OrderIndex: exercise.OrderIndex, Habilitado: exercise.OrderIndex == 1}
		for _, progress := range saved {
			if progress.ExerciciosBaseID == exercise.ID {
				data.CodigoRodou = progress.CodigoRodou
				data.Habilitado = data.Habilitado || progress.Habilitado || progress.CodigoRodou
				if includeCode {
					data.CodigoBase = progress.CodigoBase
				}
				break
			}
		}
		if includeCode && data.Habilitado && data.CodigoBase == "" {
			data.CodigoBase = exercise.CodigoBase
		}
		result = append(result, models.ExercicioDashboardDiscente{Key: strconv.Itoa(int(exercise.OrderIndex)), Data: data})
	}
	return result
}
