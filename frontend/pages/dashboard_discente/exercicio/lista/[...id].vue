<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />

    <!-- Top Navigation / Breadcrumb -->
    <div class="flex items-center justify-between pb-2 border-b border-white/[0.08]">
      <NuxtLink 
        to="/dashboard_discente" 
        class="inline-flex items-center gap-2 text-sm font-medium text-slate-400 hover:text-white transition"
      >
        <i class="pi pi-arrow-left text-xs"></i>
        <span>Voltar ao Painel de Trilhas</span>
      </NuxtLink>

      <div class="flex items-center gap-2">
        <span class="mimo-badge mimo-badge-indigo text-[11px]">
          Lista: {{ lista?.name || 'Carregando...' }}
        </span>
      </div>
    </div>

    <!-- Stepper / Exercise Runner Layout -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      
      <!-- Left Sidebar: Exercise Sequence Steps (Cols 1-4) -->
      <div class="lg:col-span-4 mimo-card p-5 border-[#363e66] bg-[#1c2038] shadow-xl space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-[#2d3458]">
          <h3 class="text-sm font-extrabold text-white uppercase tracking-wider flex items-center gap-2">
            <i class="pi pi-list text-purple-400 text-xs"></i>
            <span>Desafios da Lista</span>
          </h3>
          <span class="text-xs text-purple-300 font-mono font-bold">
            {{ (currentStep + 1) }}/{{ lista?.exercicios?.length || 0 }}
          </span>
        </div>

        <!-- Exercise Step List -->
        <div class="space-y-2 max-h-[580px] overflow-y-auto pr-1">
          <div 
            v-for="(exercicio, index) in lista?.exercicios" 
            :key="index"
            @click="exercicio.data.habilitado && (currentStep = index, exercicioSelecionado(exercicio.data))"
            class="p-3.5 rounded-xl border transition-all duration-200 flex items-center justify-between gap-3"
            :class="[
              currentStep === index 
                ? 'bg-[#7c4ed4]/20 border-[#7c4ed4] text-white font-bold shadow-md shadow-purple-500/10' 
                : exercicio.data.habilitado 
                  ? 'bg-[#15172a] border-[#2d3458] hover:border-[#424c7e] text-slate-200 cursor-pointer' 
                  : 'bg-[#121422] border-[#232742] text-slate-500 cursor-not-allowed opacity-60'
            ]"
          >
            <div class="flex items-center gap-3 truncate">
              <span 
                class="w-6 h-6 rounded-lg flex items-center justify-center text-xs font-mono font-bold shrink-0"
                :class="currentStep === index ? 'bg-[#7c4ed4] text-white' : 'bg-[#222742] text-slate-400'"
              >
                {{ index + 1 }}
              </span>
              <span class="text-xs truncate">
                {{ exercicio.data.name }}
              </span>
            </div>

            <span class="shrink-0 text-xs">
              <i v-if="!exercicio.data.habilitado" class="pi pi-lock text-slate-500"></i>
              <i v-else-if="currentStep === index" class="pi pi-play text-purple-400"></i>
              <i v-else class="pi pi-check-circle text-emerald-400"></i>
            </span>
          </div>
        </div>
      </div>

      <!-- Right Main: Monaco Editor Workstation (Cols 5-12) -->
      <div class="lg:col-span-8">
        <ExercicioQuestao 
          :exercicio="exercicioSelecionada" 
          @call-back-atualizar-lista="atualizarLista" 
          callBackSubmitText="Submeter Código"  
          @show-toast="showToastFeedBack"
        />
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import type { Listas, PrimeVueDataExercicio } from '~/types';
import { useToast } from 'primevue/usetoast';

const toast = useToast();
const route = useRoute();
const { $authService } = useNuxtApp();

const currentStep = ref(0);
const idLista = ref('');

const showToastFeedBack = async (e: PrimeVueDataExercicio) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/feedback_by_exercicio_id/' + e.id + '/' + config.public.METRICAS_FEEDBACK.EXERCICIO, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    });

    if (response.ok) {
      const data = await response.json();
      toast.add({ severity: 'info', summary: 'Dica do Assistente', detail: data.feedback, life: 5000 });
    } else {
      toast.add({ severity: 'error', summary: 'Erro', detail: "Não foi possível obter o feedback no momento.", life: 3000 });
    }
  } catch (error) {    
    console.log(error);
  }
};

const lista = ref<Listas>({
  id: '',
  name: '',
  exercicios: [],
  isEditing: false
});

const exercicioSelecionada = ref<PrimeVueDataExercicio>({
  name: '',
  id: '',
  inputValue: '',
  type: 0,
  codigo_rodou: false,
  codigo_base: '',
  habilitado: false,
  order_index: 0,
});

const exercicioSelecionado = (exercicio: PrimeVueDataExercicio) => {
  exercicioSelecionada.value = exercicio;
  onSubmitMetricaExercicioInicio(exercicio);
};

const atualizarLista = async () => {
  await getExerciciosHabilitadosByLista();
};

const getExerciciosHabilitadosByLista = async () => {
  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/exercicios_habilitados_by_lista/' + idLista.value, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + $authService.getToken()
      }
    });

    if (response.ok) {
      const resData = await response.json();
      if (resData && resData.length > 0) {
        lista.value = resData[0];
        if (lista.value.exercicios && lista.value.exercicios[currentStep.value]) {
          exercicioSelecionada.value = lista.value.exercicios[currentStep.value].data;
        }
      }
    } else {
      toast.add({ severity: 'error', summary: 'Erro', detail: 'Falha ao carregar a lista de exercícios.', life: 3000 });
    }
  } catch (error) {
    console.log(error);
  }
};

const onSubmitMetricaExercicioInicio = async (exercicio: PrimeVueDataExercicio) => {
  try {
    await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/metricas_exercicio', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        exercicio_id: exercicio.id,
        tipo_metrica: useRuntimeConfig().public.METRICAS.INICIO
      })
    });
  } catch (error) {
    console.log(error);
  }
};

onMounted(async () => {
  if (route.params.id && Array.isArray(route.params.id)) {
    idLista.value = route.params.id[0];
    currentStep.value = Number(route.params.id[1]) || 0;
  }
  await getExerciciosHabilitadosByLista();
});
</script>