<template>
  <div class="card flex justify-center">
    <Dialog :visible="modelValue" @update:visible="(value) => $emit('update:modelValue', value)" 
      modal
      header="Adicionar trilha"
      :style="{ width: '25rem' }"
    >
      <div class="card flex justify-center">
        <Toast />
        <ConfirmDialog></ConfirmDialog>

        <Form v-slot="$form" :valoresIniciais :resolver @submit="onFormSubmit" class="flex flex-col gap-4 w-full sm:w-2/3">
          <div class="flex flex-col gap-1">
            <label for="multiple-ac-1" class="font-bold mb-2 block">Título trilha</label>
            <InputText name="titulo_trilha" v-model="valoresIniciais.titulo_trilha" type="text" placeholder="Título trilha" fluid />
            <Message v-if="$form.titulo_trilha?.invalid" severity="error" size="small" variant="simple">{{ $form.titulo_trilha.error?.message }}</Message>
          </div>
          <div class="flex flex-col gap-1">
            <label for="multiple-ac-1" class="font-bold mb-2 block">Tipo da linguagem</label>
            <div class="card flex justify-center">
              <div class="flex flex-col gap-2">
                <RadioButtonGroup name="tipo_da_linguagem" class="flex flex-wrap gap-4" v-model="valoresIniciais.tipo_da_linguagem">
                  <div class="flex items-center gap-2">
                    <RadioButton inputId="1" :value="1" />
                    <label for="1">Golang</label>
                  </div>
                  <div class="flex items-center gap-2">
                    <RadioButton inputId="2" :value="2" />
                    <label for="2">Python</label>
                  </div>
                </RadioButtonGroup>
                <Message v-if="$form.tipo_da_linguagem?.invalid" severity="error" size="small" variant="simple">{{ $form.tipo_da_linguagem.error?.message }}</Message>
              </div>
            </div>
          </div>
          <div class="flex justify-center gap-2">
            <Button type="button" label="Cancel" severity="secondary" @click="closeDialog"></Button>
            <Button type="submit" severity="secondary" :label="isUpdate ? 'Salvar' : 'Adicionar'" />
          </div>
        </Form>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import Dialog from 'primevue/dialog';
import { ref } from "vue";
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { useToast } from 'primevue/usetoast';

const toast = useToast();
const { $authService } = useNuxtApp()

const valoresIniciais = ref({
  titulo_trilha: '',
  tipo_da_linguagem: null
})

const resolver = zodResolver(
  z.object({
    titulo_trilha: z.string().min(1, { message: 'Título trilha é necessário.' }).nullable(),
    tipo_da_linguagem: z.number().min(1, { message: 'Tipo da linguagem é necessário.' }).max(2, { message: 'Tipo da linguagem é necessário.' }).nullable(),
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

const emit = defineEmits(['update:modelValue', 'someEvent'] );

const closeDialog = () => {
  valoresIniciais.value.titulo_trilha = ""
  emit('update:modelValue', false);
};

watch(
  () => props.modelValue,
  (newValue) => {
    if (newValue) {
      valoresIniciais.value.titulo_trilha = ""
      if (props.isUpdate) {
        onGetTrilha()
      }
    }
  }
)

const onFormSubmit = async ({ valid }: { valid: boolean }) => {
  if (valid) {
    toast.add({
      severity: 'success',
      summary: 'Form enviado.',
      life: 3000
    });

    const config = useRuntimeConfig()
    
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
        })
  
        if (response.ok) {
          toast.add({ severity: 'success', summary: 'Trilha atualizada.', life: 3000 });
          emit('someEvent')
          closeDialog();
        } else {
          toast.add({ severity: 'error', summary: 'Falha ao criar atualizada.', life: 3000 });
        }
      } catch (error) {
        console.error(error)
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
        })
  
        if (response.ok) {
          toast.add({ severity: 'success', summary: 'Trilha criada.', life: 3000 });
          emit('someEvent')
          closeDialog();
        } else {
          toast.add({ severity: 'error', summary: 'Falha ao criar Trilha.', life: 3000 });
        }
      } catch (error) {
        console.error(error)
      }
    }
    
  }
}

const onGetTrilha = async () => {
  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + '/v1/get_trilha/' + props.trilhaID, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        "Authorization": "Bearer " + $authService.getToken()
      }
    })

    if (response.ok) {
      const resData = await response.json()
      valoresIniciais.value.titulo_trilha = resData.name
      valoresIniciais.value.tipo_da_linguagem = resData.tipo_da_linguagem
    } else {
      toast.add({ severity: 'error', summary: 'Falha ao carregar a trilha.', life: 3000 });
    }

  } catch(error) {
    console.log(error)
  }
}
</script>