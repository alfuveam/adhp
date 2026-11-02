<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />

    <div class="py-3 sm:py-5 border-b border-[#dce1d9]">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-5">
        <div class="space-y-1.5">
          <div class="inline-flex items-center gap-1.5 text-sm text-sky-800 font-semibold">
            <i class="pi pi-user text-[11px]"></i>
            <span>Seu espaço de prática</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-semibold text-slate-800 tracking-tight">
            Meus estudos
          </h1>
          <p class="text-sm sm:text-sm text-slate-600 max-w-xl leading-relaxed">
            Escolha uma lista para praticar ou retome um exercício nas suas revisões.
          </p>
        </div>

        <div class="flex items-center gap-3 shrink-0">
          <div
            class="px-4 py-2.5 text-center min-w-[80px]"
          >
            <div class="text-[11px] text-slate-600 font-medium">Trilhas</div>
            <div class="text-xl font-semibold text-slate-800">{{ trilhas?.length || 0 }}</div>
          </div>
          <div
            class="px-4 py-2.5 text-center min-w-[80px]"
          >
            <div class="text-[11px] text-sky-800 font-medium">Revisões agendadas</div>
            <div class="text-xl font-semibold text-sky-800">
              {{ exerciciosRepeticao?.length || 0 }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <div class="lg:col-span-7 space-y-4">
        <div class="flex items-center justify-between pb-1">
          <h2 class="text-base font-semibold text-slate-800 flex items-center gap-2">
            <i class="pi pi-folder text-sky-800 text-sm"></i>
            <span>Trilhas de exercícios</span>
          </h2>
          <span class="study-badge study-badge-purple text-[10px]">Passo a passo</span>
        </div>

        <div
          v-if="!trilhas || trilhas.length === 0"
          class="study-card p-10 text-center border-[#dce1d9] bg-[#ffffff]"
        >
          <p class="text-slate-600 text-sm">
            Suas trilhas aparecerão aqui quando forem cadastradas pelo professor.
          </p>
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="trilha in trilhas"
            :key="trilha.id"
            class="study-card p-6 border-[#dce1d9] bg-[#ffffff] space-y-5 shadow-none"
          >
            <div class="flex items-center justify-between pb-3 border-b border-[#dce1d9]">
              <div class="flex items-center gap-3.5">
                <span
                  class="w-10 h-10 rounded-lg flex items-center justify-center font-semibold text-sm"
                  :class="
                    trilha.tipo_da_linguagem === 1
                      ? 'bg-cyan-500/20 text-sky-800 border border-cyan-500/30'
                      : 'bg-amber-400/20 text-amber-800 border border-amber-400/30'
                  "
                >
                  {{ trilha.tipo_da_linguagem === 1 ? "Go" : "Py" }}
                </span>
                <div>
                  <h3 class="text-lg font-semibold text-slate-800">{{ trilha.name }}</h3>
                  <span class="text-sm text-slate-600 font-medium">
                    {{ trilha.tipo_da_linguagem === 1 ? "Golang" : "Python" }} •
                    {{ trilha.listas?.length || 0 }} Listas de Exercícios
                  </span>
                </div>
              </div>
            </div>

            <div class="space-y-2">
              <Accordion :value="['0']" multiple>
                <AccordionPanel v-for="lista in trilha.listas" :key="lista.id" :value="lista.id">
                  <AccordionHeader>
                    <div class="flex items-center justify-between w-full pr-3">
                      <div class="flex items-center gap-2.5 text-sm font-semibold text-slate-800">
                        <i class="pi pi-list text-sky-800 text-sm"></i>
                        <span>{{ lista.name }}</span>
                      </div>
                      <span
                        class="text-sm px-2.5 py-0.5 rounded-full bg-[#ffffff] text-slate-600 border border-[#dce1d9] font-medium"
                      >
                        {{ lista.exercicios?.length || 0 }} exercícios
                      </span>
                    </div>
                  </AccordionHeader>

                  <AccordionContent>
                    <div class="p-3 space-y-4">
                      <TreeTable :value="lista.exercicios" class="w-full text-sm">
                        <Column field="name" header="Exercício" expander>
                          <template #body="{ node }">
                            <div class="flex items-center gap-2 py-1.5">
                              <i
                                :class="
                                  node.data.habilitado
                                    ? 'pi pi-check text-emerald-700 font-semibold'
                                    : 'pi pi-lock text-slate-500'
                                "
                                class="text-sm"
                              ></i>
                              <span
                                :class="
                                  node.data.habilitado
                                    ? 'text-slate-800 font-medium'
                                    : 'text-slate-500'
                                "
                              >
                                {{ node.data.name }}
                              </span>
                            </div>
                          </template>
                        </Column>

                        <Column style="width: 7.5rem; text-align: right">
                          <template #body="{ node }">
                            <div class="flex justify-end">
                              <NuxtLink
                                v-if="node.data.habilitado"
                                :to="`/dashboard_discente/exercicio/lista/${lista.id}/${
                                  node.data.order_index - 1
                                }/${trilha.tipo_da_linguagem}`"
                              >
                                <button class="study-btn-primary px-3 py-1.5 text-sm">
                                  Resolver
                                </button>
                              </NuxtLink>
                              <span v-else class="text-sm text-slate-500 font-medium">
                                Bloqueado
                              </span>
                            </div>
                          </template>
                        </Column>
                      </TreeTable>

                      <div class="pt-2 border-t border-[#dce1d9] flex justify-end">
                        <NuxtLink
                          :to="`/dashboard_discente/exercicio/lista/${lista.id}/0/${trilha.tipo_da_linguagem}`"
                        >
                          <button
                            class="study-btn-secondary px-3.5 py-1.5 text-sm flex items-center gap-1.5"
                          >
                            <i class="pi pi-play text-[10px]"></i>
                            <span>Começar esta lista</span>
                          </button>
                        </NuxtLink>
                      </div>
                    </div>
                  </AccordionContent>
                </AccordionPanel>
              </Accordion>
            </div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-5 space-y-5">
        <div class="study-card p-6 border-[#dce1d9] bg-[#ffffff] space-y-4 shadow-none">
          <div class="flex items-center justify-between pb-3 border-b border-[#dce1d9]">
            <h2 class="text-base font-semibold text-slate-800 flex items-center gap-2">
              <i class="pi pi-history text-sky-800 text-sm"></i>
              <span>Exercícios para revisar</span>
            </h2>
          </div>

          <div
            v-if="!exerciciosRepeticao || exerciciosRepeticao.length === 0"
            class="py-8 text-center text-slate-600 space-y-1.5"
          >
            <div
              class="w-10 h-10 mx-auto rounded-full bg-sky-900/5 text-sky-800 flex items-center justify-center mb-2"
            >
              <i class="pi pi-check text-base"></i>
            </div>
            <p class="text-sm font-semibold text-slate-800">Nenhuma revisão pendente no momento</p>
            <p class="text-sm text-slate-600">
              Conclua exercícios das trilhas para agendar novas repetições.
            </p>
          </div>

          <div v-else class="space-y-3 max-h-[420px] overflow-y-auto pr-1">
            <div
              v-for="rep in exerciciosRepeticao"
              :key="rep.id"
              class="p-4 rounded-lg bg-[#ffffff] border border-[#dce1d9] space-y-2.5"
            >
              <div class="flex items-start justify-between gap-2">
                <span class="text-sm font-semibold text-slate-800 leading-snug">
                  {{ rep.descricao_exercicio }}
                </span>
                <span
                  class="shrink-0 text-[10px] font-semibold px-2.5 py-0.5 rounded-full"
                  :class="
                    verificaSePodeRepetir(rep)
                      ? 'bg-amber-400/15 text-amber-800 border border-amber-400/30'
                      : 'bg-emerald-500/15 text-emerald-800 border border-emerald-500/30'
                  "
                >
                  {{ verificaSePodeRepetir(rep) ? "Aguardando" : "Liberado" }}
                </span>
              </div>

              <div class="flex items-center justify-between text-[11px] text-slate-600 pt-1">
                <span class="font-mono text-sky-800">Próxima: {{ rep.proxima_repeticao }}</span>

                <NuxtLink :to="`/dashboard_discente/exercicio/repetir/${rep.id}`">
                  <button
                    :disabled="verificaSePodeRepetir(rep)"
                    class="transition"
                    :class="
                      verificaSePodeRepetir(rep)
                        ? 'px-3 py-1 rounded-lg text-sm font-semibold bg-[#ffffff] text-slate-500 border border-[#dce1d9] cursor-not-allowed'
                        : 'study-btn-emerald px-3.5 py-1 text-sm'
                    "
                  >
                    Revisar
                  </button>
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>

        <details class="study-card p-5 border-[#dce1d9] bg-[#ffffff] space-y-3">
          <summary class="cursor-pointer text-sm font-semibold text-slate-800">
            Como funcionam minhas revisões
          </summary>
          <p class="text-sm text-slate-600">
            Estes são os intervalos previstos entre as revisões dos exercícios.
          </p>
          <div class="grid grid-cols-5 gap-2 text-center text-sm font-semibold">
            <div
              class="p-2 rounded-lg border transition"
              :class="
                v_mode_tempo_da_repeticao === '1'
                  ? 'bg-[#285b75] border-[#285b75] text-white'
                  : 'bg-[#ffffff] border-[#dce1d9] text-slate-600'
              "
            >
              1h
            </div>
            <div
              class="p-2 rounded-lg border transition"
              :class="
                v_mode_tempo_da_repeticao === '2'
                  ? 'bg-[#285b75] border-[#285b75] text-white'
                  : 'bg-[#ffffff] border-[#dce1d9] text-slate-600'
              "
            >
              9h
            </div>
            <div
              class="p-2 rounded-lg border transition"
              :class="
                v_mode_tempo_da_repeticao === '3'
                  ? 'bg-[#285b75] border-[#285b75] text-white'
                  : 'bg-[#ffffff] border-[#dce1d9] text-slate-600'
              "
            >
              1d
            </div>
            <div
              class="p-2 rounded-lg border transition"
              :class="
                v_mode_tempo_da_repeticao === '4'
                  ? 'bg-[#285b75] border-[#285b75] text-white'
                  : 'bg-[#ffffff] border-[#dce1d9] text-slate-600'
              "
            >
              6d
            </div>
            <div
              class="p-2 rounded-lg border transition"
              :class="
                v_mode_tempo_da_repeticao === '5'
                  ? 'bg-[#285b75] border-[#285b75] text-white'
                  : 'bg-[#ffffff] border-[#dce1d9] text-slate-600'
              "
            >
              31d
            </div>
          </div>
        </details>

        <details class="study-card p-5 border-[#dce1d9] bg-[#ffffff] space-y-2">
          <summary class="cursor-pointer text-sm font-semibold">
            Entenda a repetição espaçada
          </summary>
          <p class="text-sm text-slate-600">
            O gráfico ilustra uma curva teórica de memória; ele não mede seu desempenho.
          </p>
          <ChartRepeticaoEspacada />
        </details>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
definePageMeta({
    middleware: 'auth-discente'
});

import { ref } from 'vue';
import type { Trilha } from '~/types';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { useToast } from 'primevue/usetoast';
import ChartRepeticaoEspacada from '@/components/ChartRepeticaoEspacada.vue';

const trilhas = ref<Trilha[]>([]);
const { $authService } = useNuxtApp();

interface ExerciciosRepeticao {
  id: string;
  descricao_exercicio: string;
  ultima_repeticao: string;
  proxima_repeticao: string;
  id_exercicio: string;
}

const exerciciosRepeticao = ref<ExerciciosRepeticao[]>([]);
const toast = useToast();
const v_mode_tempo_da_repeticao = ref('');

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
    // Optional submission endpoint
  }
};

const getTrilhasListaExercicios = async () => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/dashboard_discente', {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    });

    if (response.ok) {
      const resData = await response.json();
      trilhas.value = resData.trilhas;
      if (resData.id_tempo_repeticao) {
        v_mode_tempo_da_repeticao.value = resData.id_tempo_repeticao.toString();
        initialValues.value.tempo_da_repeticao = resData.id_tempo_repeticao.toString();
      }
    } else {
      toast.add({ severity: 'error', summary: 'Erro ao carregar trilhas', detail: 'Falha ao buscar dados do painel.', life: 3000 });
    }
  } catch (error) {
    console.error(error);
  }
};

const getExerciciosRepeticao = async () => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/exercicios_repeticao_by_user', {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    });

    if (response.ok) {
      exerciciosRepeticao.value = await response.json();
    }
  } catch (error) {
    console.log(error);
  }
};

const verificaSePodeRepetir = (repeticao: ExerciciosRepeticao) => {
  if (!repeticao.proxima_repeticao) return false;
  const date = new Date(repeticao.proxima_repeticao);
  const now = new Date();
  return date.getTime() > now.getTime();
};

getTrilhasListaExercicios();
getExerciciosRepeticao();
</script>