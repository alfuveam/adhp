<template>
  <Toast />
  <ConfirmDialog></ConfirmDialog>
  <ManipuladorTrilha 
    v-model="visibleFormManipulador" 
    :isUpdate="isUpdateDialog"
    :trilhaID="trilhaID"
    @some-event="getTrilhasListaExercicios" @blur="formManipulador()" />

  <template v-for="trilha in trilhas">
    <div class="card">
      <div class="flex flex-col justify-center items-center">
        <div>
          <h1 class="mb-4 text-2xl font-extrabold text-center">{{ trilha.name }}</h1>
        </div>
        <div class="flex flex-wrap gap-2">
          <Button @click="isUpdateDialog = true; trilhaID = trilha.id; visibleFormManipulador = true" label="Editar" severity="war"></Button>
          <Button @click="confirmarRemoverTrilha(trilha.id)" label="Excluir" severity="danger" outlined></Button>
        </div>
      </div>

      <Accordion :value="['0']" multiple v-for="lista in trilha.listas">
        <AccordionPanel :value="lista.id">
          <AccordionHeader>
            {{lista.name}}
          </AccordionHeader>
          <AccordionContent>
              <div class="card">
                <TreeTable :value="lista.exercicios" tableStyle="min-width: 50px">
                  <template #header>
                    <div class="flex flex-row justify-center">
                      <div class="flex flex-col justify-center items-center w-auto">
                        <div>
                          <div class="mb-4 text-xl font-bold" @click="editListaName(lista, false)" v-if="!lista.isEditing">{{ lista.name }}</div>
                          <InputText class="mb-4" v-else v-model="lista.name" @blur="saveListaName(lista)" />
                        </div>
                        <div class="flex flex-wrap gap-2">
                          <Button @click="editListaName(lista, true)" :label="lista.isEditing ? 'Salvar' : 'Editar'" severity="war"></Button>
                          <Button @click="confirmarRemoverLista(lista.id)" label="Remover" severity="danger" outlined/>
                          <Button @click="alterarIndexLista(trilha.id, lista.id, -1)" label="^" severity="secondary" outlined></Button>
                          <Button @click="alterarIndexLista(trilha.id, lista.id, 1)" label="v" severity="secondary" outlined></Button>
                        </div>
                      </div>
                    </div>
                  </template>
                  <Column field="name" header="Name" expander style="width: 250px">
                    <div class="mt-2">
                      <InputText placeholder="Enter new value" />
                    </div>
                  </Column>
                  <Column style="width: 10rem">
                    <template #body="{ node }">
                      <div class="flex flex-wrap gap-2">
                        <NuxtLink :to="`/dashboard_docente/exercicio/atualizar/${trilha.tipo_da_linguagem}/${node.data.id}`">
                          <Button label="Atualizar"/>
                        </NuxtLink>
                        <Button @click="removerExercicio(node.data.id)" label="Remover" severity="danger" outlined></Button>
                        <Button @click="alterarIndexExercicio(lista.id, node.data.id, -1)" label="^" severity="secondary" outlined></Button>
                        <Button @click="alterarIndexExercicio(lista.id, node.data.id, 1)" label="v" severity="secondary" outlined></Button>
                      </div>
                    </template>
                  </Column>
                  <template #footer>
                    <NuxtLink class="flex flex-row justify-center" :to="`/dashboard_docente/exercicio/adicionar/${trilha.tipo_da_linguagem}/${lista.id}`">
                      <Button label="Adicionar Exercício" severity="warn" />
                    </NuxtLink>
                  </template>
                </TreeTable>
              </div>
            </AccordionContent>
          </AccordionPanel>
        </Accordion>
      <div class="flex flex-wrap gap-2 items-center justify-center pt-2">
        <Button label="Adicionar Lista" severity="warn" class="items-center" @click="addLista(trilha.id)" />
      </div>
  </div>
  <br>
</template>
<div class="flex w-full justify-center">
  <div class="flex flex-col w-54 justify-center">
    <Button icon="pi pi-refresh" label="Adicionar nova trilha" severity="warn" class="mt-2" @click="formManipulador()" />
  </div>
</div>
</template>

<script setup lang="ts">
definePageMeta({
    middleware: 'auth-docente'
})

import { ref } from 'vue';
import type { Listas, Trilha } from '~/types';
import { useConfirm } from "primevue/useconfirm";
import { useToast } from "primevue/usetoast";

const visibleFormManipulador = ref(false)
const confirm = useConfirm();
const toast = useToast();

const isUpdateDialog = ref(false)
const trilhaID = ref('')

const trilhas = ref<Trilha[]>([]);

const { $authService } = useNuxtApp();

const formManipulador = () => {
  isUpdateDialog.value = false
  trilhaID.value = ''
  visibleFormManipulador.value = !visibleFormManipulador.value
}

const onRemoveExercicioById = async (exercicios_id: string) => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/remover_exercicio/' + exercicios_id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Exercício removido.', life: 3000 });
      getTrilhasListaExercicios()
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao remover o exercício.', life: 3000 });
    }
  } catch(err) {
    console.error(err)
  }  
}

const removerExercicio = (lista_id: string) => {
    confirm.require({
        message: 'Você deseja remover esse exercício?',
        header: 'Zona de Perigo',
        icon: 'pi pi-info-circle',
        rejectLabel: 'Cancelar',
        rejectProps: {
            label: 'Cancelar',
            severity: 'secondary',
            outlined: true
        },
        acceptProps: {
            label: 'Deletar',
            severity: 'danger'
        },
        accept: () => {
            onRemoveExercicioById(lista_id)
        },
        reject: () => {
            toast.add({ severity: 'error', summary: 'Rejeitado', detail: 'Você rejeitou', life: 3000 });
        }
    });
};

const confirmarRemoverLista = (lista_id: string) => {
  confirm.require({
      message: 'Você deseja remover essa lista?',
      header: 'Zona Perigosa',
      icon: 'pi pi-info-circle',
      rejectLabel: 'Cancelar',
      rejectProps: {
          label: 'Cancellar',
          severity: 'secondary',
          outlined: true
      },
      acceptProps: {
          label: 'Remover',
          severity: 'danger'
      },
      accept: () => {
          removeItemLista(lista_id)
      },
      reject: () => {
          toast.add({ severity: 'error', summary: 'Rejeitado', detail: 'Você rejeitou', life: 3000 });
      }
  });
};

const removeItemLista = async (lista_id: string) => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/remove_lista/' + lista_id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Lista removida.', life: 3000 });
      getTrilhasListaExercicios()
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao remover Lista.', life: 3000 });
    }
  } catch(err) {
    console.error(err)
  }
}

const addLista = async (idTrilha: string) => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/add_lista', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        trilha_id: idTrilha,
      })
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Lista criada.', life: 3000 });
      getTrilhasListaExercicios()
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao criar Lista.', life: 3000 });
    }
  } catch (error) {
    console.error(error)
  }
}

const getTrilhasListaExercicios = async () => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/trilhas_lista_exercicios', {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    })

    if (response.ok) {
      const resData = await response.json()
      trilhas.value = resData
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao carregar o dashboard.', life: 3000 });
    }
  } catch(error) {
    console.error(error)
  }
}

const confirmarRemoverTrilha = (trilha_id: string) => {
  confirm.require({
      message: 'Você deseja remover essa trilha?',
      header: 'Zona Perigosa',
      icon: 'pi pi-info-circle',
      rejectLabel: 'Cancelar',
      rejectProps: {
          label: 'Cancellar',
          severity: 'secondary',
          outlined: true
      },
      acceptProps: {
          label: 'Remover',
          severity: 'danger'
      },
      accept: () => {
          removeItemTrilha(trilha_id)
      },
      reject: () => {
          toast.add({ severity: 'error', summary: 'Rejeitado', detail: 'Você rejeitou', life: 3000 });
      }
  });
};

const removeItemTrilha = async (trilha_id: string) => {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/remover_trilha/' + trilha_id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Trilha removida.', life: 3000 });
      getTrilhasListaExercicios()
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao remover Trilha.', life: 3000 });
    }
  } catch(err) {
    console.error(err)
  }
}

const editListaName = (lista: Listas, from_button: boolean) => {
  if (from_button &&lista.isEditing) {
    saveListaName(lista)
    return
  }
  lista.isEditing = !lista.isEditing;
}

const saveListaName = async (lista: Listas) => {
  lista.isEditing = false;
  try {
    const config = useRuntimeConfig()
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
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Lista atualizada.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao atualizar a lista.', life: 3000 });
    }
  } catch(error) {
    console.error(error)
  }
}

const alterarIndexExercicio = async (lista_id: string, exercicio_id: string, posicoes_a_trocar: number) => {
  try {
    const config = useRuntimeConfig()
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
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Ordem do exercício atualizada.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao atualizar a ordem do exercício.', life: 3000 });
    }
  } catch (error) {
    console.error(error)
  }
}

const alterarIndexLista = async (trilha_id: string, lista_id: string, posicoes_a_trocar: number) => {
  try {
    const config = useRuntimeConfig()
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
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Ordem do lista atualizada.', life: 3000 });
      getTrilhasListaExercicios();
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao atualizar a ordem do lista.', life: 3000 });
    }
  } catch (error) {
    console.error(error)
  }
}
await getTrilhasListaExercicios()

</script>