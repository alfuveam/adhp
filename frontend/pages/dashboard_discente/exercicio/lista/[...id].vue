<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />
    <NuxtLink to="/dashboard_discente" class="inline-flex items-center gap-2 text-sm text-sky-800"
      ><i class="pi pi-arrow-left" aria-hidden="true"></i> Voltar aos meus estudos</NuxtLink
    >
    <p v-if="loading" role="status" class="study-card p-6">Carregando os exercícios da lista…</p>
    <div v-else-if="loadError" role="alert" class="study-card p-6 space-y-3">
      <p>{{ loadError }}</p>
      <button type="button" class="study-btn-secondary px-4 py-2" @click="loadList">
        Tentar novamente
      </button>
    </div>
    <p v-else-if="!lista.exercicios.length" class="study-card p-6">
      Esta lista ainda não tem exercícios. Volte ao painel para escolher outra lista.
    </p>
    <div v-else class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <nav aria-label="Exercícios da lista" class="lg:col-span-4 study-card p-5 space-y-4 min-w-0">
        <div>
          <h1 class="font-semibold text-lg">
            {{ lista.name === '-' ? 'Exercícios da lista' : lista.name }}
          </h1>
          <p class="text-sm text-slate-600">
            {{ completedCount }} de {{ lista.exercicios.length }} concluídos
          </p>
        </div>
        <ol class="space-y-2 max-h-[580px] overflow-y-auto p-1">
          <li v-for="(item, index) in lista.exercicios" :key="item.data.id">
            <button
              type="button"
              class="w-full p-3 rounded-lg border flex items-start gap-3 text-left"
              :class="currentStep === index ? 'bg-sky-900/5 border-sky-800' : 'border-slate-200'"
              :disabled="busy || !item.data.habilitado"
              :aria-current="currentStep === index ? 'step' : undefined"
              @click="selectExercise(index)"
            >
              <span class="text-sm text-slate-500">{{ index + 1 }}.</span>
              <span class="min-w-0 flex-1"
                ><span class="block text-sm break-words">{{ item.data.name }}</span
                ><span class="block text-xs mt-1" :class="exerciseStatus(item.data).color"
                  ><i :class="exerciseStatus(item.data).icon" aria-hidden="true"></i>
                  {{ exerciseStatus(item.data).label }}</span
                ></span
              >
            </button>
          </li>
        </ol>
      </nav>
      <div class="lg:col-span-8 min-w-0">
        <ExercicioQuestao
          v-if="selected?.habilitado"
          :exercicio="selected"
          callBackSubmitText="Testar minha solução"
          @busy="busy = $event"
          @show-toast="showHint"
          @call-back-atualizar-lista="refreshList"
        />
        <p v-else class="study-card p-6" role="status">
          Este exercício ainda está bloqueado. Escolha um exercício disponível na lista para
          continuar.
        </p>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import type { Listas, PrimeVueDataExercicio } from '~/types';
import { exerciseStatus } from '~/utils/learning';
import { useToast } from 'primevue/usetoast';
definePageMeta({ middleware: 'auth-discente' });
const route = useRoute();
const { $authService } = useNuxtApp();
const toast = useToast();
const loading = ref(true);
const loadError = ref('');
const busy = ref(false);
const currentStep = ref(0);
const lista = ref<Listas>({ id: '', name: '', isEditing: false, exercicios: [] });
const selected = computed(() => lista.value.exercicios[currentStep.value]?.data);
const completedCount = computed(
  () => lista.value.exercicios.filter((item) => item.data.codigo_rodou).length
);
const listId = computed(() => (Array.isArray(route.params.id) ? route.params.id[0] : ''));
async function loadList() {
  if (import.meta.server) return;
  loading.value = true;
  loadError.value = '';
  try {
    await fetchList();
  } catch {
    loadError.value =
      'Não foi possível carregar esta lista. Seu progresso continua salvo. Tente novamente.';
  } finally {
    loading.value = false;
  }
}
async function fetchList() {
  const response = await fetch(
    useRuntimeConfig().public.BACKEND_API_URL +
      '/v1/exercicios_habilitados_by_lista/' +
      listId.value,
    {
      signal: AbortSignal.timeout(20000),
      headers: { Authorization: 'Bearer ' + $authService.getToken() },
    }
  );
  if (!response.ok) throw new Error('request');
  const data = await response.json();
  if (!Array.isArray(data)) throw new Error('response');
  lista.value = data[0] || { id: listId.value, name: '', isEditing: false, exercicios: [] };
  lista.value.exercicios ||= [];
  if (currentStep.value >= lista.value.exercicios.length || currentStep.value < 0)
    currentStep.value = 0;
}
async function refreshList() {
  try {
    await fetchList();
  } catch {
    toast.add({
      severity: 'warn',
      summary: 'Sua solução foi aprovada',
      detail:
        'Não foi possível atualizar a lista agora. Recarregue a lista para ver o próximo exercício.',
      life: 6000,
    });
  }
}
function selectExercise(index: number) {
  if (busy.value || !lista.value.exercicios[index]?.data.habilitado) return;
  currentStep.value = index;
  const exercise = selected.value;
  if (!exercise) return;
  fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/metricas_exercicio', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: 'Bearer ' + $authService.getToken(),
    },
    body: JSON.stringify({
      exercicio_id: exercise.id,
      tipo_metrica: useRuntimeConfig().public.METRICAS.INICIO,
    }),
  }).catch(() => {});
}
async function showHint(exercise: PrimeVueDataExercicio) {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(
      config.public.BACKEND_API_URL +
        '/v1/feedback_by_exercicio_id/' +
        exercise.id +
        '/' +
        config.public.METRICAS_FEEDBACK.EXERCICIO,
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
watch(
  () => route.params.id,
  () => {
    currentStep.value = Math.max(
      0,
      Math.floor(Number(Array.isArray(route.params.id) ? route.params.id[1] : 0) || 0)
    );
    loadList();
  },
  { immediate: true }
);
</script>
