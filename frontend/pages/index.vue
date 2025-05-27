<template>
  <div class="flex flex-col md:flex-row lg:flex-row items-top md:items-center pt-10 md:pt-0 justify-center min-h-screen w-full">
    <div v-if="!isAuthenticated" class="bg-zinc-800 grid justify-items-center w-5/5 h-96 md:w-5/5 md:h-96 lg:mb-96 2xl:mb-0 rounded space-y-0.2 p-28">
      <NuxtLink to="/account/login">
        <Button class="w-32 md:w-40 h-12 rounded text-center p-3" type="submit" label="Login" />
      </NuxtLink>
      <NuxtLink to="/account/createaccount">
        <Button class="w-32 md:w-40 h-12 rounded text-center p-3" type="submit" label="Criar Conta" />
      </NuxtLink>
    </div>
    <div v-else  class="bg-zinc-800 grid justify-items-center w-5/5 h-96 md:w-5/5 md:h-96 lg:mb-96 2xl:mb-0 rounded space-y-0.2 p-28">
      <Button class="w-32 md:w-40 h-12 rounded text-center p-3" type="submit" label="Dashboard" @click="redirectToDashboard" />
    </div>
    <div class="bg-zinc-800 grid justify-items-center w-5/5 h-96 md:w-5/5 md:h-96 lg:mb-96 2xl:mb-0 rounded space-y-0.2 p-28 grow overflow-y-auto ml-0 md:ml-2 lg:ml-2 my-2 md:my-0 lg:my-0">
      <p class="text-white text-center justify">
        Bem Vindo! <br>
        Essa ferramenta tem a função de ajudar o instrutor, retirando o trabalho repetitivo de corrigir questões de programação de alunos que estão iniciando na área de programação.
        Além de ajudar os alunos com a utilização de feedback automatizado para realizar o termino do exercicios de forma mais eficaz.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch } from 'vue';
const route = useRoute();
const isAuthenticated = ref(false)
const { $authService } = useNuxtApp();

const redirectToDashboard = () => {
  if ($authService.isDocente()) {
    navigateTo("/dashboard_docente")
  } else if ($authService.isDiscente()) {
    navigateTo("/dashboard_discente")
  }
}

const checkIsAuthenticated = () => {
  if ($authService.isAuthenticated()) {
    const jwtPayload = $authService.parseJwt($authService.getToken() || '');
    if (jwtPayload.exp < Date.now() / 1000) {
      $authService.logout();
    }
  }
  isAuthenticated.value = $authService.isAuthenticated();
}

watch(route, (to, from) => {
  checkIsAuthenticated()
});

onMounted(() => {
  checkIsAuthenticated()
})

</script>