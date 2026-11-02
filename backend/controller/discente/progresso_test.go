package discente

import (
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"testing"
)

func TestStudentExercisesPreservesFirstCompletion(t *testing.T) {
	first := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	second := pgtype.UUID{Bytes: [16]byte{2}, Valid: true}
	base := []generated.ExerciciosBase{{ID: second, OrderIndex: 2, CodigoBase: "template-2"}, {ID: first, OrderIndex: 1, CodigoBase: "template-1"}}
	saved := []generated.ExerciciosDiscente{{ExerciciosBaseID: first, CodigoRodou: true, Habilitado: true, CodigoBase: "solution"}, {ExerciciosBaseID: second, Habilitado: true}}
	dashboard := studentExercises(base, saved, false)
	runner := studentExercises(base, saved, true)
	if dashboard[0].Data.Id != first.String() || !dashboard[0].Data.CodigoRodou {
		t.Fatal("first exercise lost its completion or order")
	}
	if !dashboard[1].Data.Habilitado || dashboard[1].Data.CodigoRodou {
		t.Fatal("available exercise must remain incomplete")
	}
	if dashboard[0].Data.CodigoBase != "" {
		t.Fatal("dashboard should not include code")
	}
	if runner[0].Data.CodigoBase != "solution" || runner[1].Data.CodigoBase != "template-2" {
		t.Fatal("runner must use saved solution or starter code")
	}
	if base[0].ID != second {
		t.Fatal("mapping mutated source order")
	}
}
func TestStudentExercisesInitialAvailability(t *testing.T) {
	base := []generated.ExerciciosBase{{OrderIndex: 1, CodigoBase: "first"}, {OrderIndex: 2, CodigoBase: "second"}}
	items := studentExercises(base, nil, true)
	if !items[0].Data.Habilitado || items[0].Data.CodigoRodou || items[1].Data.Habilitado {
		t.Fatal("only the first exercise starts available, without completion")
	}
	if items[0].Data.CodigoBase != "first" || items[1].Data.CodigoBase != "" {
		t.Fatal("only available starter code should be included")
	}
	if len(studentExercises(nil, nil, false)) != 0 {
		t.Fatal("empty list must remain empty")
	}
}
