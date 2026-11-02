<template>
  <Dialog 
    :visible="modelValue" 
    @update:visible="(value) => $emit('update:modelValue', value)" 
    modal
    :header="isUpdate ? 'Editar Trilha de Aprendizado' : 'Adicionar Nova Trilha'"
    :style="{ width: '32rem', maxWidth: '95vw' }"
    class="mimo-dialog"
  >
    <Toast />
    <ConfirmDialog />

    <Form 
      v-slot="$form" 
      :initialValues="valoresIniciais" 
      :resolver="resolver" 
      @submit="onFormSubmit" 
      class="flex flex-col gap-5 pt-2"
    >
      <!-- Título da Trilha -->
      <div class="flex flex-col gap-1.5">
        <label class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
          <i class="pi pi-tag text-purple-400 text-xs"></i>
          <span>Título da Trilha</span>
        </label>
        <InputText 
          name="titulo_trilha" 
          v-model="valoresIniciais.titulo_trilha" 
          type="text" 
          placeholder="Ex: Introdução a Algoritmos em Python" 
          class="w-full !p-3 !bg-[#15172a] !border-[#363e66] focus:!border-[#7c4ed4] rounded-xl text-sm"
          fluid 
        />
        <Message v-if="$form.titulo_trilha?.invalid" severity="error" size="small" variant="simple" class="!text-xs !text-rose-400 mt-1">
          {{ $form.titulo_trilha.error?.message }}
        </Message>
      </div>

      <!-- Tipo da Linguagem (Golang vs Python) -->
      <div class="flex flex-col gap-2">
        <label class="text-xs font-bold text-slate-300 uppercase tracking-wider flex items-center gap-1.5">
          <i class="pi pi-code text-purple-400 text-xs"></i>
          <span>Linguagem de Programação</span>
        </label>

        <RadioButtonGroup name="tipo_da_linguagem" v-model="valoresIniciais.tipo_da_linguagem" class="grid grid-cols-2 gap-3">
          <div 
            @click="valoresIniciais.tipo_da_linguagem = 1" 
            class="p-4 rounded-xl border cursor-pointer transition-all flex items-center gap-3.5"
            :class="valoresIniciais.tipo_da_linguagem === 1 ? 'bg-cyan-500/20 border-cyan-400 shadow-md shadow-cyan-500/10' : 'bg-[#15172a] border-[#363e66] hover:border-[#4f598c]'"
          >
            <RadioButton inputId="lang-go" :value="1" class="hidden" />
            <span class="w-9 h-9 rounded-xl bg-cyan-500/20 text-cyan-300 flex items-center justify-center font-extrabold text-xs">
              Go
            </span>
            <div>
              <div class="text-sm font-bold text-white">Golang</div>
              <div class="text-[11px] text-slate-400">Compilada</div>
            </div>
          </div>

          <div 
            @click="valoresIniciais.tipo_da_linguagem = 2" 
            class="p-4 rounded-xl border cursor-pointer transition-all flex items-center gap-3.5"
            :class="valoresIniciais.tipo_da_linguagem === 2 ? 'bg-amber-400/20 border-amber-400 shadow-md shadow-amber-500/10' : 'bg-[#15172a] border-[#363e66] hover:border-[#4f598c]'"
          >
            <RadioButton inputId="lang-py" :value="2" class="hidden" />
            <span class="w-9 h-9 rounded-xl bg-amber-400/20 text-amber-300 flex items-center justify-center font-extrabold text-xs">
              Py
            </span>
            <div>
              <div class="text-sm font-bold text-white">Python</div>
              <div class="text-[11px] text-slate-400">Interpretada</div>
            </div>
          </div>
        </RadioButtonGroup>

        <Message v-if="$form.tipo_da_linguagem?.invalid" severity="error" size="small" variant="simple" class="!text-xs !text-rose-400 mt-1">
          {{ $form.tipo_da_linguagem.error?.message }}
        </Message>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center justify-end gap-3 pt-3 border-t border-[#2d3458]">
        <button 
          type="button" 
          @click="closeDialog" 
          class="mimo-btn-secondary px-4 py-2 text-xs"
        >
          Cancelar
        </button>
        <button 
          type="submit" 
          class="mimo-btn-primary px-5 py-2 text-xs"
        >
          {{ isUpdate ? 'Salvar Alterações' : 'Criar Trilha' }}
        </button>
      </div>
    </Form>
  </Dialog>
</template>

<script setup lang="ts">
import Dialog from 'primevue/dialog';
import { ref, watch } from "vue";
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { useToast } from 'primevue/usetoast';

const toast = useToast();
const { $authService } = useNuxtApp();

const valoresIniciais = ref<{
  titulo_trilha: string;
  tipo_da_linguagem: number | null;
}>({
  titulo_trilha: '',
  tipo_da_linguagem: 1
});

const resolver = zodResolver(
  z.object({
    titulo_trilha: z.string().min(1, { message: 'Título da trilha é necessário.' }),
    tipo_da_linguagem: z.number().min(1, { message: 'Selecione a linguagem.' }).max(2, { message: 'Selecione a linguagem.' })
  })
);

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true
  },
  isUpdate: {
    type: Boolean,
    required: true
  },
  trilhaID: {
    type: String,
    required: true
  }
});

const emit = defineEmits(['update:modelValue', 'someEvent']);

const closeDialog = () => {
  valoresIniciais.value.titulo_trilha = "";
  emit('update:modelValue', false);
};

watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue) {
      valoresIniciais.value.titulo_trilha = "";
      if (props.isUpdate) {
        onGetTrilha();
      }
    }
  }
);

const onFormSubmit = async ({ valid }: { valid: boolean }) => {
  if (valid) {
    const config = useRuntimeConfig();
    
    if (props.isUpdate) {
      try {
        const response = await fetch(config.public.BACKEND_API_URL + '/v1/update_trilha', {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + $authService.getToken()
          },
          body: JSON.stringify({
            name: valoresIniciais.value.titulo_trilha,
            id: props.trilhaID,
            tipo_da_linguagem: Number(valoresIniciais.value.tipo_da_linguagem),
          })
        });

        if (response.ok) {
          toast.add({ severity: 'success', summary: 'Trilha atualizada.', life: 3000 });
          emit('someEvent');
          closeDialog();
        } else {
          toast.add({ severity: 'error', summary: 'Erro', detail: 'Falha ao atualizar trilha.', life: 3000 });
        }
      } catch (error) {
        console.error(error);
      }
    } else {
      try {
        const response = await fetch(config.public.BACKEND_API_URL + '/v1/add_trilha', {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            "Authorization": "Bearer " + $authService.getToken()
          },
          body: JSON.stringify({
            name: valoresIniciais.value.titulo_trilha,
            tipo_da_linguagem: Number(valoresIniciais.value.tipo_da_linguagem),
          })
        });

        if (response.ok) {
          toast.add({ severity: 'success', summary: 'Trilha criada com sucesso.', life: 3000 });
          emit('someEvent');
          closeDialog();
        } else {
          toast.add({ severity: 'error', summary: 'Erro', detail: 'Falha ao criar trilha.', life: 3000 });
        }
      } catch (error) {
        console.error(error);
      }
    }
  }
};

const onGetTrilha = async () => {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/get_trilha/' + props.trilhaID, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    });

    if (response.ok) {
      const resData = await response.json();
      valoresIniciais.value.titulo_trilha = resData.name;
      valoresIniciais.value.tipo_da_linguagem = resData.tipo_da_linguagem;
    } else {
      toast.add({ severity: 'error', summary: 'Erro', detail: 'Falha ao carregar a trilha.', life: 3000 });
    }
  } catch(error) {
    console.error(error);
  }
};
</script>