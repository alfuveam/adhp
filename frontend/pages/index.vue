<template>
  <div class="flex flex-col md:flex-row items-center justify-center min-h-screen w-full p-4 md:p-8 max-w-7xl mx-auto">
    <div v-if="!isAuthenticated" class="bg-zinc-800 flex flex-col items-center justify-center w-full md:w-80 h-auto min-h-64 md:h-96 rounded-xl space-y-4 p-6 md:p-8 shadow-lg shrink-0">
      <NuxtLink to="/account/login">
        <Button class="w-36 md:w-40 h-12 rounded text-center" type="submit" label="Login" />
      </NuxtLink>
      <NuxtLink to="/account/createaccount">
        <Button class="w-36 md:w-40 h-12 rounded text-center" type="submit" label="Criar Conta" />
      </NuxtLink>
    </div>
    <div v-else class="bg-zinc-800 flex flex-col items-center justify-center w-full md:w-80 h-auto min-h-64 md:h-96 rounded-xl p-6 md:p-8 shadow-lg shrink-0">
      <Button class="w-36 md:w-40 h-12 rounded text-center" type="submit" label="Dashboard" @click="redirectToDashboard" />
    </div>
    <div class="bg-zinc-800 flex items-center justify-center w-full h-auto min-h-64 md:h-96 rounded-xl p-6 md:p-10 grow ml-0 md:ml-4 my-4 md:my-0 shadow-lg">
      <p class="text-white text-center md:text-left text-base md:text-lg leading-relaxed">
        <strong>Bem Vindo!</strong><br><br>
        Essa ferramenta tem a função de ajudar o instrutor, retirando o trabalho repetitivo de corrigir questões de programação de alunos que estão iniciando na área de programação. Além de ajudar os alunos com a utilização de feedback automatizado para realizar o termino dos exercicios de forma mais eficaz, sendo então explorado o uso da repetição espaçada.
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