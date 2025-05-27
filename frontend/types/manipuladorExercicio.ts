
export interface Feedback {
    id: string,
    descricao: string,
}

export interface ExercicioManipulador {
	id: string,
	titulo: string,
	lista_id: string,
	order_index: number,
	codigo_base: string,
	codigo_teste: string,
	feedbacks: Feedback[]
}

export interface RetornoExecucao {
  out_put_from_user: {
    success: boolean,
    error: string,
    output: string,
  },
  out_put_unit_teste: {
    success: boolean,
    error: string,
    output: string,
  },
}