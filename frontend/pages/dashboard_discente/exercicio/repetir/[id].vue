<template>
  <div class="card flex justify-center">
    <Toast />
    <ConfirmDialog></ConfirmDialog>

    <Form v-slot="$form" :initialValues :resolver @submit="onFormSubmit" class="flex flex-col gap-4 w-2/3 sm:w-2/3">
      <div class="flex flex-col gap-1">
        <label for="multiple-ac-1" class="font-bold mb-2 block">Título exercício</label>
        <p>
          {{exercicio.name}}
        </p>
      </div>
      <div class="flex flex-col gap-1">
        <div class="card flex justify-between">
          <label for="multiple-ac-1" class="font-bold block">Código para Submeter</label>
          <Button class="mb-2" label="Solicitar feedback" @click="showFeedback(exercicio)" />
        </div>
        <MonacoEditor name="codigo_base" v-model="exercicio.codigo_base" :options="optionsEditor" :lang="tipoDaLinguagem == 1 ? 'go' : 'python'" class="h-[600px] z-10" />
        <!-- <Editor v-model="exercicio.codigo_base" editorStyle="height: 220px" spellcheck="false" ref="refQuillCodigoBase">
          <template v-slot:toolbar>
            <span class="ql-formats">
              <button v-tooltip.bottom="'Bold'" class="ql-bold"></button>
              <button v-tooltip.bottom="'Italic'" class="ql-italic"></button>
              <button v-tooltip.bottom="'Underline'" class="ql-underline"></button>
            </span>
          </template>
        </Editor> -->

      </div>
      <div class="flex flex-col gap-1">
        <label for="multiple-ac-1" class="font-bold mb-2 block">Retorno da execução</label>
        <Textarea v-model="retorno_codigo" rows="5" cols="40" />
      </div>
      <div class="flex flex-col gap-1">
        <label for="multiple-ac-1" class="font-bold mb-2 block">Retorno do Teste</label>
        <Textarea v-model="retorno_teste" rows="5" cols="40" />
      </div>
      <Button type="submit" severity="secondary" label="Enviar" />
    </Form>
  </div>
</template>
  
<script setup lang="ts">
import * as Monaco from 'monaco-editor'

const optionsEditor = ref<Monaco.editor.IStandaloneEditorConstructionOptions>({
  theme: 'vs-dark'
});

import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import type { PrimeVueDataExercicio, RetornoExecucao } from '~/types'
import { reactive } from 'vue';
import { useToast } from 'primevue/usetoast';

const route = useRoute();
const toast = useToast();
const { $authService } = useNuxtApp();
let refQuillCodigoBase = ref(null)
const exercicioIDRep = ref('')
const tipoDaLinguagem = ref(0)

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
})

const retorno_codigo = ref('')
const retorno_teste = ref('')
const codigoBaseToSend = ref('')

const showFeedback = async (e: PrimeVueDataExercicio) => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/feedback_by_exercicio_id/'+e.id+'/'+config.public.METRICAS_FEEDBACK.REPETICAO, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    })

    if (response.ok) {
      const data = await response.json()
      toast.add({ severity: 'info', summary: 'Info', detail: data.feedback, life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Error', detail: "Error ao requisitar o feedback", life: 3000 });
    }
  } catch (error) {    
    console.log(error)
  }
};

const resolver = zodResolver(
  z.object({
    codigo_base: z.string().min(1, { message: 'Código base exercício é necessario.' })
  })
);

const onFormSubmit = async ({ valid }: { valid: boolean }) => {
  if (valid) {
    toast.add({
      severity: 'success',
      summary: 'Form foi enviado.',
      life: 3000
    });
  }

  if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
    codigoBaseToSend.value = exercicio.value.codigo_base
  } else {
    // @ts-ignore
    codigoBaseToSend.value = refQuillCodigoBase.value?.quill.getText()
  }

  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/discente_submit_exerc_rep_espacada', {
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
    })

    const exercicioResponse: RetornoExecucao = await response.json()
    if (response.ok) {
      retorno_codigo.value = exercicioResponse.out_put_from_user.output
      retorno_teste.value = exercicioResponse.out_put_unit_teste.output

      if (!exercicioResponse.out_put_unit_teste.success) {
        retorno_teste.value = exercicioResponse.out_put_unit_teste.error
      }
      if (exercicioResponse.out_put_from_user.success && exercicioResponse.out_put_unit_teste.success) {
        toast.add({ severity: 'success', summary: 'Exercício submetido com sucesso', life: 3000 });
      }
    } else {
      retorno_codigo.value = exercicioResponse.out_put_from_user.error
      retorno_teste.value = exercicioResponse.out_put_unit_teste.error
      toast.add({ severity: 'error', summary: 'Error ao submeter o exercício', life: 3000 });
    }
  } catch(error) {
    console.log(error)
  }
};

const loadRepeticaoEspacadaExercicio = async () => {
  // @ts-ignore
  exercicioIDRep.value = route.params.id
  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/discente_get_exerc_rep_espacada/'+exercicioIDRep.value, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    })

    if (response.ok) {
      const data = await response.json()
      exercicio.value = data.exercicio
      tipoDaLinguagem.value = data.tipo_da_linguagem
      toast.add({ severity: 'success', summary: 'Exercício carregado.', life: 3000 });
      await onSubmitMetricaExercicioInicio()
    } else {
      toast.add({ severity: 'error', summary: 'Error ao carregar o exercício.', life: 3000 });
    }
  } catch(error) {
    console.log(error)
  }
}

const onSubmitMetricaExercicioInicio = async () => {
  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/metricas_repeticao_espacada', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        exercicio_id: exercicio.value.id,
        tipo_metrica: useRuntimeConfig().public.METRICAS.INICIO
      })
    })

    if (response.ok) {
      console.log("Métrica da repetição enviada com sucesso")
    } else {
      console.log("Erro ao enviar a métrica da repetição")
    }
  } catch (error) {
    console.log(error)
  }
}

loadRepeticaoEspacadaExercicio()
</script>