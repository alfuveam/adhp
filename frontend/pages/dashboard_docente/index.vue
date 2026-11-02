<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
    <Toast />
    <ConfirmDialog />

    <ManipuladorTrilha
      v-model="visibleFormManipulador"
      :isUpdate="isUpdateDialog"
      :trilhaID="trilhaID"
      @some-event="getTrilhasListaExercicios"
      @blur="formManipulador()"
    />

    <div class="study-card p-6 sm:p-7 border-[#dce1d9] bg-[#ffffff] shadow-none">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-5">
        <div class="space-y-1.5">
          <div class="inline-flex items-center gap-1.5 text-sm text-sky-800 font-semibold">
            <i class="pi pi-user-edit text-[11px]"></i>
            <span>Painel do Docente / Instrutor</span>
          </div>
          <h1 class="text-xl sm:text-2xl font-semibold text-slate-800 tracking-tight">
            Gestão de Trilhas e Listas de Exercícios
          </h1>
          <p class="text-sm sm:text-sm text-slate-600 max-w-xl leading-relaxed">
            Crie novos conteúdos, estruture as listas de exercícios sequenciais e configure a suíte
            de testes unitários.
          </p>
        </div>

        <button
          @click="formManipulador()"
          class="study-btn-primary px-5 py-2.5 text-sm font-semibold flex items-center gap-2"
        >
          <i class="pi pi-plus text-sm"></i>
          <span>Nova Trilha</span>
        </button>
      </div>
    </div>

    <div
      v-if="!trilhas || trilhas.length === 0"
      class="study-card p-12 text-center border border-slate-200"
    >
      <i class="pi pi-folder-open text-4xl text-slate-500 mb-3 block"></i>
      <p class="text-slate-600 text-sm mb-4">Nenhuma trilha cadastrada ainda.</p>
      <button
        @click="formManipulador()"
        class="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-semibold text-slate-800 bg-blue-600 hover:bg-blue-500 transition"
      >
        <i class="pi pi-plus text-sm"></i>
        <span>Criar Primeira Trilha</span>
      </button>
    </div>

    <div v-else class="space-y-8">
      <div
        v-for="trilha in trilhas"
        :key="trilha.id"
        class="study-card p-6 sm:p-8 border-[#dce1d9] bg-[#ffffff] shadow-none space-y-6"
      >
        <div
          class="flex flex-col sm:flex-row sm:items-center justify-between pb-5 border-b border-[#dce1d9] gap-4"
        >
          <div class="flex items-center gap-3.5">
            <span
              class="w-12 h-12 rounded-lg flex items-center justify-center font-semibold text-base shadow-md"
              :class="
                trilha.tipo_da_linguagem === 1
                  ? 'bg-cyan-500/20 text-sky-800 border border-cyan-500/30'
                  : 'bg-amber-400/20 text-amber-800 border border-amber-400/30'
              "
            >
              {{ trilha.tipo_da_linguagem === 1 ? "Go" : "Py" }}
            </span>
            <div>
              <h2 class="text-xl font-semibold text-slate-800 tracking-tight">{{ trilha.name }}</h2>
              <span class="text-sm text-slate-600 font-medium">
                {{ trilha.tipo_da_linguagem === 1 ? "Golang" : "Python" }} •
                {{ trilha.listas?.length || 0 }} Listas
              </span>
            </div>
          </div>

          <div class="flex items-center gap-2">
            <button
              @click="
                isUpdateDialog = true;
                trilhaID = trilha.id;
                visibleFormManipulador = true;
              "
              class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-sm font-semibold text-amber-800 bg-amber-400/10 border border-amber-400/25 hover:bg-amber-400/20 transition"
            >
              <i class="pi pi-pencil text-[11px]"></i>
              <span>Editar Trilha</span>
            </button>
            <button
              @click="confirmarRemoverTrilha(trilha.id)"
              class="inline-flex items-center gap-1.5 px-3.5 py-2 rounded-lg text-sm font-semibold text-rose-700 bg-rose-500/10 border border-rose-500/25 hover:bg-rose-500/20 transition"
            >
              <i class="pi pi-trash text-[11px]"></i>
              <span>Excluir</span>
            </button>
          </div>
        </div>

        <div class="space-y-4">
          <Accordion :value="['0']" multiple>
            <AccordionPanel v-for="lista in trilha.listas" :key="lista.id" :value="lista.id">
              <AccordionHeader>
                <div class="flex items-center justify-between w-full pr-4">
                  <div class="flex items-center gap-2.5">
                    <i class="pi pi-folder text-sky-800 text-sm"></i>
                    <span class="font-semibold text-slate-800">{{ lista.name }}</span>
                  </div>
                  <span
                    class="text-sm px-2.5 py-0.5 rounded-full bg-[#ffffff] text-slate-600 border border-[#dce1d9] font-medium"
                  >
                    {{ lista.exercicios?.length || 0 }} exercícios
                  </span>
                </div>
              </AccordionHeader>

              <AccordionContent>
                <div class="p-4 space-y-5">
                  <div
                    class="p-4 rounded-lg bg-[#ffffff] border border-[#dce1d9] flex flex-col md:flex-row md:items-center justify-between gap-4"
                  >
                    <div class="flex items-center gap-3">
                      <div v-if="!lista.isEditing" class="text-base font-semibold text-slate-800">
                        {{ lista.name }}
                      </div>
                      <div v-else class="flex items-center gap-2">
                        <InputText
                          v-model="lista.name"
                          class="!py-1.5 !px-3 text-sm"
                          @blur="saveListaName(lista)"
                        />
                      </div>
                    </div>

                    <div class="flex flex-wrap items-center gap-2">
                      <button
                        @click="editListaName(lista, true)"
                        class="px-3 py-1.5 rounded-lg text-sm font-semibold text-slate-600 bg-[#ffffff] border border-[#dce1d9] hover:bg-[#ffffff] transition"
                      >
                        {{ lista.isEditing ? "Salvar Nome" : "Renomear" }}
                      </button>
                      <button
                        @click="alterarIndexLista(trilha.id, lista.id, -1)"
                        title="Subir lista"
                        class="p-1.5 w-8 h-8 rounded-lg text-sm text-slate-600 bg-[#ffffff] border border-[#dce1d9] hover:bg-[#ffffff] transition flex items-center justify-center"
                      >
                        <i class="pi pi-arrow-up"></i>
                      </button>
                      <button
                        @click="alterarIndexLista(trilha.id, lista.id, 1)"
                        title="Descer lista"
                        class="p-1.5 w-8 h-8 rounded-lg text-sm text-slate-600 bg-[#ffffff] border border-[#dce1d9] hover:bg-[#ffffff] transition flex items-center justify-center"
                      >
                        <i class="pi pi-arrow-down"></i>
                      </button>
                      <button
                        @click="confirmarRemoverLista(lista.id)"
                        class="px-3 py-1.5 rounded-lg text-sm font-semibold text-rose-700 bg-rose-500/10 border border-rose-500/30 hover:bg-rose-500/20 transition"
                      >
                        Remover Lista
                      </button>
                    </div>
                  </div>

                  <TreeTable :value="lista.exercicios" class="w-full">
                    <Column field="name" header="Exercício" expander>
                      <template #body="{ node }">
                        <span class="font-medium text-slate-700 text-sm">
                          {{ node.data.name }}
                        </span>
                      </template>
                    </Column>

                    <Column header="Ações" style="width: 14rem; text-align: right">
                      <template #body="{ node }">
                        <div class="flex items-center justify-end gap-1.5">
                          <NuxtLink
                            :to="`/dashboard_docente/exercicio/atualizar/${trilha.tipo_da_linguagem}/${node.data.id}`"
                          >
                            <button
                              class="px-2.5 py-1 rounded-md text-sm font-semibold text-sky-800 bg-cyan-500/10 border border-cyan-500/25 hover:bg-cyan-500/20 transition"
                            >
                              Editar
                            </button>
                          </NuxtLink>
                          <button
                            @click="alterarIndexExercicio(lista.id, node.data.id, -1)"
                            class="p-1 w-7 h-7 rounded-md text-sm text-slate-600 bg-[#ffffff] border border-[#dce1d9] hover:bg-[#ffffff] transition flex items-center justify-center"
                            title="Subir ordem"
                          >
                            <i class="pi pi-arrow-up text-[10px]"></i>
                          </button>
                          <button
                            @click="alterarIndexExercicio(lista.id, node.data.id, 1)"
                            class="p-1 w-7 h-7 rounded-md text-sm text-slate-600 bg-[#ffffff] border border-[#dce1d9] hover:bg-[#ffffff] transition flex items-center justify-center"
                            title="Descer ordem"
                          >
                            <i class="pi pi-arrow-down text-[10px]"></i>
                          </button>
                          <button
                            @click="removerExercicio(node.data.id)"
                            class="p-1 w-7 h-7 rounded-md text-sm text-rose-700 bg-rose-500/10 border border-rose-500/25 hover:bg-rose-500/20 transition flex items-center justify-center"
                            title="Excluir exercício"
                          >
                            <i class="pi pi-trash text-[10px]"></i>
                          </button>
                        </div>
                      </template>
                    </Column>
                  </TreeTable>

                  <div class="pt-4 border-t border-[#dce1d9] flex justify-end">
                    <NuxtLink
                      :to="`/dashboard_docente/exercicio/adicionar/${trilha.tipo_da_linguagem}/${lista.id}`"
                    >
                      <button
                        class="study-btn-primary px-3.5 py-1.5 text-sm flex items-center gap-1.5"
                      >
                        <i class="pi pi-plus text-[10px]"></i>
                        <span>Adicionar Exercício</span>
                      </button>
                    </NuxtLink>
                  </div>
                </div>
              </AccordionContent>
            </AccordionPanel>
          </Accordion>

          <div class="pt-2 flex justify-center">
            <button
              @click="addLista(trilha.id)"
              class="study-btn-secondary px-4 py-2 text-sm flex items-center gap-2"
            >
              <i class="pi pi-folder-plus text-sm text-sky-800"></i>
              <span>Adicionar Nova Lista nesta Trilha</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="flex justify-center pt-2">
      <button
        @click="formManipulador()"
        class="study-btn-primary px-6 py-3 text-sm flex items-center gap-2"
      >
        <i class="pi pi-plus text-sm"></i>
        <span>Adicionar Nova Trilha</span>
      </button>
    </div>
  </div>
</template>


<script setup lang="ts">
definePageMeta({
    middleware: 'auth-docente'
});

import { ref } from 'vue';
import type { Listas, Trilha } from '~/types';
import { useConfirm } from "primevue/useconfirm";
import { useToast } from "primevue/usetoast";

const visibleFormManipulador = ref(false);
const confirm = useConfirm();
const toast = useToast();

const isUpdateDialog = ref(false);
const trilhaID = ref('');

const trilhas = ref<Trilha[]>([]);
const { $authService } = useNuxtApp();

const formManipulador = () => {
  isUpdateDialog.value = false;
  trilhaID.value = '';
  visibleFormManipulador.value = !visibleFormManipulador.value;
};

const onRemoveExercicioById = async (exercicios_id: string) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/remover_exercicio/' + exercicios_id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
    });

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Exercício removido.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao remover o exercício.', life: 3000 });
    }
  } catch(err) {
    console.error(err);
  }
};

const removerExercicio = (exercicio_id: string) => {
  confirm.require({
    message: 'Você deseja realmente remover este exercício?',
    header: 'Confirmação de Exclusão',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancelar',
    rejectProps: {
      label: 'Cancelar',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: 'Excluir',
      severity: 'danger'
    },
    accept: () => {
      onRemoveExercicioById(exercicio_id);
    },
    reject: () => {
      // noop
    }
  });
};

const confirmarRemoverLista = (lista_id: string) => {
  confirm.require({
    message: 'Você deseja remover esta lista de exercícios?',
    header: 'Confirmação de Exclusão',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancelar',
    rejectProps: {
      label: 'Cancelar',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: 'Remover',
      severity: 'danger'
    },
    accept: () => {
      removeItemLista(lista_id);
    },
    reject: () => {
      // noop
    }
  });
};

const removeItemLista = async (lista_id: string) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/remove_lista/' + lista_id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
    });

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Lista removida.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao remover Lista.', life: 3000 });
    }
  } catch(err) {
    console.error(err);
  }
};

const addLista = async (idTrilha: string) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/add_lista', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        trilha_id: idTrilha,
      })
    });

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Lista criada com sucesso.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao criar Lista.', life: 3000 });
    }
  } catch (error) {
    console.error(error);
  }
};

const getTrilhasListaExercicios = async () => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/trilhas_lista_exercicios', {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    });

    if (response.ok) {
      const resData = await response.json();
      trilhas.value = resData;
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao carregar trilhas.', life: 3000 });
    }
  } catch(error) {
    console.error(error);
  }
};

const confirmarRemoverTrilha = (trilha_id: string) => {
  confirm.require({
    message: 'Você deseja realmente remover esta trilha completa e todas suas listas?',
    header: 'Zona de Perigo',
    icon: 'pi pi-exclamation-triangle',
    rejectLabel: 'Cancelar',
    rejectProps: {
      label: 'Cancelar',
      severity: 'secondary',
      outlined: true
    },
    acceptProps: {
      label: 'Excluir Trilha',
      severity: 'danger'
    },
    accept: () => {
      removeItemTrilha(trilha_id);
    },
    reject: () => {
      // noop
    }
  });
};

const removeItemTrilha = async (trilha_id: string) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/remover_trilha/' + trilha_id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
    });

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Trilha removida com sucesso.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao remover Trilha.', life: 3000 });
    }
  } catch(err) {
    console.error(err);
  }
};

const editListaName = (lista: Listas, from_button: boolean) => {
  if (from_button && lista.isEditing) {
    saveListaName(lista);
    return;
  }
  lista.isEditing = !lista.isEditing;
};

const saveListaName = async (lista: Listas) => {
  lista.isEditing = false;
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/update_lista', {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        name: lista.name,
        id: lista.id,
      })
    });

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Lista atualizada com sucesso.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao atualizar a lista.', life: 3000 });
    }
  } catch(error) {
    console.error(error);
  }
};

const alterarIndexExercicio = async (lista_id: string, exercicio_id: string, posicoes_a_trocar: number) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/update_exercicio_index', {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        lista_id: lista_id,
        exercicio_id: exercicio_id,
        posicoes_a_trocar: posicoes_a_trocar,
      })
    });

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Ordem do exercício alterada.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao alterar ordem do exercício.', life: 3000 });
    }
  } catch (error) {
    console.error(error);
  }
};

const alterarIndexLista = async (trilha_id: string, lista_id: string, posicoes_a_trocar: number) => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/update_lista_index', {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        trilha_id: trilha_id,
        lista_id: lista_id,
        posicoes_a_trocar: posicoes_a_trocar,
      })
    });

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Ordem da lista alterada.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao alterar ordem da lista.', life: 3000 });
    }
  } catch (error) {
    console.error(error);
  }
};

await getTrilhasListaExercicios();
</script>