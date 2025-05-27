<template>
  <div class="flex flex-col md:flex-row lg:flex-row items-top md:items-center pt-10 md:pt-0 justify-center min-h-screen w-full">
    <div class="grid justify-items-center w-5/5 h-auto 2xl:mb-0 rounded space-y-0.2 p-3 md:p-5">
      <h1 class="mb-4 text-2xl font-extrabold">Login</h1>
      <div class="card flex justify-center mt-2">
        <Toast />

        <Form v-slot="$form" :initialValues="initialValues" :resolver="resolver" @submit="onFormSubmit" class="flex flex-col gap-4 w-full sm:w-60">
          <div class="flex flex-col gap-1">
              <InputText v-model="initialValues.email" name="email" type="text" placeholder="Email" fluid />
              <Message v-if="$form.email?.invalid" severity="error" size="small" variant="simple">{{ $form.email.error.message }}</Message>
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
    email: '',
    password: ''
});

const resolver = zodResolver(
    z.object({
        email: z.string().min(1, { message: 'Email é necessário.' }),
        password: z
            .string()
            .min(3, { message: 'Mínimo de 3 characters.' })
            .max(15, { message: 'Máximo de 15 characters.' })
    })
);

const onFormSubmit = async (e: any) => {
    if (e.valid) {
        toast.add({ severity: 'success', summary: 'Form enviado.', life: 3000 });
    }

    await sendLogin()
};

async function sendLogin() {
  try {
    const response = await fetch(useRuntimeConfig().public.BACKEND_API_URL + "/v1/login", {
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
      toast.add({ severity: 'success', summary: 'Login aprovado.', life: 3000 });
      initialValues.value.email = ''
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
      toast.add({ severity: 'error', summary: 'Login falhou, verifique as credenciais.', life: 3000 });
    }
  } catch (error) {
    console.error(error)
  }
}
</script>
