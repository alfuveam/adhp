<template>
  <div class="card flex justify-center">    
    <ConfirmDialog></ConfirmDialog>

    <Form v-slot="$form" :initialValues :resolver @submit="onFormSubmit" class="flex flex-col gap-4 p-4 w-screen sm:w-screen">
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
      <Button type="submit" severity="secondary" :label="callBackSubmitText" />
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
const route = useRoute()

const toast = useToast();
const { $authService } = useNuxtApp();
let refQuillCodigoBase = ref(null)
const exercicioID = ref('')

const tipoDaLinguagem = ref(0)
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
})

const emit = defineEmits(['show-toast', 'call-back-atualizar-lista']);

const retorno_codigo = ref('')
const retorno_teste = ref('')
const codigoBaseToSend = ref('')

const showFeedback = (e: PrimeVueDataExercicio) => {
  emit('show-toast', e);
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
    codigoBaseToSend.value = props.exercicio.codigo_base
  } else {
    // @ts-ignore
    codigoBaseToSend.value = refQuillCodigoBase.value?.quill.getText()
  }

  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/discente_submit_exercicio', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        exercicio_id: props.exercicio.id,
        codigo_base: btoa(encodeURIComponent(codigoBaseToSend.value))
      })
    })

    const exercicioResponse: RetornoExecucao = await response.json()
    if (response.ok) {
      retorno_codigo.value = exercicioResponse.out_put_from_user.output
      retorno_teste.value = exercicioResponse.out_put_unit_teste.output

      if (exercicioResponse.out_put_from_user.success && exercicioResponse.out_put_unit_teste.success) {
        toast.add({ severity: 'success', summary: 'Exercício submetido com sucesso', life: 3000 });
        emit('call-back-atualizar-lista')
      } else {
        retorno_codigo.value = exercicioResponse.out_put_from_user.error
        retorno_teste.value = exercicioResponse.out_put_unit_teste.error
        toast.add({ severity: 'error', summary: 'Error ao submeter o exercício', life: 3000 });
      }

    } else {
      toast.add({ severity: 'error', summary: 'Error ao submeter o exercício', life: 3000 });
    }
  } catch(error) {
    console.log(error)
  }
};

// -_-
watch(async () => props.exercicio, async (newExercicio) => {
  try {
    const exercicioProps = await newExercicio
    if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
      props.exercicio.codigo_base = decodeURIComponent(atob(props.exercicio.codigo_base))
    } else {
      const v = decodeURIComponent(atob(exercicioProps.codigo_base))
      // @ts-ignore
      refQuillCodigoBase.value?.quill.setText(v)
    }
    if (exercicioID.value != '' && exercicioID.value != exercicioProps.id) {
      retorno_codigo.value = ""
      retorno_teste.value = ""
    }

    exercicioID.value = exercicioProps.id
  } catch (error) {
    console.log(error)
  }
})

onMounted(() => {
  tipoDaLinguagem.value = Number(route.params.id[2])
  const intervalId = setInterval(checkIfQuillIsLoad,100)

  function checkIfQuillIsLoad() {
    if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
      clearInterval(intervalId);
      const v = decodeURIComponent(atob(props.exercicio.codigo_base))
      props.exercicio.codigo_base = v
    } else {
      // @ts-ignore
      if (refQuillCodigoBase && refQuillCodigoBase.value?.quill.root) {
        clearInterval(intervalId);

        const v = decodeURIComponent(atob(props.exercicio.codigo_base))
        // @ts-ignore
        refQuillCodigoBase.value?.quill.setText(v)
      }
    }
  }
})
</script>