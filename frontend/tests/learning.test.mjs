import test from 'node:test';
import assert from 'node:assert/strict';
import jiti from 'jiti';
const { exerciseStatus, nextExercise, executionGuidance, isExecutionResult } = jiti(
  import.meta.url
)('../utils/learning.ts');
const exercise = (id, enabled, done) => ({
  key: id,
  data: { id, name: id, habilitado: enabled, codigo_rodou: done },
});
const tracks = (items) => [
  {
    name: 'Python',
    tipo_da_linguagem: 2,
    listas: [{ id: 'list', name: 'Fundamentos', exercicios: items }],
  },
];
const execution = (error, runs = false, passes = false) => ({
  out_put_from_user: { success: runs, error, output: '' },
  out_put_unit_teste: { success: passes, error: '', output: '' },
});
test('availability does not imply completion', () => {
  assert.equal(exerciseStatus(exercise('1', true, false).data).label, 'Disponível');
  assert.equal(exerciseStatus(exercise('1', true, true).data).label, 'Concluído');
  assert.equal(exerciseStatus(exercise('1', false, false).data).label, 'Bloqueado');
});
test('next activity skips completed and locked exercises and uses the list position', () => {
  const next = nextExercise(
    tracks([
      exercise('done', true, true),
      exercise('locked', false, false),
      exercise('next', true, false),
    ])
  );
  assert.equal(next.name, 'next');
  assert.equal(next.to, '/dashboard_discente/exercicio/lista/list/2/2');
  assert.equal(nextExercise(tracks([exercise('done', true, true)])), null);
  assert.equal(nextExercise([]), null);
});
test('explanations identify Python and Go locations without treating test code as student errors', () => {
  assert.match(
    executionGuidance(execution('File "main.py", line 3\nIndentationError: unexpected indent'))
      .detail,
    /recuo.*linha 3/s
  );
  assert.match(
    executionGuidance(execution('main.go:12:4: undefined: total')).detail,
    /nome.*linha 12/s
  );
  const result = execution('', true, false);
  result.out_put_unit_teste.error = 'SyntaxError: line 99';
  assert.match(executionGuidance(result).detail, /não passou por todas/);
  assert.doesNotMatch(executionGuidance(result).detail, /linha 99/);
});
test('success requires both execution and checks; unknown errors remain available as technical details', () => {
  assert.equal(executionGuidance(execution('', true, true)).kind, 'success');
  assert.equal(executionGuidance(execution('internal detail', false, true)).kind, 'warning');
  assert.match(executionGuidance(execution('unknown')).detail, /mensagem técnica/);
  assert.equal(isExecutionResult({ error: 'unavailable' }), false);
  assert.equal(isExecutionResult(null), false);
  assert.equal(
    isExecutionResult({
      ...execution('', true, true),
      out_put_from_user: { success: false, error: 123 },
    }),
    false
  );
  assert.equal(isExecutionResult(execution('', true, true)), true);
});
