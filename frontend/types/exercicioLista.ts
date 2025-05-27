export interface PrimeVueDataExercicio {
  name: string;
  id: string;
  inputValue: string;
  type: number;
  codigo_rodou: boolean;
  codigo_base: string;
  habilitado: boolean;
  order_index: number;
}

export interface PrimeVueDataExercicioKey {
  key: string;
  data: PrimeVueDataExercicio;
}

export interface Listas {
    name: string;
    id: string;
    isEditing: boolean;
    exercicios: PrimeVueDataExercicioKey[];
}

export interface Trilha {
  name: string;
  id: string;
  tipo_da_linguagem: number;
  isEditing: boolean;
  listas: Listas[];
}

export interface ExerciciosLista {
    id: string,
    questao: string,
    codigo_usuario: string,
    retorno_codigo: string,
}