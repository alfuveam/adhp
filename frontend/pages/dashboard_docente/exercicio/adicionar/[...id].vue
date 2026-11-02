<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-6">
    <Toast />
    <ConfirmDialog />

    <!-- Top Navigation -->
    <div class="flex items-center justify-between pb-2 border-b border-white/[0.08]">
      <NuxtLink 
        to="/dashboard_docente" 
        class="inline-flex items-center gap-2 text-sm font-medium text-slate-400 hover:text-white transition"
      >
        <i class="pi pi-arrow-left text-xs"></i>
        <span>Voltar ao Painel Docente</span>
      </NuxtLink>

      <span class="mimo-badge mimo-badge-indigo text-[11px]">
        {{ tipoDaLinguagem === 1 ? 'Linguagem Golang' : 'Linguagem Python' }}
      </span>
    </div>

    <!-- Header Card -->
    <div class="mimo-card p-6 sm:p-7 border-[#363e66] bg-[#1c2038] shadow-xl space-y-1.5">
      <h1 class="text-2xl font-extrabold text-white tracking-tight">
        Cadastrar Novo Exercício
      </h1>
      <p class="text-xs sm:text-sm text-slate-300 leading-relaxed">
        Defina o enunciado, o código base com boilerplate e a suíte de testes unitários para avaliação automática.
      </p>
    </div>

    <!-- Form -->
    <Form 
      v-slot="$form" 
      :initialValues="valoresIniciais" 
      :resolver="resolver" 
      @submit="onFormSubmit" 
      class="flex flex-col gap-6"
    >
      <!-- Título do Exercício -->
      <div class="mimo-card p-6 border-[#363e66] bg-[#1c2038] shadow-xl flex flex-col gap-2">
        <label class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
          <i class="pi pi-bookmark text-purple-400 text-xs"></i>
          <span>Título do Exercício</span>
        </label>
        <InputText 
          name="titulo_exercicio" 
          v-model="valoresIniciais.titulo_exercicio" 
          type="text" 
          placeholder="Ex: Soma de Elementos em uma Lista" 
          class="w-full !p-3 !bg-[#15172a] !border-[#363e66] focus:!border-[#7c4ed4] rounded-xl text-sm"
          fluid 
        />
        <Message v-if="$form.titulo_exercicio?.invalid" severity="error" size="small" variant="simple" class="!text-xs !text-rose-400 mt-1">
          {{ $form.titulo_exercicio.error?.message }}
        </Message>
      </div>

      <!-- Código Base -->
      <div class="rounded-2xl border border-white/10 overflow-hidden shadow-2xl bg-[#0d121f]">
        <div class="flex items-center justify-between px-4 py-2.5 bg-[#0a0e17] border-b border-white/[0.08]">
          <div class="flex items-center gap-1.5">
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
          </div>
          <span class="text-xs font-mono font-bold text-slate-300">Código Base do Aluno</span>
          <span class="text-[11px] text-slate-500 font-mono">{{ tipoDaLinguagem === 1 ? 'Go' : 'Python' }}</span>
        </div>
        <MonacoEditor 
          name="codigo_base" 
          v-model="valoresIniciais.codigo_base" 
          :options="optionsEditor" 
          :lang="tipoDaLinguagem == 1 ? 'go' : 'python'" 
          class="h-[460px] z-10" 
        />
        <Message v-if="$form.codigo_base?.invalid" severity="error" size="small" variant="simple" class="!text-xs !text-rose-400 p-3">
          {{ $form.codigo_base.error?.message }}
        </Message>
      </div>

      <!-- Retorno da Execução do Código Base -->
      <div v-if="(returnExec.out_put_from_user.success && returnExec.out_put_from_user.output !== '') || (!returnExec.out_put_from_user.success && returnExec.out_put_from_user.error !== '')" class="mimo-card p-4 border border-white/10">
        <label class="text-xs font-bold uppercase tracking-wider text-slate-300 mb-2 block">
          {{ returnExec.out_put_from_user.success ? 'Saída da execução do código base:' : 'Erro da execução do código base:' }}
        </label>
        <Textarea 
          :model-value="returnExec.out_put_from_user.success ? returnExec.out_put_from_user.output : returnExec.out_put_from_user.error" 
          rows="4" 
          readonly 
          class="w-full !bg-slate-900/90 !font-mono !text-xs !border-white/10" 
        />
      </div>

      <!-- Código de Teste -->
      <div class="rounded-2xl border border-white/10 overflow-hidden shadow-2xl bg-[#0d121f]">
        <div class="flex items-center justify-between px-4 py-2.5 bg-[#0a0e17] border-b border-white/[0.08]">
          <div class="flex items-center gap-1.5">
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
            <span class="w-2.5 h-2.5 rounded-full bg-slate-600"></span>
          </div>
          <span class="text-xs font-mono font-bold text-slate-300">Código da Suíte de Testes</span>
          <span class="text-[11px] text-slate-500 font-mono">{{ tipoDaLinguagem === 1 ? 'Go Test' : 'PyTest' }}</span>
        </div>
        <MonacoEditor 
          name="codigo_teste" 
          v-model="valoresIniciais.codigo_teste" 
          :options="optionsEditor" 
          :lang="tipoDaLinguagem == 1 ? 'go' : 'python'" 
          class="h-[460px]" 
        />
        <Message v-if="$form.codigo_teste?.invalid" severity="error" size="small" variant="simple" class="!text-xs !text-rose-400 p-3">
          {{ $form.codigo_teste.error?.message }}
        </Message>
      </div>

      <!-- Retorno da Execução do Código Teste -->
      <div v-if="(returnExec.out_put_unit_teste.success && returnExec.out_put_unit_teste.output !== '') || (!returnExec.out_put_unit_teste.success && returnExec.out_put_unit_teste.error !== '')" class="mimo-card p-4 border border-white/10">
        <label class="text-xs font-bold uppercase tracking-wider text-slate-300 mb-2 block">
          {{ returnExec.out_put_unit_teste.success ? 'Saída dos testes unitários:' : 'Erro dos testes unitários:' }}
        </label>
        <Textarea 
          :model-value="returnExec.out_put_unit_teste.success ? returnExec.out_put_unit_teste.output : returnExec.out_put_unit_teste.error" 
          rows="4" 
          readonly 
          class="w-full !bg-slate-900/90 !font-mono !text-xs !border-white/10" 
        />
      </div>

      <!-- Feedbacks Automatizados -->
      <div class="mimo-card p-6 border border-white/10 shadow-xl space-y-4">
        <div class="flex items-center justify-between pb-3 border-b border-white/[0.08]">
          <label class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-2">
            <i class="pi pi-comments text-indigo-400"></i>
            <span>Dicas e Feedbacks Pré-Cadastrados</span>
          </label>
          <button 
            type="button" 
            @click="adicionarFeedback()" 
            class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold text-white bg-indigo-600 hover:bg-indigo-500 transition"
          >
            <i class="pi pi-plus text-[10px]"></i>
            <span>Nova Dica</span>
          </button>
        </div>

        <DataTable :value="feedbacks" class="w-full">
          <Column field="descricao" header="Descrição do Feedback">
            <template #body="{ data }">
              <div @click="editFeedback(data)" v-if="!data.isEditing" class="cursor-pointer text-sm text-slate-200 hover:text-indigo-300">
                {{ data.descricao || '(Clique para editar a dica)' }}
              </div>
              <InputText v-else v-model="data.descricao" @blur="saveFeedback(data)" class="!p-2 text-sm w-full" />
            </template>
          </Column>
          <Column style="width: 5rem; text-align: right;">
            <template #body="{ data }">
              <button 
                type="button" 
                @click="removerFeedback(data.id)" 
                class="p-1.5 rounded-lg text-rose-400 hover:bg-rose-500/10 transition"
                title="Remover dica"
              >
                <i class="pi pi-trash text-xs"></i>
              </button>
            </template>
          </Column>
        </DataTable>
      </div>

      <!-- Submit Action -->
      <div class="flex justify-end pt-2">
        <button 
          type="submit" 
          class="mimo-btn-primary px-8 py-3 text-sm font-bold flex items-center gap-2"
        >
          <i class="pi pi-check text-xs"></i>
          <span>Cadastrar Exercício</span>
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
import type { Feedback, RetornoExecucao} from '~/types';
import { useRuntimeConfig } from 'nuxt/app';

const route = useRoute();
const confirm = useConfirm();
const toast = useToast();

const tipoDaLinguagem = ref(0)
const listaID = ref('')

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

const refQuillCodigoBase = ref(null);
const refQuillCodigoTeste = ref(null);

const { $authService } = useNuxtApp()

let codigoBaseToSend = ref('')
let codigoTesteToSend = ref('')

const valoresIniciais = ref({
  titulo_exercicio: '',
  codigo_base: '',
  codigo_teste: '',
})

const resolver = zodResolver(
  z.object({
    titulo_exercicio: z.string().min(1, { message: 'Título exercício é necessário.' }),
    codigo_base: z.string().min(1, { message: 'Código base é necessário.' }),
    codigo_teste: z.string().min(1, { message: 'Código teste é necessário.' }),
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

    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + "/v1/adicionar_exericicio", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      },
      body: JSON.stringify({
        titulo: valoresIniciais.value.titulo_exercicio,
        lista_id: listaID.value,
        codigo_base: btoa(encodeURIComponent(codigoBaseToSend.value)),
        codigo_teste: btoa(encodeURIComponent(codigoTesteToSend.value)),
        feedbacks: feedbacks.value,
      }),
    })

    returnExec.value = await response.json()
    if (response.ok) {
      toast.add({ severity: 'info', summary: 'Confirmed', detail: 'Adicionado com sucesso.', life: 3000 });
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

const removeFeedbackById = (id: string)  => {
  feedbacks.value = feedbacks.value.filter((feedback_item) => feedback_item.id !== id);
}

const adicionarFeedback = () => {
  feedbacks.value.push({id: Math.floor(Math.random() * 500).toString(), descricao: "alterar"})
}

const removerFeedback = (id: string) => {
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
        toast.add({ severity: 'info', summary: 'Confirmado', detail: 'Feedback deletado', life: 3000 });
        removeFeedbackById(id)
    },

    reject: () => {
        toast.add({ severity: 'error', summary: 'Rejeitado', detail: 'Você rejeitou.', life: 3000 });
    }
  });
};

onMounted(() => {
  tipoDaLinguagem.value = Number(route.params.id[0])
  listaID.value = route.params.id[1]
});
</script>
