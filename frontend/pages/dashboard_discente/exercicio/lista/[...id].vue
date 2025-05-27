<template>
  <Toast />
  <div class="card flex justify-center">
    <Stepper :value="currentStep" orientation="vertical" class="flex flex-col sm:flex-row w-full ">
      <div class="bg-zinc-800 overflow-y-auto overflow-x-auto w-full h-44 sm:h-[600px]" style="flex-direction: column">
        <Step class="text-center" v-for="(exercicio, index) in lista?.exercicios" 
          :key="index" 
          :value="index" 
          @click="exercicioSelecionado(exercicio.data);" 
          :disabled="!exercicio.data.habilitado"
          >
          <div class="card max-w-[600px]">
            <Panel class="h-auto">
              <article class="text-wrap max-w-auto" style="text-wrap: wrap;">
                <p class="w-full">
                  {{exercicio.data.name}}
                </p>
              </article>
            </Panel>
          </div>
        </Step>
      </div>
      <StepPanels class="w-auto sm:w-[496px] lg:w-[696px]">
        <ExercicioQuestao class="w-auto sm:w-[496px] lg:w-[696px]" :exercicio="exercicioSelecionada" @call-back-atualizar-lista="atualizarLista" callBackSubmitText="Submeter Código"  @show-toast="showToastFeedBack"/>
      </StepPanels>
    </Stepper>
  </div>
</template>

<script setup lang="ts">
import type { Listas, PrimeVueDataExercicio } from '~/types'
import { useToast } from 'primevue/usetoast';
const toast = useToast();
const route = useRoute()
const { $authService } = useNuxtApp()

const currentStep = ref(0)
const idLista = ref('')

const showToastFeedBack = async (e: PrimeVueDataExercicio) => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/feedback_by_exercicio_id/'+e.id+'/'+config.public.METRICAS_FEEDBACK.EXERCICIO, {
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

const lista = ref<Listas>({
  id: '',
  name: '',
  exercicios: [],
  isEditing: false
})

const exercicioSelecionada = ref<PrimeVueDataExercicio>({
  name: '',
  id: '',
  inputValue: '',
  type: 0,
  codigo_rodou: false,
  codigo_base: '',
  habilitado: false,
  order_index: 0,
})

const exercicioSelecionado = (exercicio: PrimeVueDataExercicio) => {
  exercicioSelecionada.value = exercicio
  onSubmitMetricaExercicioInicio(exercicio)
}

const atualizarLista = async () => {
  await getExerciciosHabilitadosByLista()
}

const getExerciciosHabilitadosByLista = async () => {
  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/exercicios_habilitados_by_lista/' + idLista.value, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer ' + $authService.getToken()
      }
    })

    if (response.ok) {
      const resData = await response.json()
      lista.value = resData[0]
      exercicioSelecionada.value = resData[0].exercicios[currentStep.value].data
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao carregar a lista.', life: 3000 });
    }
  } catch (error) {
    console.log(error)
  }
}

const onSubmitMetricaExercicioInicio = async (exercicio: PrimeVueDataExercicio) => {
  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/metricas_exercicio', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        exercicio_id: exercicio.id,
        tipo_metrica: useRuntimeConfig().public.METRICAS.INICIO
      })
    })

    if (response.ok) {
      console.log("Métrica enviada com sucesso")
    } else {
      console.log("Erro ao enviar a métrica")
    }
  } catch (error) {
    console.log(error)
  }
}

onMounted( async () => {
  idLista.value = route.params.id[0]
  currentStep.value = Number(route.params.id[1])

  await getExerciciosHabilitadosByLista()
})

</script>