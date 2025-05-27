<template>
  <div class="flex flex-col md:flex-row lg:flex-row items-top md:items-center pt-10 md:pt-0 justify-center min-h-screen w-full">
    <div class="grid justify-items-center w-5/5 h-auto md:w-5/5 lg:mb-96 2xl:mb-0 rounded space-y-0.2 p-3 md:p-5">
      <h1 class="mb-4 text-2xl font-extrabold">Criar Conta</h1>
      <div class="card flex justify-center mt-2">
        <Toast />

        <Form v-slot="$form" :initialValues="initialValues" :resolver="resolver" @submit="onFormSubmit" class="flex flex-col gap-4 w-full sm:w-60">
          <div class="flex flex-col gap-1">
              <InputText v-model="initialValues.completename" name="completename" type="text" placeholder="Nome completo" fluid />
              <Message v-if="$form.completename?.invalid" severity="error" size="small" variant="simple">{{ $form.completename.error.message }}</Message>
          </div>
          <div class="flex flex-col gap-1">
              <InputText v-model="initialValues.email" name="email" type="text" placeholder="Email" fluid />
              <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple">{{ $form.email.error.message }}</Message>
          </div>
          <IftaLabel>
            <label for="repeticao_espacada_minutos" style="font-size: 1rem;">Repetição Espacada</label>
          </IftaLabel>
          <span class="flex items-center gap-2"></span>
          <div class="flex flex-col items-center gap-2">
            <RadioButtonGroup name="repeticao_espacada_minutos" v-model="initialValues.repeticao_espacada_minutos" class="flex flex-wrap gap-4">
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
            <Message v-if="$form.repeticao_espacada_minutos?.invalid" severity="error" size="small" variant="simple">{{ $form.repeticao_espacada_minutos.error?.message }}</Message>
          </div>

          <div class="flex flex-col gap-1">
              <Password v-model="initialValues.password" name="password" placeholder="Password" :feedback="false" toggleMask fluid />
              <Message v-if="$form.password?.invalid" severity="error" size="small" variant="simple">
                  <ul class="my-0 px-4 flex flex-col gap-1">
                      <li v-for="(error, index) of $form.password.errors" :key="index">{{ error.message }}</li>
                  </ul>
              </Message>
          </div>
          <Button type="submit" severity="secondary" label="Submit" />
        </Form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { zodResolver } from '@primevue/forms/resolvers/zod';
import { z } from 'zod';
import { useToast } from 'primevue/usetoast';

const { $authService } = useNuxtApp();
const toast = useToast();

interface LoginToken {
  token: string
}

const initialValues = ref({
  completename: '',
  cpf: '',
  idade: '',
  sexo: '',
  experiencia_programacao: '',
  phone: '',
  email: '',
  repeticao_espacada_minutos: '',
  password: '',
});

const resolver = zodResolver(
    z.object({
        completename: z.string().min(1, { message: 'Completename é requerido.' }),
        cpf: z.string().min(1, { message: 'Cpf é requerido.' }),
        idade: z.string().min(1, { message: 'Idade é requerido.' }),
        sexo: z.string().min(1, { message: 'Sexo é requerido.' }),
        experiencia_programacao: z.string().min(1, { message: 'Experiencia_programacao é requerido.' }),
        phone: z.string().min(1, { message: 'Telefone é requerido.' }),
        email: z.string().min(1, { message: 'Email é requerido.' }),
        repeticao_espacada_minutos: z.string().min(1, { message: 'Tempo da repetição é requerido.' }),
        password: z
            .string()
            .min(9, { message: 'Mínimo 9 caracteres.' })
            .max(100, { message: 'Máximo 100 caracteres.' })
    })
);

const onFormSubmit = async (e: any) => {
    if (e.valid) {
        toast.add({ severity: 'success', summary: 'Form enviado.', life: 3000 });
        await createAccount();
    }
};

async function createAccount() {
  try {
    const config = useRuntimeConfig()
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/registeraccount', {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        completename: initialValues.value.completename,
        cpf: initialValues.value.cpf,
        idade: Number(initialValues.value.idade),
        sexo: Number(initialValues.value.sexo),
        experiencia_programacao: initialValues.value.experiencia_programacao == "true" ? true : false,
        phone: initialValues.value.phone,
        email: initialValues.value.email,
        repeticao_espacada_minutos: Number(initialValues.value.repeticao_espacada_minutos),
        password: initialValues.value.password
      })
    })

    if (response.ok) {
      const lToken: LoginToken = await response.json();
      toast.add({ severity: 'success', summary: 'Conta criada.', life: 3000 });

      initialValues.value.completename = ''
      initialValues.value.cpf = ''
      initialValues.value.idade = ''
      initialValues.value.sexo = ''
      initialValues.value.experiencia_programacao = ''
      initialValues.value.phone = ''
      initialValues.value.email = ''
      initialValues.value.repeticao_espacada_minutos = ''
      initialValues.value.password = ''

      await $authService.login(lToken.token);
      if ($authService.isAuthenticated()) {
        if ($authService.isDiscente()) {
          navigateTo('/dashboard_discente')
        } else {
          navigateTo('/dashboard_docente')
        }
      }
    } else {
      const errorData = await response.json();
      let textError = ''
      // Atualize os erros no formulário com base na resposta da API
      if (errorData.error_completename) {
        textError = errorData.error_completename
      }

      if (errorData.error_cpf) {
        textError = errorData.error_cpf
      }

      if (errorData.error_idade) {
        textError = errorData.error_idade
      }

      if (errorData.error_sexo) {
        textError = errorData.error_sexo
      }

      if (errorData.error_experiencia_programacao) {
        textError = errorData.error_experiencia_programacao
      }

      if (errorData.error_phone) {
        textError = errorData.error_phone
      }

      if (errorData.error_email) {
        textError = errorData.error_email
      }

      if (errorData.error_repeticao_espacada_minutos) {
        textError = errorData.error_repeticao_espacada_minutos
      }

      if (errorData.error_password) {
        textError = errorData.error_password
      }

      toast.add({ severity: 'error', summary: 'Erro ao criar à conta. ' + textError, life: 3000 });
    }
  } catch (error) {
    console.error(error)
  }
}
</script>
