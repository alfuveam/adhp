<template>
  <div class="min-h-[calc(100vh-8rem)] flex items-center justify-center p-4 sm:p-6">
    <Toast />

    <div class="w-full max-w-lg study-card p-8 border-[#dce1d9] bg-[#ffffff] shadow-none space-y-6">
      <div class="space-y-2 text-center">
        <div
          class="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-sky-900/5 border border-slate-200 text-sky-800 mb-1"
        >
          <i class="pi pi-user-plus text-lg"></i>
        </div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Criar Conta de Aluno</h1>
        <p class="text-sm text-slate-600">
          Acesse os desafios práticos das trilhas e acompanhe sua curva de retenção
        </p>
      </div>

      <Form
        v-slot="$form"
        :initialValues="initialValues"
        :resolver="resolver"
        @submit="onFormSubmit"
        class="flex flex-col gap-4"
      >
        <div class="flex flex-col gap-1.5">
          <label for="createaccount-completename" class="text-sm font-semibold text-slate-600"
            >Nome Completo</label
          >
          <InputText
            id="createaccount-completename"
            autocomplete="name"
            v-model="initialValues.completename"
            name="completename"
            type="text"
            placeholder="Nome e Sobrenome"
            class="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-sm"
            fluid
          />
          <Message
            v-if="$form.completename?.invalid"
            severity="error"
            size="small"
            variant="simple"
            class="!text-sm !text-rose-700 mt-1"
          >
            {{ $form.completename.error.message }}
          </Message>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="createaccount-email" class="text-sm font-semibold text-slate-600"
            >E-mail</label
          >
          <InputText
            v-model="initialValues.email"
            id="createaccount-email"
            autocomplete="email"
            name="email"
            type="email"
            placeholder="aluno@universidade.edu.br"
            class="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-sm"
            fluid
          />
          <Message
            v-if="$form.email?.invalid"
            severity="error"
            size="small"
            variant="simple"
            class="!text-sm !text-rose-700 mt-1"
          >
            {{ $form.email.error.message }}
          </Message>
        </div>

        <div class="flex flex-col gap-2">
          <div class="flex items-center justify-between">
            <span id="review-interval-label" class="text-sm font-semibold text-slate-600"
              >Intervalo de Repetição Espaçada</span
            >
            <span class="text-[11px] text-sky-800 font-semibold font-mono">Curva Ebbinghaus</span>
          </div>

          <RadioButtonGroup
            aria-labelledby="review-interval-label"
            name="repeticao_espacada_minutos"
            v-model="initialValues.repeticao_espacada_minutos"
            class="grid grid-cols-2 sm:grid-cols-5 gap-2"
          >
            <label
              for="rep-1"
              class="flex flex-col items-center justify-center p-2.5 rounded-lg border cursor-pointer transition text-center focus-within:ring-2 focus-within:ring-sky-800 focus-within:ring-offset-2"
              :class="
                initialValues.repeticao_espacada_minutos === '1'
                  ? 'bg-[#285b75] border-[#285b75] text-white font-semibold border-b-2 border-b-[#204c63]'
                  : 'bg-[#ffffff] border-[#dce1d9] hover:border-[#dce1d9] text-slate-600'
              "
            >
              <RadioButton inputId="rep-1" value="1" class="sr-only" />
              <span class="text-sm font-semibold">1 Hora</span>
              <span class="text-[10px] opacity-75">Curto</span>
            </label>

            <label
              for="rep-2"
              class="flex flex-col items-center justify-center p-2.5 rounded-lg border cursor-pointer transition text-center focus-within:ring-2 focus-within:ring-sky-800 focus-within:ring-offset-2"
              :class="
                initialValues.repeticao_espacada_minutos === '2'
                  ? 'bg-[#285b75] border-[#285b75] text-white font-semibold border-b-2 border-b-[#204c63]'
                  : 'bg-[#ffffff] border-[#dce1d9] hover:border-[#dce1d9] text-slate-600'
              "
            >
              <RadioButton inputId="rep-2" value="2" class="sr-only" />
              <span class="text-sm font-semibold">9 Horas</span>
              <span class="text-[10px] opacity-75">Mesmo dia</span>
            </label>

            <label
              for="rep-3"
              class="flex flex-col items-center justify-center p-2.5 rounded-lg border cursor-pointer transition text-center focus-within:ring-2 focus-within:ring-sky-800 focus-within:ring-offset-2"
              :class="
                initialValues.repeticao_espacada_minutos === '3'
                  ? 'bg-[#285b75] border-[#285b75] text-white font-semibold border-b-2 border-b-[#204c63]'
                  : 'bg-[#ffffff] border-[#dce1d9] hover:border-[#dce1d9] text-slate-600'
              "
            >
              <RadioButton inputId="rep-3" value="3" class="sr-only" />
              <span class="text-sm font-semibold">1 Dia</span>
              <span class="text-[10px] opacity-75">Padrão</span>
            </label>

            <label
              for="rep-4"
              class="flex flex-col items-center justify-center p-2.5 rounded-lg border cursor-pointer transition text-center focus-within:ring-2 focus-within:ring-sky-800 focus-within:ring-offset-2"
              :class="
                initialValues.repeticao_espacada_minutos === '4'
                  ? 'bg-[#285b75] border-[#285b75] text-white font-semibold border-b-2 border-b-[#204c63]'
                  : 'bg-[#ffffff] border-[#dce1d9] hover:border-[#dce1d9] text-slate-600'
              "
            >
              <RadioButton inputId="rep-4" value="4" class="sr-only" />
              <span class="text-sm font-semibold">6 Dias</span>
              <span class="text-[10px] opacity-75">Semanal</span>
            </label>

            <label
              for="rep-5"
              class="flex flex-col items-center justify-center p-2.5 rounded-lg border cursor-pointer transition text-center focus-within:ring-2 focus-within:ring-sky-800 focus-within:ring-offset-2"
              :class="
                initialValues.repeticao_espacada_minutos === '5'
                  ? 'bg-[#285b75] border-[#285b75] text-white font-semibold border-b-2 border-b-[#204c63]'
                  : 'bg-[#ffffff] border-[#dce1d9] hover:border-[#dce1d9] text-slate-600'
              "
            >
              <RadioButton inputId="rep-5" value="5" class="sr-only" />
              <span class="text-sm font-semibold">31 Dias</span>
              <span class="text-[10px] opacity-75">Mensal</span>
            </label>
          </RadioButtonGroup>

          <Message
            v-if="$form.repeticao_espacada_minutos?.invalid"
            severity="error"
            size="small"
            variant="simple"
            class="!text-sm !text-rose-700 mt-1"
          >
            {{ $form.repeticao_espacada_minutos.error?.message }}
          </Message>
        </div>

        <div class="flex flex-col gap-1.5">
          <label for="createaccount-password" class="text-sm font-semibold text-slate-600"
            >Senha (mínimo 9 caracteres)</label
          >
          <Password
            inputId="createaccount-password"
            autocomplete="new-password"
            v-model="initialValues.password"
            name="password"
            placeholder="Digite uma senha forte"
            :feedback="false"
            toggleMask
            class="w-full"
            inputClass="w-full !p-3 !bg-[#ffffff] !border-[#dce1d9] focus:!border-[#285b75] rounded-lg text-sm"
            fluid
          />
          <Message
            v-if="$form.password?.invalid"
            severity="error"
            size="small"
            variant="simple"
            class="!text-sm !text-rose-700 mt-1"
          >
            <ul class="my-0 px-2 flex flex-col gap-0.5">
              <li v-for="(error, index) of $form.password.errors" :key="index">
                {{ error.message }}
              </li>
            </ul>
          </Message>
        </div>

        <button
          type="submit"
          :disabled="submitting"
          :aria-busy="submitting"
          class="study-btn-primary w-full mt-2 py-3 rounded-lg text-sm font-semibold"
        >
          <span>{{ submitting ? 'Criando sua conta…' : 'Criar minha conta' }}</span>
          <i class="pi pi-arrow-right text-sm"></i>
        </button>
      </Form>

      <div class="pt-4 border-t border-[#dce1d9] text-center">
        <p class="text-sm text-slate-600">
          Já tem conta cadastrada?
          <NuxtLink to="/account/login" class="font-semibold text-sky-800 hover:text-sky-900 ml-1">
            Faça login aqui
          </NuxtLink>
        </p>
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
const submitting = ref(false);

interface LoginToken {
  token: string;
}

const initialValues = ref({
  completename: '',
  email: '',
  repeticao_espacada_minutos: '3', // 1 Dia default
  password: '',
  cpf: '000.000.000-00',
  idade: '20',
  sexo: '1',
  experiencia_programacao: 'false',
  phone: '000000000',
});

const resolver = zodResolver(
  z.object({
    completename: z
      .string()
      .min(5, { message: 'O nome completo deve ter entre 5 e 50 caracteres.' })
      .max(50, { message: 'Máximo de 50 caracteres.' }),
    email: z.string().email({ message: 'Insira um email válido.' }),
    repeticao_espacada_minutos: z.string().min(1, { message: 'Selecione o tempo da repetição.' }),
    password: z
      .string()
      .min(9, { message: 'A senha deve conter no mínimo 9 caracteres.' })
      .max(100, { message: 'Máximo 100 caracteres.' }),
  })
);

const onFormSubmit = async (e: any) => {
  if (!e.valid || submitting.value) return;
  submitting.value = true;
  try {
    await createAccount();
  } finally {
    submitting.value = false;
  }
};

async function createAccount() {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/registeraccount', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        completename: initialValues.value.completename,
        cpf: initialValues.value.cpf || '000.000.000-00',
        idade: Number(initialValues.value.idade) || 20,
        sexo: Number(initialValues.value.sexo) || 1,
        experiencia_programacao: initialValues.value.experiencia_programacao === 'true',
        phone: initialValues.value.phone || '000000000',
        email: initialValues.value.email,
        repeticao_espacada_minutos: Number(initialValues.value.repeticao_espacada_minutos),
        password: initialValues.value.password,
      }),
    });

    if (response.ok) {
      const lToken: LoginToken = await response.json();
      toast.add({ severity: 'success', summary: 'Conta criada com sucesso!', life: 3000 });

      initialValues.value.completename = '';
      initialValues.value.email = '';
      initialValues.value.password = '';

      await $authService.login(lToken.token);
      if ($authService.isAuthenticated()) {
        if ($authService.isDiscente()) {
          navigateTo('/dashboard_discente');
        } else {
          navigateTo('/dashboard_docente');
        }
      }
    } else {
      const errorData = await response.json().catch(() => ({}));
      let textError =
        errorData.error_completename ||
        errorData.error_email ||
        errorData.error_password ||
        errorData.error_repeticao_espacada_minutos ||
        errorData.error ||
        'Verifique os dados informados.';

      toast.add({
        severity: 'error',
        summary: 'Erro ao criar conta',
        detail: textError,
        life: 4000,
      });
    }
  } catch (error) {
    console.error(error);
    toast.add({
      severity: 'error',
      summary: 'Erro de conexão',
      detail: 'Não foi possível conectar ao servidor.',
      life: 3000,
    });
  }
}
</script>
