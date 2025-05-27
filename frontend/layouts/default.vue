<template>
  <div class="bg-zinc-900 min-h-screen text-gray-100">
    <header>
      <div class="w-full h-14 bg-zinc-950">
        <div class="p-1 flex justify-between">
          <NuxtLink to="/" class="flex flex-center">
            <img src="~/assets/images/logo_tcc.webp" class="h-12 w-12" alt="">
            <span class="justify-center text-3xl font-bold text-gray-300 p-1">TCC</span>
          </NuxtLink>
          <div></div>
          <div v-if="!isAuthenticated" class="flex flex-center justify-between p-3">
            <NuxtLink to="/account/login" class="mr-4">Entrar</NuxtLink>
            <NuxtLink to="/account/createaccount">Criar Conta</NuxtLink>
          </div>
          <div v-else class="flex flex-center justify-between p-3">
            <button @click="dashboardRedirec" class="mr-4">Dashboard</button>
            <button @click="logout" class="mr-4">Logout</button>
          </div>
        </div>
      </div>
    </header>

    <div class="w-full min-h-screen p-2">
      <slot />
    </div>

    <footer>
      <div class="w-full h-24 bg-zinc-950">
        <div class="flex justify-between">
            <div></div>
            <div class="flex justify-between p-3">
              <div class="p-3"></div>
              <NuxtLink to="/sitemap">Sitemap</NuxtLink>
            </div>
        </div>
        <hr class="ml-4 mr-4 border-gray-600"/>
        <div class="flex justify-center mt-2">TCC - Todos os direitos reservados</div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">

import { watch } from 'vue';
const isAuthenticated = ref(false)
const route = useRoute();

const { $authService } = useNuxtApp();

const dashboardRedirec = () => {
  if ($authService.isDocente()) {
    navigateTo("/dashboard_docente")
  } else if ($authService.isDiscente()) {
    navigateTo("/dashboard_discente")
  }
}

const logout = () => {
  navigateTo("/account/login")
  $authService.logout();
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