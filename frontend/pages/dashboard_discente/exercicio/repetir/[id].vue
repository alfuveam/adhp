<template>
  <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />
    <ConfirmDialog />

    <!-- Top Navigation / Breadcrumb -->
    <div class="flex items-center justify-between pb-2 border-b border-white/[0.08]">
      <NuxtLink 
        to="/dashboard_discente" 
        class="inline-flex items-center gap-2 text-sm font-medium text-slate-400 hover:text-white transition"
      >
        <i class="pi pi-arrow-left text-xs"></i>
        <span>Voltar ao Painel</span>
      </NuxtLink>

      <span class="mimo-badge mimo-badge-indigo text-[11px]">
        Modo de Repetição Espaçada
      </span>
    </div>

    <!-- FORM WORKSPACE -->
    <Form 
      v-slot="$form" 
      :initialValues="initialValues" 
      :resolver="resolver" 
      @submit="onFormSubmit" 
      class="w-full flex flex-col gap-6"
    >
      <!-- Exercise Header Card (Mimo Style) -->
      <div class="mimo-card p-6 border-[#363e66] bg-[#1c2038] shadow-lg flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="space-y-1.5">
          <div class="flex items-center gap-2.5">
            <span 
              class="px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider"
              :class="tipoDaLinguagem === 1 ? 'mimo-badge-blue' : 'mimo-badge-amber'"
            >
              {{ tipoDaLinguagem === 1 ? 'Golang' : 'Python' }}
            </span>
            <span class="text-xs text-purple-300 font-medium">• Revisão Agendada</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-extrabold text-white tracking-tight">
            {{ exercicio.name || 'Exercício' }}
          </h1>
        </div>

        <button 
          type="button" 
          @click="showFeedback(exercicio)"
          class="mimo-btn-secondary px-4 py-2 text-xs flex items-center gap-2"
        >
          <i class="pi pi-info-circle text-xs text-purple-300"></i>
          <span>Ver Dica do Exercício</span>
        </button>
      </div>

      <!-- Monaco Editor Workspace Frame (Mimo Style) -->
      <div class="rounded-2xl border border-[#363e66] overflow-hidden bg-[#15172a] shadow-xl">
        <!-- Window Bar -->
        <div class="flex items-center justify-between px-4 py-2.5 bg-[#1e233d] border-b border-[#2d3458]">
          <div class="flex items-center gap-1.5">
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
          </div>

          <div class="flex items-center gap-2 px-3 py-1 rounded-lg bg-[#141628] text-xs font-mono text-purple-300 border border-[#2d3458]">
            <i class="pi pi-file text-purple-400 text-xs"></i>
            <span>{{ tipoDaLinguagem === 1 ? 'repeticao.go' : 'repeticao.py' }}</span>
          </div>

          <span class="text-[11px] text-slate-400 font-mono hidden sm:inline-block">Monaco Editor</span>
        </div>

        <!-- Monaco Editor Mount -->
        <div class="relative w-full">
          <MonacoEditor 
            name="codigo_base" 
            v-model="exercicio.codigo_base" 
            :options="optionsEditor" 
            :lang="tipoDaLinguagem === 1 ? 'go' : 'python'" 
            class="h-[460px] sm:h-[520px] w-full z-10" 
          />
        </div>
      </div>

      <!-- Output Consoles (Stdout & Tests) -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        
        <!-- Console 1: Execução -->
        <div class="rounded-2xl border border-[#363e66] overflow-hidden bg-[#16182c] shadow-lg flex flex-col">
          <div class="flex items-center justify-between px-4 py-2.5 bg-[#1e233d] border-b border-[#2d3458]">
            <div class="flex items-center gap-2">
              <i class="pi pi-terminal text-slate-400 text-xs"></i>
              <span class="text-xs font-bold text-slate-200">Saída da Execução (Stdout/Stderr)</span>
            </div>
          </div>
          <div class="p-3.5 flex-1 bg-[#121424]">
            <textarea 
              v-model="retorno_codigo" 
              readonly 
              placeholder="A saída padrão do seu código aparecerá aqui..." 
              class="w-full h-32 bg-transparent text-slate-300 font-mono text-xs resize-none focus:outline-none placeholder-slate-600 leading-relaxed"
            ></textarea>
          </div>
        </div>

        <!-- Console 2: Testes Unitários -->
        <div class="rounded-2xl border border-[#363e66] overflow-hidden bg-[#16182c] shadow-lg flex flex-col">
          <div class="flex items-center justify-between px-4 py-2.5 bg-[#1e233d] border-b border-[#2d3458]">
            <div class="flex items-center gap-2">
              <i class="pi pi-verified text-emerald-400 text-xs"></i>
              <span class="text-xs font-bold text-slate-200">Resultado dos Testes Unitários</span>
            </div>
          </div>
          <div class="p-3.5 flex-1 bg-[#121424]">
            <textarea 
              v-model="retorno_teste" 
              readonly 
              placeholder="O relatório da bateria de testes aparecerá aqui..." 
              class="w-full h-32 bg-transparent text-emerald-400 font-mono text-xs resize-none focus:outline-none placeholder-slate-600 leading-relaxed"
            ></textarea>
          </div>
        </div>

      </div>

      <!-- Action Button -->
      <div class="flex justify-end pt-2">
        <button 
          type="submit" 
          class="mimo-btn-primary px-8 py-3 text-sm font-bold flex items-center justify-center gap-2"
        >
          <i class="pi pi-check text-xs"></i>
          <span>Submeter Revisão Espaçada</span>
        </button>
      </div>
    </Form>
  </div>
</template>

<script setup lang="ts">
import * as Monaco from 'monaco-editor';
import { ref, reactive } from 'vue';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import type { PrimeVueDataExercicio, RetornoExecucao } from '~/types';
import { useToast } from 'primevue/usetoast';

const optionsEditor = ref<Monaco.editor.IStandaloneEditorConstructionOptions>({
  theme: 'vs-dark',
  fontSize: 14,
  lineHeight: 22,
  fontFamily: '"JetBrains Mono", Menlo, Monaco, "Courier New", monospace',
  minimap: { enabled: false },
  scrollBeyondLastLine: false,
  automaticLayout: true,
  padding: { top: 12, bottom: 12 }
});

const route = useRoute();
const toast = useToast();
const { $authService } = useNuxtApp();
let refQuillCodigoBase = ref(null);
const exercicioIDRep = ref('');
const tipoDaLinguagem = ref(0);

const initialValues = reactive({
    codigo_base: ''
});

const exercicio = ref<PrimeVueDataExercicio>({
  name: '',
  id: '',
  inputValue: '',
  type: 0,
  codigo_rodou: false,
  codigo_base: '',
  habilitado: false,
  order_index: 0,
});

const retorno_codigo = ref('');
const retorno_teste = ref('');
const codigoBaseToSend = ref('');

const showFeedback = async (e: PrimeVueDataExercicio) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/feedback_by_exercicio_id/' + e.id + '/' + config.public.METRICAS_FEEDBACK.REPETICAO, {
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
      toast.add({ severity: 'error', summary: 'Erro', detail: "Erro ao requisitar o feedback.", life: 3000 });
    }
  } catch (error) {    
    console.log(error);
  }
};

const resolver = zodResolver(
  z.object({
    codigo_base: z.string().min(1, { message: 'Código base exercício é necessário.' })
  })
);

const onFormSubmit = async ({ valid }: { valid: boolean }) => {
  if (valid) {
    toast.add({
      severity: 'info',
      summary: 'Executando testes...',
      life: 2000
    });
  }

  if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
    codigoBaseToSend.value = exercicio.value.codigo_base;
  } else {
    // @ts-ignore
    codigoBaseToSend.value = refQuillCodigoBase.value?.quill?.getText();
  }

  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/discente_submit_exerc_rep_espacada', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        repeticao_espacada_id: exercicioIDRep.value,
        exercicio_id: exercicio.value.id,
        codigo_base: btoa(encodeURIComponent(codigoBaseToSend.value))
      })
    });

    const exercicioResponse: RetornoExecucao = await response.json();
    if (response.ok) {
      retorno_codigo.value = exercicioResponse.out_put_from_user.output;
      retorno_teste.value = exercicioResponse.out_put_unit_teste.output;

      if (!exercicioResponse.out_put_unit_teste.success) {
        retorno_teste.value = exercicioResponse.out_put_unit_teste.error || exercicioResponse.out_put_unit_teste.output;
      }
      if (exercicioResponse.out_put_from_user.success && exercicioResponse.out_put_unit_teste.success) {
        toast.add({ severity: 'success', summary: 'Revisão Concluída!', detail: 'Exercício validado com sucesso.', life: 3500 });
      } else {
        toast.add({ severity: 'error', summary: 'Atenção aos testes', detail: 'Verifique as falhas nos testes.', life: 3500 });
      }
    } else {
      retorno_codigo.value = exercicioResponse.out_put_from_user.error || 'Erro na execução';
      retorno_teste.value = exercicioResponse.out_put_unit_teste.error || 'Erro nos testes';
      toast.add({ severity: 'error', summary: 'Erro ao submeter', detail: 'Falha no processamento.', life: 3000 });
    }
  } catch(error) {
    console.log(error);
  }
};

const loadRepeticaoEspacadaExercicio = async () => {
  // @ts-ignore
  exercicioIDRep.value = route.params.id;
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/discente_get_exerc_rep_espacada/' + exercicioIDRep.value, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    });

    if (response.ok) {
      const data = await response.json();
      exercicio.value = data.exercicio;
      tipoDaLinguagem.value = data.tipo_da_linguagem;
      if (exercicio.value.codigo_base) {
        try {
          exercicio.value.codigo_base = decodeURIComponent(atob(exercicio.value.codigo_base));
        } catch {
          // already decoded
        }
      }
      toast.add({ severity: 'success', summary: 'Exercício pronto para revisão', life: 2500 });
      await onSubmitMetricaExercicioInicio();
    } else {
      toast.add({ severity: 'error', summary: 'Erro', detail: 'Falha ao carregar o exercício de repetição.', life: 3000 });
    }
  } catch(error) {
    console.log(error);
  }
};

const onSubmitMetricaExercicioInicio = async () => {
  try {
    const config = useRuntimeConfig();
    await fetch(config.public.BACKEND_API_URL + '/v1/metricas_repeticao_espacada', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        exercicio_id: exercicio.value.id,
        tipo_metrica: config.public.METRICAS.INICIO
      })
    });
  } catch (error) {
    console.log(error);
  }
};

loadRepeticaoEspacadaExercicio();
</script>