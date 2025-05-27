package models

type ExercicioDataDiscente struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	CodigoRodou bool   `json:"codigo_rodou"`
	Habilitado  bool   `json:"habilitado"`
	CodigoBase  string `json:"codigo_base"`
	OrderIndex  int16  `json:"order_index"`
}

type ExercicioDashboardDiscente struct {
	Key  string                `json:"key"`
	Data ExercicioDataDiscente `json:"data"`
}

type ListaDashboardDiscente struct {
	Name       string                       `json:"name"`
	Id         string                       `json:"id"`
	Exercicios []ExercicioDashboardDiscente `json:"exercicios"`
}

type TrilhaDashboardDiscente struct {
	Name            string                   `json:"name"`
	Id              string                   `json:"id"`
	TipoDaLinguagem int16                    `json:"tipo_da_linguagem"`
	Listas          []ListaDashboardDiscente `json:"listas"`
}
