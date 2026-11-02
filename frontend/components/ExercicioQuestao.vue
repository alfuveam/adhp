<template>
  <div class="w-full max-w-full">
    <ConfirmDialog />

    <Form 
      v-slot="$form" 
      :initialValues="initialValues" 
      :resolver="resolver" 
      @submit="onFormSubmit" 
      class="w-full flex flex-col gap-6"
    >
      <!-- EXERCISE HEADER CARD (Mimo Style) -->
      <div class="mimo-card p-6 border-[#363e66] bg-[#1c2038] shadow-lg flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="space-y-1.5">
          <div class="flex items-center gap-2.5">
            <span 
              class="px-2.5 py-0.5 rounded-full text-xs font-bold uppercase tracking-wider"
              :class="tipoDaLinguagem === 1 ? 'mimo-badge-blue' : 'mimo-badge-amber'"
            >
              {{ tipoDaLinguagem === 1 ? 'Golang' : 'Python' }}
            </span>
            <span class="text-xs text-slate-400 font-medium">• Desafio Prático</span>
          </div>
          <h2 class="text-xl sm:text-2xl font-extrabold text-white tracking-tight">
            {{ exercicio.name || 'Exercício' }}
          </h2>
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

      <!-- CODE EDITOR WORKSPACE (Mimo Style Frame) -->
      <div class="rounded-2xl border border-[#363e66] overflow-hidden bg-[#15172a] shadow-xl">
        <!-- Editor Window Bar -->
        <div class="flex items-center justify-between px-4 py-2.5 bg-[#1e233d] border-b border-[#2d3458]">
          <div class="flex items-center gap-1.5">
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
          </div>

          <div class="flex items-center gap-2 px-3 py-1 rounded-lg bg-[#141628] text-xs font-mono text-purple-300 border border-[#2d3458]">
            <i class="pi pi-file text-purple-400 text-xs"></i>
            <span>{{ tipoDaLinguagem === 1 ? 'main.go' : 'main.py' }}</span>
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

      <!-- EXECUTION & TEST OUTPUT CONSOLES (Grid) -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        
        <!-- Console 1: Retorno da Execução -->
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

        <!-- Console 2: Retorno dos Testes -->
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

      <!-- SUBMIT BUTTON -->
      <div class="flex justify-end pt-2">
        <button 
          type="submit" 
          class="mimo-btn-primary px-8 py-3 text-sm font-bold flex items-center justify-center gap-2"
        >
          <i class="pi pi-send text-xs"></i>
          <span>{{ callBackSubmitText || 'Submeter Código' }}</span>
        </button>
      </div>
    </Form>
  </div>
</template>

<script setup lang="ts">
import * as Monaco from 'monaco-editor';
import { ref, reactive, watch, onMounted } from 'vue';
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
const exercicioID = ref('');

const tipoDaLinguagem = ref(0);
const initialValues = reactive({
    codigo_base: ''
});

const props = defineProps({
  exercicio: {
    type: Object as () => PrimeVueDataExercicio,
    required: true
  },
  callBackSubmitText: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['show-toast', 'call-back-atualizar-lista']);

const retorno_codigo = ref('');
const retorno_teste = ref('');
const codigoBaseToSend = ref('');

const showFeedback = (e: PrimeVueDataExercicio) => {
  emit('show-toast', e);
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
      summary: 'Executando código...',
      life: 2000
    });
  }

  if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
    codigoBaseToSend.value = props.exercicio.codigo_base;
  } else {
    // @ts-ignore
    codigoBaseToSend.value = refQuillCodigoBase.value?.quill?.getText();
  }

  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/discente_submit_exercicio', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        exercicio_id: props.exercicio.id,
        codigo_base: btoa(encodeURIComponent(codigoBaseToSend.value))
      })
    });

    const exercicioResponse: RetornoExecucao = await response.json();
    if (response.ok) {
      retorno_codigo.value = exercicioResponse.out_put_from_user.output;
      retorno_teste.value = exercicioResponse.out_put_unit_teste.output;

      if (exercicioResponse.out_put_from_user.success && exercicioResponse.out_put_unit_teste.success) {
        toast.add({ severity: 'success', summary: 'Parabéns!', detail: 'Exercício resolvido com sucesso!', life: 3500 });
        emit('call-back-atualizar-lista');
      } else {
        retorno_codigo.value = exercicioResponse.out_put_from_user.error || exercicioResponse.out_put_from_user.output;
        retorno_teste.value = exercicioResponse.out_put_unit_teste.error || exercicioResponse.out_put_unit_teste.output;
        toast.add({ severity: 'error', summary: 'Atenção aos testes', detail: 'Um ou mais testes não passaram.', life: 3500 });
      }
    } else {
      toast.add({ severity: 'error', summary: 'Erro ao submeter', detail: 'Falha ao processar o código no servidor.', life: 3000 });
    }
  } catch(error) {
    console.log(error);
  }
};

watch(async () => props.exercicio, async (newExercicio) => {
  try {
    const exercicioProps = await newExercicio;
    if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
      props.exercicio.codigo_base = decodeURIComponent(atob(props.exercicio.codigo_base));
    } else {
      const v = decodeURIComponent(atob(exercicioProps.codigo_base));
      // @ts-ignore
      refQuillCodigoBase.value?.quill?.setText(v);
    }
    if (exercicioID.value != '' && exercicioID.value != exercicioProps.id) {
      retorno_codigo.value = "";
      retorno_teste.value = "";
    }

    exercicioID.value = exercicioProps.id;
  } catch (error) {
    console.log(error);
  }
});

onMounted(() => {
  if (route.params.id && Array.isArray(route.params.id) && route.params.id.length > 2) {
    tipoDaLinguagem.value = Number(route.params.id[2]);
  }
  const intervalId = setInterval(checkIfQuillIsLoad, 100);

  function checkIfQuillIsLoad() {
    if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
      clearInterval(intervalId);
      if (props.exercicio?.codigo_base) {
        try {
          const v = decodeURIComponent(atob(props.exercicio.codigo_base));
          props.exercicio.codigo_base = v;
        } catch {
          // already decoded
        }
      }
    } else {
      // @ts-ignore
      if (refQuillCodigoBase && refQuillCodigoBase.value?.quill?.root) {
        clearInterval(intervalId);
        try {
          const v = decodeURIComponent(atob(props.exercicio.codigo_base));
          // @ts-ignore
          refQuillCodigoBase.value?.quill?.setText(v);
        } catch {
          // already decoded
        }
      }
    }
  }
});
</script>