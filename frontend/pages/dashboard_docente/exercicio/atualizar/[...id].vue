<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />
    <ConfirmDialog />

    <div class="flex items-center justify-between pb-2 border-b border-slate-200">
      <NuxtLink
        to="/dashboard_docente"
        class="inline-flex items-center gap-2 text-sm font-medium text-slate-600 hover:text-slate-800 transition"
      >
        <i class="pi pi-arrow-left text-sm"></i>
        <span>Voltar ao Painel Docente</span>
      </NuxtLink>

      <span class="study-badge study-badge-indigo text-[11px]">
        {{ tipoDaLinguagem === 1 ? "Linguagem Golang" : "Linguagem Python" }}
      </span>
    </div>

    <div class="study-card p-6 sm:p-7 border-[#dce1d9] bg-[#ffffff] shadow-none space-y-1.5">
      <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Atualizar Exercício</h1>
      <p class="text-sm sm:text-sm text-slate-600 leading-relaxed">
        Edite o título, código base fornecido ao estudante e a suíte de testes unitários.
      </p>
    </div>

    <Form
      v-slot="$form"
      :initialValues="valoresIniciais"
      :resolver="resolver"
      @submit="onFormSubmit"
      class="flex flex-col gap-6"
    >
      <div class="study-card p-6 border-[#dce1d9] bg-[#ffffff] shadow-none flex flex-col gap-2">
        <label class="text-sm font-semibold text-slate-600 flex items-center gap-1.5">
          <i class="pi pi-bookmark text-sky-800 text-sm"></i>
          <span>Título do Exercício</span>
        </label>
        <InputText
          name="titulo_exercicio"
          v-model="valoresIniciais.titulo_exercicio"
          type="text"
          placeholder="Título do exercício"
          class="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-sm"
          fluid
        />
        <Message
          v-if="$form.titulo_exercicio?.invalid"
          severity="error"
          size="small"
          variant="simple"
          class="!text-sm !text-rose-700 mt-1"
        >
          {{ $form.titulo_exercicio.error?.message }}
        </Message>
      </div>

      <div class="rounded-lg border border-slate-200 overflow-hidden shadow-none bg-[#ffffff]">
        <div
          class="flex items-center justify-between px-4 py-2.5 bg-[#ffffff] border-b border-slate-200"
        >
          <div class="flex items-center gap-1.5">
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
          </div>
          <span class="text-sm font-mono font-semibold text-slate-600">Código Base do Aluno</span>
          <span class="text-[11px] text-slate-500 font-mono">{{
            tipoDaLinguagem === 1 ? "Go" : "Python"
          }}</span>
        </div>
        <MonacoEditor
          name="codigo_base"
          v-model="valoresIniciais.codigo_base"
          :options="optionsEditor"
          :lang="tipoDaLinguagem == 1 ? 'go' : 'python'"
          class="h-[460px] z-10"
        />
        <Message
          v-if="$form.codigo_base?.invalid"
          severity="error"
          size="small"
          variant="simple"
          class="!text-sm !text-rose-700 p-3"
        >
          {{ $form.codigo_base.error?.message }}
        </Message>
      </div>

      <div
        v-if="
          (returnExec.out_put_from_user.success && returnExec.out_put_from_user.output !== '') ||
          (!returnExec.out_put_from_user.success && returnExec.out_put_from_user.error !== '')
        "
        class="study-card p-4 border border-slate-200"
      >
        <label class="text-sm font-semibold text-slate-600 mb-2 block">
          {{
            returnExec.out_put_from_user.success
              ? "Saída da execução do código base:"
              : "Erro da execução do código base:"
          }}
        </label>
        <Textarea
          :model-value="
            returnExec.out_put_from_user.success
              ? returnExec.out_put_from_user.output
              : returnExec.out_put_from_user.error
          "
          rows="4"
          readonly
          class="w-full !bg-slate-50 !font-mono !text-sm !border-slate-200"
        />
      </div>

      <div class="rounded-lg border border-slate-200 overflow-hidden shadow-none bg-[#ffffff]">
        <div
          class="flex items-center justify-between px-4 py-2.5 bg-[#ffffff] border-b border-slate-200"
        >
          <div class="flex items-center gap-1.5">
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
          </div>
          <span class="text-sm font-mono font-semibold text-slate-600"
            >Código da Suíte de Testes</span
          >
          <span class="text-[11px] text-slate-500 font-mono">{{
            tipoDaLinguagem === 1 ? "Go Test" : "PyTest"
          }}</span>
        </div>
        <MonacoEditor
          name="codigo_teste"
          v-model="valoresIniciais.codigo_teste"
          :options="optionsEditor"
          :lang="tipoDaLinguagem == 1 ? 'go' : 'python'"
          class="h-[460px]"
        />
        <Message
          v-if="$form.codigo_teste?.invalid"
          severity="error"
          size="small"
          variant="simple"
          class="!text-sm !text-rose-700 p-3"
        >
          {{ $form.codigo_teste.error?.message }}
        </Message>
      </div>

      <div
        v-if="
          (returnExec.out_put_unit_teste.success && returnExec.out_put_unit_teste.output !== '') ||
          (!returnExec.out_put_unit_teste.success && returnExec.out_put_unit_teste.error !== '')
        "
        class="study-card p-4 border border-slate-200"
      >
        <label class="text-sm font-semibold text-slate-600 mb-2 block">
          {{
            returnExec.out_put_unit_teste.success
              ? "Saída dos testes unitários:"
              : "Erro dos testes unitários:"
          }}
        </label>
        <Textarea
          :model-value="
            returnExec.out_put_unit_teste.success
              ? returnExec.out_put_unit_teste.output
              : returnExec.out_put_unit_teste.error
          "
          rows="4"
          readonly
          class="w-full !bg-slate-50 !font-mono !text-sm !border-slate-200"
        />
      </div>

      <div class="study-card p-6 border border-slate-200 shadow-none space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-slate-200">
          <label class="text-sm font-semibold text-slate-600 flex items-center gap-2">
            <i class="pi pi-comments text-indigo-400"></i>
            <span>Dicas e Feedbacks Pré-Cadastrados</span>
          </label>
          <button
            type="button"
            @click="adicionarFeedback()"
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-semibold text-slate-800 bg-indigo-600 hover:bg-indigo-500 transition"
          >
            <i class="pi pi-plus text-[10px]"></i>
            <span>Nova Dica</span>
          </button>
        </div>

        <DataTable :value="feedbacks" class="w-full">
          <Column field="descricao" header="Descrição do Feedback">
            <template #body="{ data }">
              <div
                @click="editFeedback(data)"
                v-if="!data.isEditing"
                class="cursor-pointer text-sm text-slate-700 hover:text-indigo-300"
              >
                {{ data.descricao || "(Clique para editar a dica)" }}
              </div>
              <InputText
                v-else
                v-model="data.descricao"
                @blur="saveFeedback(data)"
                class="!p-2 text-sm w-full"
              />
            </template>
          </Column>
          <Column style="width: 5rem; text-align: right">
            <template #body="{ data }">
              <button
                type="button"
                @click="removerFeedback(data)"
                class="p-1.5 rounded-lg text-rose-700 hover:bg-rose-500/10 transition"
                title="Remover dica"
              >
                <i class="pi pi-trash text-sm"></i>
              </button>
            </template>
          </Column>
        </DataTable>
      </div>

      <div class="flex justify-end pt-2">
        <button
          type="submit"
          class="study-btn-primary px-8 py-3 text-sm font-semibold flex items-center gap-2"
        >
          <i class="pi pi-save text-sm"></i>
          <span>Atualizar Exercício</span>
        </button>
      </div>
    </Form>
  </div>
</template>


<script setup lang="ts">
definePageMeta({
    middleware: 'auth-docente'
})

import * as Monaco from 'monaco-editor'

const optionsEditor = ref<Monaco.editor.IStandaloneEditorConstructionOptions>({
  theme: 'vs-dark'
});

import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { useToast } from 'primevue/usetoast';
import { useConfirm } from "primevue/useconfirm";
import type { Feedback, RetornoExecucao, ExercicioManipulador} from '~/types';

const route = useRoute();
const confirm = useConfirm();
const toast = useToast();

const tipoDaLinguagem = ref(0)
const exercicioID = ref('')

const feedbacks = ref<Feedback[]>([]);
const returnExec = ref<RetornoExecucao>(
  {
    out_put_from_user: {
      success: false,
      error: '',
      output: '',
    },
    out_put_unit_teste: {
      success: false,
      error: '',
      output: '',
    },
  }
);

const refQuillCodigoBase = ref(null)
const refQuillCodigoTeste = ref(null)

const { $authService } = useNuxtApp()

let codigoBaseToSend = ref('')
let codigoTesteToSend = ref('')

const valoresIniciais = ref({
    titulo_exercicio: '',
    codigo_base: '',
    codigo_teste: '',
});

const resolver = zodResolver(
    z.object({
        titulo_exercicio: z.string().min(5, { message: 'Título exercício é necessário.' }).nullish(),
        codigo_base: z.string().min(10, { message: 'Código base é necessário.' }).nullish(),
        codigo_teste: z.string().min(10, { message: 'Código teste é necessário.' }).nullish(),
    })
);

const onFormSubmit = async ({ valid }: { valid: boolean }) => {
  if (valid) {
      toast.add({
        severity: 'success',
        summary: 'Form enviado.',
        life: 3000
    });

    if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
      codigoBaseToSend.value = valoresIniciais.value.codigo_base
      codigoTesteToSend.value = valoresIniciais.value.codigo_teste
    } else {
      // @ts-ignore
      codigoBaseToSend.value = refQuillCodigoBase.value?.quill.getText() || ''
      // @ts-ignore
      codigoTesteToSend.value = refQuillCodigoTeste.value?.quill.getText() || ''
    }

    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + "/v1/atualizar_exericicio", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        id: exercicioID.value,
        titulo: valoresIniciais.value.titulo_exercicio,
        codigo_base: btoa(encodeURIComponent(codigoBaseToSend.value)),
        codigo_teste: btoa(encodeURIComponent(codigoTesteToSend.value)),
        feedbacks: feedbacks.value,
      }),
    })

    returnExec.value = await response.json()
    if (response.ok) {
      toast.add({ severity: 'info', summary: 'Confirmed', detail: 'Atualizado com sucesso.', life: 3000 });
    } else {
      toast.add({ severity: 'error', summary: 'Rejected', detail: "Erro ao submeter a questão.", life: 12000 });
    }
  }
};

const editFeedback = (data: any) => {
  data.isEditing = true;
};

const saveFeedback = (data: any) => {
  data.isEditing = false;
};

const removeFeedbackById = (feedback: Feedback)  => {
  if (feedback.id === "0") {
    //  if has a same reference
    feedbacks.value = feedbacks.value.filter((feedback_item) => feedback_item !== feedback);
  } else {
    feedbacks.value = feedbacks.value.filter((feedback_item) => feedback_item.id !== feedback.id);
  }
}

const adicionarFeedback = () => {
  if (feedbacks.value === null) {
    feedbacks.value = []
  }
  feedbacks.value.push({id: "0", descricao: "alterar"})
}

const onRemoveFeedback = async (feedback: Feedback) => {
  if (feedback.id === "0") return
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/remove_feedback/' + feedback.id, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
    })

    if (response.ok) {
      toast.add({ severity: 'success', summary: 'Feedback removida.', life: 3000 });

    } else {
      toast.add({ severity: 'error', summary: 'Falha ao remover Feedback.', life: 3000 });
    }
  } catch (error) {
      console.log(error)
  }
}

const removerFeedback = (feedback: Feedback) => {
  confirm.require({
    message: 'Você tem certeza que deseja excluir o feedback?',
    header: 'Zona de perigo',
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
        removeFeedbackById(feedback)
        onRemoveFeedback(feedback)
    },

    reject: () => {
        toast.add({ severity: 'error', summary: 'Rejeitado', detail: 'Você rejeitou.', life: 3000 });
    }
  });
};

const getExercicioById = async () => {
  const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + "/v1/exercicio/" + exercicioID.value, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      "Authorization": "Bearer " + $authService.getToken()
    },
  })

  if (response.ok) {
    const resData: ExercicioManipulador = await response.json()
    if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
      valoresIniciais.value.codigo_base = decodeURIComponent(atob(resData.codigo_base))
      valoresIniciais.value.codigo_teste = decodeURIComponent(atob(resData.codigo_teste))
    } else {
      // @ts-ignore
      refQuillCodigoBase.value?.quill.setText(decodeURIComponent(atob(resData.codigo_base)))
      // @ts-ignore
      refQuillCodigoTeste.value?.quill.setText(decodeURIComponent(atob(resData.codigo_teste)))
    }

    valoresIniciais.value.titulo_exercicio = resData.titulo
    feedbacks.value = resData.feedbacks
  } else {
    toast.add({ severity: 'error', summary: 'Falha ao carregar o dashboard.', life: 3000 });
  }
}

onMounted(() => {
  tipoDaLinguagem.value = Number(route.params.id[0])
  exercicioID.value = route.params.id[1]

  const intervalId = setInterval(checkIfQuillIsLoad,100)

  function checkIfQuillIsLoad() {
    if (useRuntimeConfig().public.USE_MONACO_EDITOR == "true") {
      clearInterval(intervalId);
      getExercicioById()
    } else {
      // @ts-ignore
      if (refQuillCodigoBase && refQuillCodigoTeste && refQuillCodigoBase.value?.quill.root && refQuillCodigoTeste.value?.quill.root) {
        clearInterval(intervalId);
        getExercicioById()
      }
    }
  }
});
</script>