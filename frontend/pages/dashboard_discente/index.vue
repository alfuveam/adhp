<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />

    <!-- STUDENT WELCOME BANNER (Mimo Style) -->
    <div class="mimo-card p-6 sm:p-7 border-[#363e66] bg-[#1c2038] shadow-xl">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-5">
        <div class="space-y-1.5">
          <div class="inline-flex items-center gap-1.5 text-xs text-purple-300 font-bold uppercase tracking-wider">
            <i class="pi pi-user text-[11px]"></i>
            <span>Painel do Estudante</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-extrabold text-white tracking-tight">
            Minhas Trilhas e Repetições
          </h1>
          <p class="text-xs sm:text-sm text-slate-300 max-w-xl leading-relaxed">
            Acompanhe o andamento das suas listas de exercícios e revise os algoritmos agendados na fila de repetição.
          </p>
        </div>

        <!-- Quick Summary Stats (Mimo Style) -->
        <div class="flex items-center gap-3 shrink-0">
          <div class="px-4 py-2.5 rounded-2xl bg-[#141728] border border-[#2d3458] text-center min-w-[95px] shadow-sm">
            <div class="text-[11px] text-slate-400 font-medium">Trilhas</div>
            <div class="text-xl font-extrabold text-white">{{ trilhas?.length || 0 }}</div>
          </div>
          <div class="px-4 py-2.5 rounded-2xl bg-[#141728] border border-[#2d3458] text-center min-w-[95px] shadow-sm">
            <div class="text-[11px] text-purple-300 font-medium">Revisões</div>
            <div class="text-xl font-extrabold text-purple-400">{{ exerciciosRepeticao?.length || 0 }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- MAIN TWO-COLUMN DASHBOARD GRID -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      
      <!-- LEFT COLUMN: TRILHAS E LISTAS DE EXERCÍCIOS (Cols 1-7) -->
      <div class="lg:col-span-7 space-y-4">
        <div class="flex items-center justify-between pb-1">
          <h2 class="text-base font-bold text-white flex items-center gap-2">
            <i class="pi pi-folder text-purple-400 text-sm"></i>
            <span>Trilhas Disponíveis</span>
          </h2>
          <span class="mimo-badge mimo-badge-purple text-[10px]">Sequencial</span>
        </div>

        <!-- Empty state if no trilhas -->
        <div v-if="!trilhas || trilhas.length === 0" class="mimo-card p-10 text-center border-[#363e66] bg-[#1c2038]">
          <p class="text-slate-300 text-sm">Nenhuma trilha cadastrada no momento.</p>
        </div>

        <!-- Trilhas List -->
        <div v-else class="space-y-4">
          <div 
            v-for="trilha in trilhas" 
            :key="trilha.id" 
            class="mimo-card p-6 border-[#363e66] bg-[#1c2038] space-y-5 shadow-lg"
          >
            <!-- Trilha Header -->
            <div class="flex items-center justify-between pb-3 border-b border-[#2d3458]">
              <div class="flex items-center gap-3.5">
                <span 
                  class="w-10 h-10 rounded-xl flex items-center justify-center font-extrabold text-sm"
                  :class="trilha.tipo_da_linguagem === 1 ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/30' : 'bg-amber-400/20 text-amber-300 border border-amber-400/30'"
                >
                  {{ trilha.tipo_da_linguagem === 1 ? 'Go' : 'Py' }}
                </span>
                <div>
                  <h3 class="text-lg font-extrabold text-white">{{ trilha.name }}</h3>
                  <span class="text-xs text-slate-400 font-medium">
                    {{ trilha.tipo_da_linguagem === 1 ? 'Golang' : 'Python' }} • {{ trilha.listas?.length || 0 }} Listas de Exercícios
                  </span>
                </div>
              </div>
            </div>

            <!-- Listas Accordion -->
            <div class="space-y-2">
              <Accordion :value="['0']" multiple>
                <AccordionPanel v-for="lista in trilha.listas" :key="lista.id" :value="lista.id">
                  <AccordionHeader>
                    <div class="flex items-center justify-between w-full pr-3">
                      <div class="flex items-center gap-2.5 text-sm font-bold text-slate-100">
                        <i class="pi pi-list text-purple-400 text-xs"></i>
                        <span>{{ lista.name }}</span>
                      </div>
                      <span class="text-xs px-2.5 py-0.5 rounded-full bg-[#151729] text-slate-300 border border-[#2d3458] font-medium">
                        {{ lista.exercicios?.length || 0 }} exercícios
                      </span>
                    </div>
                  </AccordionHeader>

                  <AccordionContent>
                    <div class="p-3 space-y-4">
                      <!-- TreeTable of Exercicios -->
                      <TreeTable :value="lista.exercicios" class="w-full text-xs">
                        <Column field="name" header="Exercício" expander>
                          <template #body="{ node }">
                            <div class="flex items-center gap-2 py-1.5">
                              <i 
                                :class="node.data.habilitado ? 'pi pi-check text-emerald-400 font-bold' : 'pi pi-lock text-slate-500'" 
                                class="text-xs"
                              ></i>
                              <span :class="node.data.habilitado ? 'text-slate-100 font-medium' : 'text-slate-500'">
                                {{ node.data.name }}
                              </span>
                            </div>
                          </template>
                        </Column>

                        <Column style="width: 7.5rem; text-align: right;">
                          <template #body="{ node }">
                            <div class="flex justify-end">
                              <NuxtLink 
                                v-if="node.data.habilitado"
                                :to="`/dashboard_discente/exercicio/lista/${lista.id}/${node.data.order_index-1}/${trilha.tipo_da_linguagem}`"
                              >
                                <button class="mimo-btn-primary px-3 py-1.5 text-xs">
                                  Resolver
                                </button>
                              </NuxtLink>
                              <span v-else class="text-xs text-slate-500 font-medium">
                                Bloqueado
                              </span>
                            </div>
                          </template>
                        </Column>
                      </TreeTable>

                      <!-- Action to start list from beginning -->
                      <div class="pt-2 border-t border-[#2d3458] flex justify-end">
                        <NuxtLink :to="`/dashboard_discente/exercicio/lista/${lista.id}/0/${trilha.tipo_da_linguagem}`">
                          <button class="mimo-btn-secondary px-3.5 py-1.5 text-xs flex items-center gap-1.5">
                            <i class="pi pi-play text-[10px]"></i>
                            <span>Iniciar do Primeiro</span>
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

      <!-- RIGHT COLUMN: REPETIÇÃO ESPAÇADA & RETENÇÃO (Cols 8-12) -->
      <div class="lg:col-span-5 space-y-5">
        
        <!-- Repetição Espaçada Deck -->
        <div class="mimo-card p-6 border-[#363e66] bg-[#1c2038] space-y-4 shadow-xl">
          <div class="flex items-center justify-between pb-3 border-b border-[#2d3458]">
            <h2 class="text-base font-extrabold text-white flex items-center gap-2">
              <i class="pi pi-history text-purple-400 text-sm"></i>
              <span>Fila de Repetição Espaçada</span>
            </h2>
            <span class="mimo-badge mimo-badge-purple text-[10px]">Ebbinghaus</span>
          </div>

          <!-- Empty queue message -->
          <div v-if="!exerciciosRepeticao || exerciciosRepeticao.length === 0" class="py-8 text-center text-slate-400 space-y-1.5">
            <div class="w-10 h-10 mx-auto rounded-full bg-purple-500/10 text-purple-300 flex items-center justify-center mb-2">
              <i class="pi pi-check text-base"></i>
            </div>
            <p class="text-sm font-bold text-white">Nenhuma revisão pendente no momento</p>
            <p class="text-xs text-slate-400">Conclua exercícios das trilhas para agendar novas repetições.</p>
          </div>

          <!-- Review items list -->
          <div v-else class="space-y-3 max-h-[420px] overflow-y-auto pr-1">
            <div 
              v-for="rep in exerciciosRepeticao" 
              :key="rep.id" 
              class="p-4 rounded-xl bg-[#141728] border border-[#2d3458] space-y-2.5"
            >
              <div class="flex items-start justify-between gap-2">
                <span class="text-xs font-bold text-slate-100 leading-snug">
                  {{ rep.descricao_exercicio }}
                </span>
                <span 
                  class="shrink-0 text-[10px] font-bold px-2.5 py-0.5 rounded-full"
                  :class="verificaSePodeRepetir(rep) ? 'bg-amber-400/15 text-amber-300 border border-amber-400/30' : 'bg-emerald-500/15 text-emerald-300 border border-emerald-500/30'"
                >
                  {{ verificaSePodeRepetir(rep) ? 'Aguardando' : 'Liberado' }}
                </span>
              </div>

              <div class="flex items-center justify-between text-[11px] text-slate-400 pt-1">
                <span class="font-mono text-purple-300">Próxima: {{ rep.proxima_repeticao }}</span>

                <NuxtLink :to="`/dashboard_discente/exercicio/repetir/${rep.id}`">
                  <button 
                    :disabled="verificaSePodeRepetir(rep)"
                    class="transition"
                    :class="verificaSePodeRepetir(rep) ? 'px-3 py-1 rounded-lg text-xs font-semibold bg-[#202540] text-slate-500 border border-[#2d3458] cursor-not-allowed' : 'mimo-btn-emerald px-3.5 py-1 text-xs'"
                  >
                    Repetir
                  </button>
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>

        <!-- Intervalo de Repetição Configuração -->
        <div class="mimo-card p-5 border-[#363e66] bg-[#1c2038] space-y-3">
          <div class="flex items-center justify-between">
            <h3 class="text-xs font-bold uppercase tracking-wider text-slate-300">Meu Intervalo Atual</h3>
            <span class="text-[11px] text-purple-300 font-mono">Ciclo Ativo</span>
          </div>
          <div class="grid grid-cols-5 gap-2 text-center text-xs font-bold">
            <div class="p-2 rounded-xl border transition" :class="v_mode_tempo_da_repeticao === '1' ? 'bg-[#7c4ed4] border-[#9061ea] text-white shadow-[0_2px_0_0_#5326ab]' : 'bg-[#141728] border-[#2d3458] text-slate-400'">1h</div>
            <div class="p-2 rounded-xl border transition" :class="v_mode_tempo_da_repeticao === '2' ? 'bg-[#7c4ed4] border-[#9061ea] text-white shadow-[0_2px_0_0_#5326ab]' : 'bg-[#141728] border-[#2d3458] text-slate-400'">9h</div>
            <div class="p-2 rounded-xl border transition" :class="v_mode_tempo_da_repeticao === '3' ? 'bg-[#7c4ed4] border-[#9061ea] text-white shadow-[0_2px_0_0_#5326ab]' : 'bg-[#141728] border-[#2d3458] text-slate-400'">1d</div>
            <div class="p-2 rounded-xl border transition" :class="v_mode_tempo_da_repeticao === '4' ? 'bg-[#7c4ed4] border-[#9061ea] text-white shadow-[0_2px_0_0_#5326ab]' : 'bg-[#141728] border-[#2d3458] text-slate-400'">6d</div>
            <div class="p-2 rounded-xl border transition" :class="v_mode_tempo_da_repeticao === '5' ? 'bg-[#7c4ed4] border-[#9061ea] text-white shadow-[0_2px_0_0_#5326ab]' : 'bg-[#141728] border-[#2d3458] text-slate-400'">31d</div>
          </div>
        </div>

        <!-- Curva de Retenção Chart Component -->
        <div class="mimo-card p-5 border-[#363e66] bg-[#1c2038] space-y-2">
          <div class="flex items-center justify-between text-xs">
            <span class="font-bold text-white">Curva de Retenção Teórica</span>
            <span class="text-purple-300 font-mono text-[11px]">Ebbinghaus</span>
          </div>
          <ChartRepeticaoEspacada />
        </div>

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