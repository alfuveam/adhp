<template>
  <div class="card flex justify-center">
    <Toast />
    <ConfirmDialog></ConfirmDialog>

    <Form v-slot="$form" :valoresIniciais :resolver @submit="onFormSubmit" class="flex flex-col gap-4 w-full sm:w-2/3">
      <div class="flex flex-col gap-1">
          <label for="multiple-ac-1" class="font-bold mb-2 block">Título exercício</label>
          <InputText name="titulo_exercicio" v-model="valoresIniciais.titulo_exercicio" type="text" placeholder="Título exercício" fluid />
          <Message v-if="$form.titulo_exercicio?.invalid" severity="error" size="small" variant="simple">{{ $form.titulo_exercicio.error?.message }}</Message>
      </div>
      <div class="flex flex-col gap-1">
        <label for="multiple-ac-1" class="font-bold mb-2 block">Código Base</label>
        <MonacoEditor name="codigo_base" v-model="valoresIniciais.codigo_base" :options="optionsEditor" :lang="tipoDaLinguagem == 1 ? 'go' : 'python'" class="h-[600px] z-10" />
        <!-- <Editor name="codigo_base" v-model="valoresIniciais.codigo_base" editorStyle="height: 220px" spellcheck="false" ref="refQuillCodigoBase">
          <template v-slot:toolbar>
            <span class="ql-formats">
              <button v-tooltip.bottom="'Bold'" class="ql-bold"></button>
              <button v-tooltip.bottom="'Italic'" class="ql-italic"></button>
              <button v-tooltip.bottom="'Underline'" class="ql-underline"></button>
            </span>
          </template>
        </Editor> -->
        <Message v-if="$form.codigo_base?.invalid" severity="error" size="small" variant="simple">{{ $form.codigo_base.error?.message }}</Message>
      </div>
      <div class="flex flex-col gap-1 p-4">
        <div v-if="returnExec.out_put_from_user.success && returnExec.out_put_from_user.output !== ''" class="flex flex-col gap-1">
          <label for="multiple-ac-1" class="font-bold mb-2 block">Saída da execução do código base</label>
          <Textarea v-model="returnExec.out_put_from_user.output" rows="5" cols="40" />
        </div>
        <div v-if="!returnExec.out_put_from_user.success && returnExec.out_put_from_user.error !== ''" class="flex flex-col gap-1">
          <label for="multiple-ac-1" class="font-bold mb-2 block">Erro da execução do código base</label>
          <Textarea v-model="returnExec.out_put_from_user.error" rows="5" cols="40" />
        </div>
      </div>
      <div class="flex flex-col gap-1">
        <label for="multiple-ac-1" class="font-bold mb-2 block">Código Teste</label>
        <MonacoEditor name="codigo_teste" v-model="valoresIniciais.codigo_teste" :options="optionsEditor" :lang="tipoDaLinguagem == 1 ? 'go' : 'python'" class="h-[600px]" />
        <!-- <Editor name="codigo_teste" v-model="valoresIniciais.codigo_teste" editorStyle="height: 220px" spellcheck="false" ref="refQuillCodigoTeste">
          <template v-slot:toolbar>
            <span class="ql-formats">
              <button v-tooltip.bottom="'Bold'" class="ql-bold"></button>
              <button v-tooltip.bottom="'Italic'" class="ql-italic"></button>
              <button v-tooltip.bottom="'Underline'" class="ql-underline"></button>
            </span>
          </template>
        </Editor> -->
        <Message v-if="$form.codigo_teste?.invalid" severity="error" size="small" variant="simple">{{ $form.codigo_teste.error?.message }}</Message>
      </div>
      <div class="flex flex-col gap-1 p-4">
        <div v-if="returnExec.out_put_unit_teste.success && returnExec.out_put_unit_teste.output !== ''" class="flex flex-col gap-1">
          <label for="multiple-ac-1" class="font-bold mb-2 block">Saída da execução do código teste</label>
          <Textarea v-model="returnExec.out_put_unit_teste.output" rows="5" cols="40" />
        </div>
        <div v-if="!returnExec.out_put_unit_teste.success && returnExec.out_put_unit_teste.error !== ''" class="flex flex-col gap-1">
          <label for="multiple-ac-1" class="font-bold mb-2 block">Erro da execução do código teste</label>
          <Textarea v-model="returnExec.out_put_unit_teste.error" rows="5" cols="40" />
        </div>
      </div>
      <div class="flex flex-col gap-1">
          <label for="multiple-ac-1" class="font-bold mb-2 block">Feedback</label>
          <div class="card">
              <DataTable :value="feedbacks" tableStyle="min-width: 50rem">
                  <Column field="descricao" header="Descrição">
                      <template #body="{ data }">
                          <div @click="editFeedback(data)" v-if="!data.isEditing">{{ data.descricao }}</div>
                          <InputText v-else v-model="data.descricao" @blur="saveFeedback(data)" />
                      </template>
                  </Column>
                  <Column class="w-24 !text-end">
                      <template #body="{ data }">
                          <Button @click="removerFeedback(data.id)" label="Apagar" severity="danger" outlined>Apagar</Button>
                      </template>
                  </Column>
              </DataTable>
          </div>
          <div class="flex justify-center my-1">
              <Button label="Adicionar feedbacks" @click="adicionarFeedback()"/>
          </div>
          <br/>
      </div>
      <Button type="submit" severity="secondary" label="Adicionar" />
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
