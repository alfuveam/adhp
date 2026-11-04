<template>
  <div class="min-h-[calc(100vh-8rem)] flex items-center justify-center p-4 sm:p-6">
    <Toast />

    <div class="w-full max-w-md study-card p-8 border-[#dce1d9] bg-[#ffffff] shadow-none space-y-6">
      <div class="space-y-2 text-center">
        <div
          class="inline-flex items-center justify-center w-12 h-12 rounded-lg bg-sky-900/5 border border-slate-200 text-sky-800 mb-1"
        >
          <i class="pi pi-lock text-lg"></i>
        </div>
        <h1 class="text-2xl font-semibold text-slate-800 tracking-tight">Bom ter você por aqui.</h1>
        <p class="text-sm text-slate-600">Entre na sua conta para continuar seus estudos.</p>
      </div>

      <Form
        v-slot="$form"
        :initialValues="initialValues"
        :resolver="resolver"
        @submit="onFormSubmit"
        class="flex flex-col gap-4"
      >
        <div class="flex flex-col gap-1.5">
          <label for="login-email" class="text-sm font-semibold text-slate-600">E-mail</label>
          <InputText
            id="login-email"
            autocomplete="email"
            v-model="initialValues.email"
            name="email"
            type="email"
            placeholder="seu.email@exemplo.com"
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

        <div class="flex flex-col gap-1.5">
          <label for="login-password" class="text-sm font-semibold text-slate-600">Senha</label>
          <Password
            inputId="login-password"
            autocomplete="current-password"
            v-model="initialValues.password"
            name="password"
            placeholder="Sua senha"
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
          <span>{{ submitting ? 'Entrando…' : 'Entrar na conta' }}</span>
          <i class="pi pi-arrow-right text-sm"></i>
        </button>
      </Form>

      <div class="pt-4 border-t border-[#dce1d9] text-center">
        <p class="text-sm text-slate-600">
          Está começando agora?
          <NuxtLink
            to="/account/createaccount"
            class="font-semibold text-sky-800 hover:text-sky-900 ml-1"
          >
            Cadastre-se aqui
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
  email: '',
  password: '',
});

const resolver = zodResolver(
  z.object({
    email: z.string().min(1, { message: 'Email é necessário.' }),
    password: z
      .string()
      .min(3, { message: 'Mínimo de 3 caracteres.' })
      .max(15, { message: 'Máximo de 15 caracteres.' }),
  })
);

const onFormSubmit = async (e: any) => {
  if (!e.valid || submitting.value) return;
  submitting.value = true;
  try {
    await sendLogin();
  } finally {
    submitting.value = false;
  }
};

async function sendLogin() {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + '/v1/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: initialValues.value.email,
        password: initialValues.value.password,
      }),
    });

    if (response.ok) {
      const lToken: LoginToken = await response.json();
      toast.add({ severity: 'success', summary: 'Login realizado com sucesso', life: 3000 });
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
      toast.add({
        severity: 'error',
        summary: 'Falha no login',
        detail: 'Verifique seu e-mail e senha.',
        life: 3000,
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
