package models

type ExercicioData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type ExercicioDashboard struct {
	Key  string        `json:"key"`
	Data ExercicioData `json:"data"`
}

type ListaDashboard struct {
	Name       string               `json:"name"`
	Id         string               `json:"id"`
	Exercicios []ExercicioDashboard `json:"exercicios"`
}

type TrilhaDashboard struct {
	Name            string           `json:"name"`
	Id              string           `json:"id"`
	TipoDaLinguagem int16            `json:"tipo_da_linguagem"`
	Listas          []ListaDashboard `json:"listas"`
}
