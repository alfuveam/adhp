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
      <div
        class="study-card p-6 border-[#dce1d9] bg-[#ffffff] shadow-none flex flex-col sm:flex-row sm:items-center justify-between gap-4"
      >
        <div class="space-y-1.5">
          <div class="flex items-center gap-2.5">
            <span
              class="px-2.5 py-0.5 rounded-full text-sm font-semibold"
              :class="tipoDaLinguagem === 1 ? 'study-badge-blue' : 'study-badge-amber'"
            >
              {{ tipoDaLinguagem === 1 ? "Golang" : "Python" }}
            </span>
            <span class="text-sm text-slate-600 font-medium">• Desafio Prático</span>
          </div>
          <h2 class="text-xl sm:text-2xl font-semibold text-slate-800 tracking-tight">
            {{ exercicio.name || "Exercício" }}
          </h2>
        </div>

        <button
          type="button"
          @click="showFeedback(exercicio)"
          class="study-btn-secondary px-4 py-2 text-sm flex items-center gap-2"
        >
          <i class="pi pi-info-circle text-sm text-sky-800"></i>
          <span>Preciso de uma dica</span>
        </button>
      </div>

      <div class="editor-panel">
        <div class="editor-toolbar">
          <span>{{ tipoDaLinguagem === 1 ? "main.go" : "main.py" }}</span
          ><span>Escreva sua solução</span>
        </div>
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

      <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
        <div
          class="rounded-lg border border-[#dce1d9] overflow-hidden bg-[#ffffff] shadow-none flex flex-col"
        >
          <div
            class="flex items-center justify-between px-4 py-2.5 bg-[#ffffff] border-b border-[#dce1d9]"
          >
            <div class="flex items-center gap-2">
              <i class="pi pi-terminal text-slate-600 text-sm"></i>
              <span class="text-sm font-semibold text-slate-700">Saída do programa</span>
            </div>
          </div>
          <div class="p-3.5 flex-1 bg-[#ffffff]">
            <textarea
              v-model="retorno_codigo"
              readonly
              placeholder="A saída padrão do seu código aparecerá aqui..."
              class="w-full h-32 bg-transparent text-slate-600 font-mono text-sm resize-none focus-visible:outline-sky-800 placeholder-slate-600 leading-relaxed"
            ></textarea>
          </div>
        </div>

        <div
          class="rounded-lg border border-[#dce1d9] overflow-hidden bg-[#ffffff] shadow-none flex flex-col"
        >
          <div
            class="flex items-center justify-between px-4 py-2.5 bg-[#ffffff] border-b border-[#dce1d9]"
          >
            <div class="flex items-center gap-2">
              <i class="pi pi-verified text-emerald-700 text-sm"></i>
              <span class="text-sm font-semibold text-slate-700">Resultado da verificação</span>
            </div>
          </div>
          <div class="p-3.5 flex-1 bg-[#ffffff]">
            <textarea
              v-model="retorno_teste"
              readonly
              placeholder="O relatório da bateria de testes aparecerá aqui..."
              class="w-full h-32 bg-transparent text-emerald-700 font-mono text-sm resize-none focus-visible:outline-sky-800 placeholder-slate-600 leading-relaxed"
            ></textarea>
          </div>
        </div>
      </div>

      <div class="flex justify-end pt-2">
        <button
          type="submit"
          class="study-btn-primary px-8 py-3 text-sm font-semibold flex items-center justify-center gap-2"
        >
          <i class="pi pi-send text-sm"></i>
          <span>{{ callBackSubmitText || "Testar minha solução" }}</span>
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