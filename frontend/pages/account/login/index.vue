<template>
  <div class="min-h-[calc(100vh-8rem)] flex items-center justify-center p-4 sm:p-6">
    <Toast />

    <!-- Login Card (Mimo Style) -->
    <div class="w-full max-w-md mimo-card p-8 border-[#2c335a] bg-[#161a33] shadow-2xl space-y-6">
      
      <!-- Card Header -->
      <div class="space-y-2 text-center">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-2xl bg-purple-500/20 border border-purple-500/30 text-purple-300 mb-1">
          <i class="pi pi-lock text-lg"></i>
        </div>
        <h1 class="text-2xl font-extrabold text-white tracking-tight">
          Acessar Plataforma
        </h1>
        <p class="text-xs text-slate-300">
          Ambiente integrado para discentes e docentes da UNESC
        </p>
      </div>

      <!-- Form -->
      <Form 
        v-slot="$form" 
        :initialValues="initialValues" 
        :resolver="resolver" 
        @submit="onFormSubmit" 
        class="flex flex-col gap-4"
      >
        <!-- Email field -->
        <div class="flex flex-col gap-1.5">
          <label class="text-xs font-bold text-slate-300">Email Acadêmico ou Pessoal</label>
          <InputText 
            v-model="initialValues.email" 
            name="email" 
            type="email" 
            placeholder="seu.email@exemplo.com" 
            class="w-full !p-3 !bg-[#0f1122] !border-[#2a3154] focus:!border-[#7c4ed4] rounded-xl text-sm"
            fluid 
          />
          <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple" class="!text-xs !text-rose-400 mt-1">
            {{ $form.email.error.message }}
          </Message>
        </div>

        <!-- Password field -->
        <div class="flex flex-col gap-1.5">
          <label class="text-xs font-bold text-slate-300">Senha</label>
          <Password 
            v-model="initialValues.password" 
            name="password" 
            placeholder="Sua senha" 
            :feedback="false" 
            toggleMask 
            class="w-full"
            inputClass="w-full !p-3 !bg-[#0f1122] !border-[#2a3154] focus:!border-[#7c4ed4] rounded-xl text-sm"
            fluid 
          />
          <Message v-if="$form.password?.invalid" severity="error" size="small" variant="simple" class="!text-xs !text-rose-400 mt-1">
            <ul class="my-0 px-2 flex flex-col gap-0.5">
              <li v-for="(error, index) of $form.password.errors" :key="index">{{ error.message }}</li>
            </ul>
          </Message>
        </div>

        <!-- Submit Button -->
        <button 
          type="submit" 
          class="mimo-btn-primary w-full mt-2 py-3 rounded-xl text-sm font-bold"
        >
          <span>Entrar na Conta</span>
          <i class="pi pi-arrow-right text-xs"></i>
        </button>
      </Form>

      <!-- Bottom Switch to Register -->
      <div class="pt-4 border-t border-[#252b4c] text-center">
        <p class="text-xs text-slate-300">
          Ainda não tem conta de aluno?
          <NuxtLink to="/account/createaccount" class="font-bold text-purple-300 hover:text-purple-200 ml-1">
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

interface LoginToken {
  token: string;
}

const initialValues = ref({
  email: '',
  password: ''
});

const resolver = zodResolver(
  z.object({
    email: z.string().min(1, { message: 'Email é necessário.' }),
    password: z
      .string()
      .min(3, { message: 'Mínimo de 3 caracteres.' })
      .max(15, { message: 'Máximo de 15 caracteres.' })
  })
);

const onFormSubmit = async (e: any) => {
  if (e.valid) {
    toast.add({ severity: 'info', summary: 'Autenticando...', life: 2000 });
  }
  await sendLogin();
};

async function sendLogin() {
  try {
    const config = useRuntimeConfig();
    const response = await fetch(config.public.BACKEND_API_URL + "/v1/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
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
      toast.add({ severity: 'error', summary: 'Falha no login', detail: 'Verifique seu e-mail e senha.', life: 3000 });
    }
  } catch (error) {
    console.error(error);
    toast.add({ severity: 'error', summary: 'Erro de conexão', detail: 'Não foi possível conectar ao servidor.', life: 3000 });
  }
}
</script>
