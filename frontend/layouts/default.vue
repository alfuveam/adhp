<template>
  <div class="site-shell">
    <a class="skip-link" href="#conteudo">Pular para o conteúdo</a>
    <header class="site-header">
      <div class="site-nav">
        <NuxtLink to="/" class="site-brand" aria-label="ADHP — início">
          <img src="/icon.svg" alt="" width="38" height="38" />
          <span
            >ADHP<span class="brand-dot">.</span
            ><small>Um espaço para aprender programação</small></span
          >
        </NuxtLink>
        <nav aria-label="Navegação principal" class="site-links">
          <NuxtLink to="/#metodologia">Como funciona</NuxtLink>
          <NuxtLink to="/#trilhas">Exercícios</NuxtLink>
          <template v-if="!isAuthenticated">
            <NuxtLink to="/account/login">Entrar</NuxtLink>
            <NuxtLink to="/account/createaccount" class="button">Criar conta</NuxtLink>
          </template>
          <template v-else>
            <button class="button" @click="dashboardRedirec">Meus estudos</button>
            <button class="nav-logout" @click="logout">Sair</button>
          </template>
        </nav>
      </div>
    </header>
    <main id="conteudo" class="site-main"><slot /></main>
    <footer class="site-footer">
      <div>
        <span>ADHP · Ambiente para o Desenvolvimento de Habilidades em Programação</span>
      </div>
      <NuxtLink to="/sitemap">Mapa do site</NuxtLink>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';

const isAuthenticated = ref(false);
const userRole = ref<string | null>(null);
const route = useRoute();
const { $authService } = useNuxtApp();

const dashboardRedirec = () => {
  if ($authService.isDocente()) {
    navigateTo("/dashboard_docente");
  } else if ($authService.isDiscente()) {
    navigateTo("/dashboard_discente");
  }
};

const logout = () => {
  $authService.logout();
  isAuthenticated.value = false;
  userRole.value = null;
  navigateTo("/account/login");
};

const checkIsAuthenticated = () => {
  if ($authService.isAuthenticated()) {
    const token = $authService.getToken() || '';
    try {
      const jwtPayload = $authService.parseJwt(token);
      if (jwtPayload && jwtPayload.exp && jwtPayload.exp < Date.now() / 1000) {
        $authService.logout();
      }
    } catch {
      // ignore parse error
    }
  }

  const isAuth = $authService.isAuthenticated();
  isAuthenticated.value = isAuth;
  if (isAuth) {
    if ($authService.isDocente()) {
      userRole.value = 'Docente';
    } else if ($authService.isDiscente()) {
      userRole.value = 'Discente';
    } else {
      userRole.value = 'Usuário';
    }
  } else {
    userRole.value = null;
  }
};

watch(route, () => {
  checkIsAuthenticated();
});

onMounted(() => {
  checkIsAuthenticated();
});
</script>