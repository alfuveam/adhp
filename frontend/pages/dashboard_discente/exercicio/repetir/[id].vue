<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />
    <NuxtLink to="/dashboard_discente" class="inline-flex items-center gap-2 text-sm text-sky-800"
      ><i class="pi pi-arrow-left" aria-hidden="true"></i> Voltar aos meus estudos</NuxtLink
    >
    <h1 class="text-xl font-semibold">Hora de revisar</h1>
    <p v-if="loading" role="status" class="study-card p-6">Carregando o exercício para revisão…</p>
    <div v-else-if="loadError" role="alert" class="study-card p-6 space-y-3">
      <p>{{ loadError }}</p>
      <button type="button" class="study-btn-secondary px-4 py-2" @click="loadReview">
        Tentar novamente
      </button>
    </div>
    <ExercicioQuestao
      v-else-if="exercise"
      :exercicio="exercise"
      :tipo-linguagem="language"
      :repeticao-id="reviewId"
      callBackSubmitText="Testar minha revisão"
      @show-toast="showHint"
      @call-back-atualizar-lista="completed = true"
    />
    <NuxtLink v-if="completed" to="/dashboard_discente" class="study-btn-secondary px-5 py-3"
      >Revisão concluída · Voltar aos meus estudos</NuxtLink
    >
  </div>
</template>
<script setup lang="ts">
import type { PrimeVueDataExercicio } from '~/types';
import { useToast } from 'primevue/usetoast';
definePageMeta({ middleware: 'auth-discente' });
const route = useRoute();
const { $authService } = useNuxtApp();
const toast = useToast();
const reviewId = computed(() => String(route.params.id || ''));
const exercise = ref<PrimeVueDataExercicio | null>(null);
const language = ref(0);
const loading = ref(true);
const loadError = ref('');
const completed = ref(false);
async function loadReview() {
  loading.value = true;
  loadError.value = '';
  completed.value = false;
  try {
    const response = await fetch(
      useRuntimeConfig().public.BACKEND_API_URL +
        '/v1/discente_get_exerc_rep_espacada/' +
        reviewId.value,
      {
        signal: AbortSignal.timeout(20000),
        headers: { Authorization: 'Bearer ' + $authService.getToken() },
      }
    );
    if (!response.ok) throw new Error('request');
    const data = await response.json();
    if (!data.exercicio?.id) throw new Error('response');
    exercise.value = data.exercicio;
    language.value = data.tipo_da_linguagem;
    fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/metricas_repeticao_espacada', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ' + $authService.getToken(),
      },
      body: JSON.stringify({
        exercicio_id: data.exercicio.id,
        tipo_metrica: useRuntimeConfig().public.METRICAS.INICIO,
      }),
    }).catch(() => {});
  } catch {
    loadError.value = 'Não foi possível carregar esta revisão. Tente novamente em instantes.';
  } finally {
    loading.value = false;
  }
}
async function showHint(item: PrimeVueDataExercicio) {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(
      config.public.BACKEND_API_URL +
        '/v1/feedback_by_exercicio_id/' +
        item.id +
        '/' +
        config.public.METRICAS_FEEDBACK.REPETICAO,
      {
        signal: AbortSignal.timeout(20000),
        headers: { Authorization: 'Bearer ' + $authService.getToken() },
      }
    );
    if (!response.ok) throw new Error('request');
    const data = await response.json();
    toast.add({
      severity: 'info',
      summary: 'Uma dica para continuar',
      detail: data.feedback || 'Este exercício ainda não tem uma dica disponível.',
      life: 8000,
    });
  } catch {
    toast.add({
      severity: 'error',
      summary: 'Não foi possível carregar a dica',
      detail: 'Tente novamente em instantes.',
      life: 4000,
    });
  }
}
onMounted(loadReview);
</script>
