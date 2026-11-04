<template>
  <form class="w-full min-w-0 space-y-5" @submit.prevent="submitSolution" :aria-busy="submitting">
    <div class="study-card p-5 sm:p-6 flex flex-wrap items-start justify-between gap-4">
      <div class="min-w-0 space-y-2">
        <span class="study-badge study-badge-blue">{{ language === 1 ? 'Go' : 'Python' }}</span>
        <h2 class="text-xl font-semibold text-slate-800 break-words">
          {{ exercicio.name || 'Exercício' }}
        </h2>
      </div>
      <button
        type="button"
        class="study-btn-secondary px-4 py-2 text-sm"
        :disabled="submitting || !exercicio.id"
        @click="$emit('show-toast', exercicio)"
      >
        <i class="pi pi-info-circle" aria-hidden="true"></i> Preciso de uma dica
      </button>
    </div>

    <div class="editor-panel">
      <div class="editor-toolbar">
        <span>{{ language === 1 ? 'main.go' : 'main.py' }}</span
        ><span>Escreva sua solução</span>
      </div>
      <MonacoEditor
        v-model="code"
        :options="editorOptions"
        :lang="language === 1 ? 'go' : 'python'"
        class="h-[400px] sm:h-[480px] w-full"
      />
    </div>
    <div class="flex flex-wrap items-center justify-between gap-3">
      <p class="text-sm text-slate-600">
        Você pode testar e ajustar sua solução quantas vezes precisar.
      </p>
      <button
        type="submit"
        class="study-btn-primary px-5 py-3 text-sm"
        :disabled="submitting || !exercicio.id"
      >
        <i :class="submitting ? 'pi pi-spin pi-spinner' : 'pi pi-play'" aria-hidden="true"></i>
        {{ submitting ? 'Verificando sua solução…' : callBackSubmitText || 'Testar minha solução' }}
      </button>
    </div>

    <section
      v-if="message"
      role="status"
      aria-live="polite"
      class="study-card p-5 border-l-4"
      :class="message.kind === 'success' ? 'border-l-emerald-700' : 'border-l-amber-700'"
    >
      <h3 class="font-semibold text-slate-800">{{ message.title }}</h3>
      <p class="mt-2 text-sm text-slate-600">{{ message.detail }}</p>
    </section>
    <section v-if="result" class="study-card p-5 space-y-4" aria-label="Resultado da tentativa">
      <div>
        <h3 class="font-semibold text-slate-800">Saída do programa</h3>
        <pre
          v-if="result.out_put_from_user.output"
          class="study-output mt-3"
          tabindex="0"
          aria-label="Saída do programa"
          >{{ result.out_put_from_user.output }}</pre
        >
        <p v-else class="mt-2 text-sm text-slate-600">
          {{
            result.out_put_from_user.success
              ? 'Seu programa não exibiu uma mensagem nesta execução.'
              : 'O programa não produziu uma saída antes do erro.'
          }}
        </p>
      </div>
      <details
        v-if="
          result.out_put_from_user.error ||
          result.out_put_unit_teste.error ||
          result.out_put_unit_teste.output
        "
        class="border-t border-slate-200 pt-4"
      >
        <summary class="cursor-pointer text-sm font-semibold text-sky-800">
          Ver mensagens técnicas e verificações
        </summary>
        <div class="space-y-3 mt-3">
          <div v-if="result.out_put_from_user.error">
            <h4 class="text-sm font-semibold">Mensagem da execução</h4>
            <pre class="study-output mt-2" tabindex="0">{{ result.out_put_from_user.error }}</pre>
          </div>
          <div v-if="result.out_put_unit_teste.error || result.out_put_unit_teste.output">
            <h4 class="text-sm font-semibold">Verificações do exercício</h4>
            <pre class="study-output mt-2" tabindex="0">{{
              result.out_put_unit_teste.error || result.out_put_unit_teste.output
            }}</pre>
          </div>
        </div>
      </details>
    </section>
  </form>
</template>

<script setup lang="ts">
import type { editor } from 'monaco-editor';
import type { PrimeVueDataExercicio, RetornoExecucao } from '~/types';
import { executionGuidance, isExecutionResult } from '~/utils/learning';

const props = defineProps<{
  exercicio: PrimeVueDataExercicio;
  callBackSubmitText?: string;
  tipoLinguagem?: number;
  repeticaoId?: string;
}>();
const emit = defineEmits(['show-toast', 'call-back-atualizar-lista', 'busy']);
const route = useRoute();
const { $authService } = useNuxtApp();
const language = computed(
  () => props.tipoLinguagem ?? Number(Array.isArray(route.params.id) ? route.params.id[2] : 0)
);
const code = ref('');
const submitting = ref(false);
const result = ref<RetornoExecucao | null>(null);
const message = ref<ReturnType<typeof executionGuidance> | null>(null);
const editorOptions = computed<editor.IStandaloneEditorConstructionOptions>(() => ({
  theme: 'vs-dark',
  fontSize: 15,
  lineHeight: 24,
  fontFamily: '"JetBrains Mono", monospace',
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  automaticLayout: true,
  padding: { top: 16, bottom: 16 },
  readOnly: submitting.value,
  ariaLabel: 'Código da sua solução',
  tabFocusMode: true,
}));

watch(
  () => [props.exercicio.id, props.exercicio.codigo_base],
  ([id, source], previous) => {
    if (previous?.[0] === id && (submitting.value || result.value)) return;
    try {
      code.value = decodeURIComponent(atob(source || ''));
    } catch {
      code.value = source || '';
    }
    result.value = null;
    message.value = null;
  },
  { immediate: true }
);

watch(code, () => {
  if (!submitting.value && result.value) {
    result.value = null;
    message.value = null;
  }
});

async function submitSolution() {
  if (submitting.value || !props.exercicio.id) return;
  if (!code.value.trim()) {
    message.value = {
      kind: 'warning',
      title: 'Escreva sua solução primeiro',
      detail: 'Use o editor acima para escrever seu código antes de testar.',
    };
    return;
  }
  submitting.value = true;
  emit('busy', true);
  message.value = null;
  result.value = null;
  const exerciseId = props.exercicio.id;
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 60000);
  try {
    const config = useRuntimeConfig();
    const endpoint = props.repeticaoId
      ? '/v1/discente_submit_exerc_rep_espacada'
      : '/v1/discente_submit_exercicio';
    const response = await fetch(config.public.BACKEND_API_URL + endpoint, {
      method: 'POST',
      signal: controller.signal,
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + $authService.getToken(),
      },
      body: JSON.stringify({
        exercicio_id: exerciseId,
        codigo_base: btoa(encodeURIComponent(code.value)),
        ...(props.repeticaoId ? { repeticao_espacada_id: props.repeticaoId } : {}),
      }),
    });
    if (!response.ok) throw new Error('request');
    const data: unknown = await response.json();
    if (!isExecutionResult(data)) throw new Error('response');
    if (props.exercicio.id !== exerciseId) return;
    result.value = data;
    message.value = executionGuidance(data);
    if (data.out_put_from_user.success && data.out_put_unit_teste.success)
      emit('call-back-atualizar-lista');
  } catch {
    if (props.exercicio.id === exerciseId)
      message.value = {
        kind: 'warning',
        title: 'Não foi possível verificar sua solução',
        detail:
          'A execução não terminou ou houve uma falha de conexão. Seu código continua no editor. Tente novamente em instantes.',
      };
  } finally {
    clearTimeout(timeout);
    submitting.value = false;
    emit('busy', false);
  }
}
</script>
