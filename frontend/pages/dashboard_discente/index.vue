<template>
<div class="card">
  <div class="mb-8 flex flex-col lg:flex-row w-full justify-center">
    <Panel class="bg-gray-700 w-auto h-auto min-w-[200px] lg:min-w-[650px]">
      <div class="card">
        <template v-for="trilha in trilhas">
          <div class="card">
            <h1 class="mb-4 text-2xl font-extrabold text-center">{{trilha.name}}</h1>
            <Accordion :value="['0']" multiple v-for="lista in trilha.listas">
                <AccordionPanel :value="lista.id">
                <AccordionHeader>
                    {{lista.name}}
                </AccordionHeader>
                <AccordionContent>
                    <div class="card">
                        <TreeTable :value="lista.exercicios" tableStyle="">
                        <template #header>
                            <div class="text-xl font-bold">{{lista.name}}</div>
                        </template>
                        <Column field="name" header="Name" expander style="">
                            <div class="mt-2">
                            <InputText placeholder="Enter new value" />
                            </div>
                        </Column>
                        <Column style="">
                            <template #body="{ node }">
                            <div class="flex flex-wrap gap-2">
                              <NuxtLink :to="`/dashboard_discente/exercicio/lista/${lista.id}/${node.data.order_index-1}/${trilha.tipo_da_linguagem}`">
                                <Button label="Responder" :disabled="!node.data.habilitado"/>
                              </NuxtLink>
                            </div>
                            </template>
                        </Column>
                        <template #footer>
                            <div class="flex justify-start">
                            <NuxtLink :to="`/dashboard_discente/exercicio/lista/${lista.id}/0/${trilha.tipo_da_linguagem}`">
                                <Button icon="pi pi-refresh" label="Responder lista" severity="warn"  />
                            </NuxtLink>
                            </div>
                        </template>
                        </TreeTable>
                    </div>
                </AccordionContent>
                </AccordionPanel>
            </Accordion>
          </div>
          <br>
        </template>
      </div>
    </Panel>
    <div class="h-[1300px] w-auto h-auto min-w-[200px] lg:min-w-[650px]">
      <Panel class="flex items-top justify-center h-[600px]">
        <h2 class="text-center mt-2">Repetição Espaçada</h2>
        <div class="card overflow-y-scroll max-h-[540px]">
          <DataTable :value="exerciciosRepeticao" paginator :rows="5" :rowsPerPageOptions="[5, 10, 20, 50]">
            <Column field="descricao_exercicio">
              <h2 field="descricao_exercicio"></h2>
              <template #body="{ data }">
                <div>{{ data.descricao_exercicio }}</div>    
                <br>
                <div class="flex flex-row mt-1 justify-center">
                  <div class="flex flex-row mt-1 justify-center">
                    <div class="mr-3" field="ultima_repeticao">Última repetição: <br> {{ data.ultima_repeticao }}</div>
                    <div class="mr-3" field="proxima_repeticao">Proxima repetição: <br> {{ data.proxima_repeticao }}</div>
                  </div>
                  <NuxtLink :to="`/dashboard_discente/exercicio/repetir/${data.id}`">
                    <Button label="Repetir" :disabled="verificaSePodeRepetir(data)"></Button>
                  </NuxtLink>
                </div>
              </template>
            </Column>
          </DataTable>
          <Toast />
        </div>
      </Panel>
      <div class="card flex flex-col justify-center mt-1 p-2 items-center">
        <h2 class="text-center p-2">Tempo da Repetição Espaçada</h2>
        <Form v-slot="$form" :resolver="resolver" :initialValues="initialValues" @submit="onFormSubmit" class="flex flex-col items-center gap-4">
          <div class="flex flex-col gap-2">
            <RadioButtonGroup name="tempo_da_repeticao" v-model="v_mode_tempo_da_repeticao" class="flex flex-wrap gap-4">
                <div class="flex items-center gap-2">
                    <RadioButton inputId="1" value="1" />
                    <label for="1">1 Hora</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton inputId="2" value="2" />
                    <label for="2">9 Horas</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton inputId="3" value="3" />
                    <label for="3">1 Dia</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton inputId="4" value="4" />
                    <label for="4">6 Dias</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton inputId="5" value="5" />
                    <label for="5">31 Dias</label>
                </div>
            </RadioButtonGroup>
            <Message v-if="$form.tempo_da_repeticao?.invalid" severity="error" size="small" variant="simple">{{ $form.tempo_da_repeticao.error?.message }}</Message>
          </div>
          <!-- <Button class="w-28" type="submit" severity="secondary" label="Submit" /> -->
        </Form>
        <Divider />
      </div>
      <ChartRepeticaoEspacada />
    </div>
  </div>
</div>
</template>

<script setup lang="ts">
definePageMeta({
    middleware: 'auth-discente'
})

import type { Trilha } from '~/types';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { useToast } from 'primevue/usetoast';
import ChartRepeticaoEspacada from '@/components/ChartRepeticaoEspacada.vue';

const trilhas = ref<Trilha[]>([]);
const { $authService } = useNuxtApp();

interface ExerciciosRepeticao {
  id: string,
  descricao_exercicio: string,
  ultima_repeticao: string,
  proxima_repeticao: string,
  id_exercicio: string,
}

const exerciciosRepeticao = ref<ExerciciosRepeticao[]>();
const toast = useToast();
const v_mode_tempo_da_repeticao = ref('')

const initialValues = ref({
  tempo_da_repeticao: ''
});

const resolver = ref(zodResolver(
  z.object({
    tempo_da_repeticao: z.string().min(1, { message: 'Tempo da repetição.' })
  })
));

const onFormSubmit = async ({ valid }: { valid: boolean }) => {
  if (valid) {
    // try {
    //   const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/discente_submit_tempo_repeticao/'+v_mode_tempo_da_repeticao.value, {
    //     method: "PUT",
    //     headers: {
    //       "Content-Type": "application/json",
    //       "Authorization": "Bearer " + $authService.getToken()
    //     },
    //   })

    //   if (response.ok) {
    //     toast.add({ severity: 'success', summary: 'Tempo da repetição alterado.', life: 3000 });
    //   } else {
    //     toast.add({ severity: 'error', summary: 'Falha ao alterar o tempo da repetição.', life: 3000 });
    //   }
    // } catch (error) {
    //   console.log(error)
    // }
  }
};

const getTrilhasListaExercicios = async () => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/dashboard_discente', {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    })

    if (response.ok) {
      const resData = await response.json()
      trilhas.value = resData.trilhas
      v_mode_tempo_da_repeticao.value = resData.id_tempo_repeticao.toString()
      initialValues.value.tempo_da_repeticao = resData.id_tempo_repeticao.toString()
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao carregar o dashboard.', life: 3000 });
    }
  } catch(error) {
    console.error(error)
  }
}

const getExerciciosRepeticao = async () => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/exercicios_repeticao_by_user', {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    })

    if (response.ok) {
      exerciciosRepeticao.value = await response.json()
    }
  } catch(error) {
    console.log(error)
  }
}

const verificaSePodeRepetir = (repeticao: ExerciciosRepeticao) => {
  const date = new Date(repeticao.proxima_repeticao);
  const now = new Date();

  return date.getTime() > now.getTime();
}

getTrilhasListaExercicios()
getExerciciosRepeticao()
</script>