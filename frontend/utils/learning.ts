import type { Trilha, PrimeVueDataExercicio, RetornoExecucao } from '../types/index';

export function exerciseStatus(
  exercise: Pick<PrimeVueDataExercicio, 'codigo_rodou' | 'habilitado'>
) {
  if (exercise.codigo_rodou)
    return { label: 'Concluído', icon: 'pi pi-check-circle', color: 'text-emerald-700' };
  if (exercise.habilitado)
    return { label: 'Disponível', icon: 'pi pi-play-circle', color: 'text-sky-800' };
  return { label: 'Bloqueado', icon: 'pi pi-lock', color: 'text-slate-500' };
}

export function nextExercise(trilhas: Trilha[]) {
  for (const trilha of trilhas) {
    for (const lista of trilha.listas || []) {
      const index = (lista.exercicios || []).findIndex(
        ({ data }) => data.habilitado && !data.codigo_rodou
      );
      if (index >= 0)
        return {
          name: lista.exercicios[index].data.name,
          context: `${trilha.name} · ${lista.name}`,
          to: `/dashboard_discente/exercicio/lista/${lista.id}/${index}/${trilha.tipo_da_linguagem}`,
        };
    }
  }
  return null;
}

export function executionGuidance(result: RetornoExecucao) {
  if (result.out_put_from_user.success && result.out_put_unit_teste.success) {
    return {
      kind: 'success',
      title: 'Você conseguiu!',
      detail: 'Sua solução passou pelas verificações deste exercício.',
    };
  }
  const error = result.out_put_from_user.error || '';
  // Only describe errors reported by the student's program, not errors in the teacher's tests.
  const line = error.match(/\bline\s+(\d+)\b/i)?.[1] || error.match(/\.go:(\d+):\d+/)?.[1];
  const location = line ? ` Confira a linha ${line}.` : '';
  let detail =
    'Seu programa ainda não executou corretamente. Confira a mensagem técnica abaixo e tente novamente.';
  if (/IndentationError|TabError/.test(error))
    detail =
      'Confira os espaços no início das linhas. Linhas do mesmo bloco precisam ter o mesmo recuo.' +
      location;
  else if (/SyntaxError|syntax error/i.test(error))
    detail = 'Confira a escrita do código: parênteses, aspas e sinais de pontuação.' + location;
  else if (/NameError|undefined:/.test(error))
    detail =
      'Há um nome que o programa não reconheceu. Confira se a variável ou função foi definida e escrita da mesma forma.' +
      location;
  else if (/TypeError|cannot use|mismatched types/.test(error))
    detail =
      'Confira os tipos dos valores usados nessa operação. Um texto e um número, por exemplo, podem precisar de conversão.' +
      location;
  else if (/ZeroDivisionError|division by zero/i.test(error))
    detail =
      'Uma divisão tentou usar zero como divisor. Confira os valores antes de dividir.' + location;
  else if (/IndexError|index out of range/i.test(error))
    detail =
      'O código tentou acessar uma posição que não existe na lista. Confira o tamanho da lista e os índices usados.' +
      location;
  else if (/timeout|timed out|deadline exceeded/i.test(error))
    detail = 'A execução demorou além do limite. Confira se os laços têm uma condição de parada.';
  else if (result.out_put_from_user.success)
    detail =
      'Seu programa executou, mas ainda não passou por todas as verificações. Compare o resultado com o enunciado e confira os detalhes abaixo.';
  return { kind: 'warning', title: 'Vamos ajustar sua solução', detail };
}

export function isExecutionResult(value: unknown): value is RetornoExecucao {
  if (!value || typeof value !== 'object') return false;
  return ['out_put_from_user', 'out_put_unit_teste'].every((key) => {
    const item = (value as Record<string, unknown>)[key];
    return (
      item &&
      typeof item === 'object' &&
      typeof (item as { success?: unknown }).success === 'boolean' &&
      ['error', 'output'].every((field) => {
        const text = (item as Record<string, unknown>)[field];
        return text == null || typeof text === 'string';
      })
    );
  });
}
