<template>
  <div class="min-h-screen bg-[#0e1020] text-slate-100 flex flex-col font-sans selection:bg-[#7c4ed4]/40 selection:text-white">
    <!-- Header / Navbar (Mimo Style) -->
    <header class="sticky top-0 z-50 bg-[#13162b]/95 backdrop-blur-md border-b border-[#232847]">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          
          <!-- Logo & Brand -->
          <NuxtLink to="/" class="flex items-center gap-3 group focus:outline-none">
            <div class="relative">
              <img 
                src="~/assets/images/logo_adhp.webp" 
                class="h-9 w-9 rounded-xl object-cover border border-[#3b436d] shadow-md shadow-black/40 group-hover:border-[#7c4ed4] transition-all" 
                alt="ADHP Logo"
              >
            </div>
            <div class="flex flex-col">
              <div class="flex items-center gap-2">
                <span class="text-lg font-extrabold tracking-tight text-white group-hover:text-[#a78bfa] transition-colors">
                  ADHP
                </span>
                <span class="px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider rounded-full bg-[#1e233d] text-purple-300 border border-[#39426f]">
                  UNESC
                </span>
              </div>
              <span class="hidden sm:inline-block text-[11px] text-slate-400 font-medium -mt-0.5">
                Ambiente de Habilidades em Programação
              </span>
            </div>
          </NuxtLink>

          <!-- Navigation Links (Center) -->
          <nav class="hidden md:flex items-center gap-1.5">
            <NuxtLink 
              to="/" 
              class="px-3.5 py-1.5 rounded-xl text-sm font-semibold text-slate-300 hover:text-white hover:bg-[#1d223f] transition-all"
            >
              Início
            </NuxtLink>
            <NuxtLink 
              to="/#trilhas" 
              class="px-3.5 py-1.5 rounded-xl text-sm font-semibold text-slate-300 hover:text-white hover:bg-[#1d223f] transition-all"
            >
              Trilhas
            </NuxtLink>
            <NuxtLink 
              to="/#metodologia" 
              class="px-3.5 py-1.5 rounded-xl text-sm font-semibold text-slate-300 hover:text-white hover:bg-[#1d223f] transition-all"
            >
              Metodologia
            </NuxtLink>
            <NuxtLink 
              to="/sitemap" 
              class="px-3.5 py-1.5 rounded-xl text-sm font-semibold text-slate-400 hover:text-white hover:bg-[#1d223f] transition-all"
            >
              Mapa do Site
            </NuxtLink>
          </nav>

          <!-- Auth Actions (Right) -->
          <div class="flex items-center gap-3">
            <template v-if="!isAuthenticated">
              <NuxtLink 
                to="/account/login" 
                class="px-3.5 py-2 rounded-xl text-sm font-semibold text-slate-300 hover:text-white hover:bg-[#1d223f] transition"
              >
                Entrar
              </NuxtLink>
              <NuxtLink 
                to="/account/createaccount" 
                class="mimo-btn-primary px-4 py-2 text-xs font-bold uppercase tracking-wider h-9"
              >
                <span>Criar Conta</span>
                <i class="pi pi-arrow-right text-[10px]"></i>
              </NuxtLink>
            </template>

            <template v-else>
              <!-- Role badge -->
              <span 
                v-if="userRole" 
                class="hidden sm:inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold"
                :class="userRole === 'Docente' ? 'mimo-badge-purple' : 'mimo-badge-emerald'"
              >
                <i :class="userRole === 'Docente' ? 'pi pi-user-edit' : 'pi pi-graduation-cap'" class="text-[11px]"></i>
                {{ userRole }}
              </span>

              <!-- Dashboard action -->
              <button 
                @click="dashboardRedirec" 
                class="mimo-btn-secondary px-3.5 py-2 text-xs font-bold uppercase tracking-wider h-9 flex items-center gap-1.5"
              >
                <i class="pi pi-th-large text-xs text-purple-300"></i>
                <span>Painel</span>
              </button>

              <!-- Logout -->
              <button 
                @click="logout" 
                class="p-2 sm:px-3 sm:py-2 rounded-xl text-sm font-semibold text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 transition"
                title="Sair"
              >
                <i class="pi pi-sign-out text-sm"></i>
                <span class="hidden sm:inline-block ml-1">Sair</span>
              </button>
            </template>
          </div>

        </div>
      </div>
    </header>

    <!-- Main Content Slot -->
    <main class="flex-1 w-full">
      <slot />
    </main>

    <!-- Clean, High-Contrast Footer -->
    <footer class="border-t border-[#262b48] bg-[#0e101f] mt-auto">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          
          <!-- Column 1: Brand & Academic Info -->
          <div class="space-y-3 md:col-span-1">
            <div class="flex items-center gap-2.5">
              <img src="~/assets/images/logo_adhp.webp" class="h-7 w-7 rounded object-cover" alt="ADHP">
              <span class="text-base font-bold text-white tracking-tight">ADHP</span>
            </div>
            <p class="text-xs text-slate-400 leading-relaxed">
              Ambiente para o Desenvolvimento de Habilidades em Programação. Prática ativa de exercícios com feedback de testes e repetição espaçada.
            </p>
            <div class="text-xs text-slate-500">
              UNESC • Curso de Ciência da Computação
            </div>
          </div>

          <!-- Column 2: Trilhas -->
          <div>
            <h4 class="text-xs font-semibold text-slate-300 uppercase tracking-wider mb-3">Linguagens</h4>
            <ul class="space-y-2 text-xs text-slate-400">
              <li>Golang (Compilada &amp; Concorrente)</li>
              <li>Python (Algoritmos &amp; Estruturas)</li>
              <li>Testes Unitários Automatizados</li>
              <li>Curva de Ebbinghaus</li>
            </ul>
          </div>

          <!-- Column 3: Navegação -->
          <div>
            <h4 class="text-xs font-semibold text-slate-300 uppercase tracking-wider mb-3">Acesso</h4>
            <ul class="space-y-2 text-xs text-slate-400">
              <li><NuxtLink to="/" class="hover:text-slate-200 transition">Página Inicial</NuxtLink></li>
              <li><NuxtLink to="/account/login" class="hover:text-slate-200 transition">Acessar Conta</NuxtLink></li>
              <li><NuxtLink to="/account/createaccount" class="hover:text-slate-200 transition">Cadastro de Aluno</NuxtLink></li>
              <li><NuxtLink to="/sitemap" class="hover:text-slate-200 transition">Mapa do Site</NuxtLink></li>
            </ul>
          </div>

          <!-- Column 4: Ambiente -->
          <div>
            <h4 class="text-xs font-semibold text-slate-300 uppercase tracking-wider mb-3">Tecnologias</h4>
            <div class="flex flex-wrap gap-1.5 text-xs text-slate-400">
              <span class="px-2.5 py-1 rounded-lg bg-[#1e233d] border border-[#363e66] text-slate-300 font-mono text-[11px]">Nuxt 3</span>
              <span class="px-2.5 py-1 rounded-lg bg-[#1e233d] border border-[#363e66] text-slate-300 font-mono text-[11px]">Vue 3</span>
              <span class="px-2.5 py-1 rounded-lg bg-[#1e233d] border border-[#363e66] text-slate-300 font-mono text-[11px]">Monaco Editor</span>
              <span class="px-2.5 py-1 rounded-lg bg-[#1e233d] border border-[#363e66] text-slate-300 font-mono text-[11px]">Go</span>
              <span class="px-2.5 py-1 rounded-lg bg-[#1e233d] border border-[#363e66] text-slate-300 font-mono text-[11px]">Python</span>
              <span class="px-2.5 py-1 rounded-lg bg-[#1e233d] border border-[#363e66] text-slate-300 font-mono text-[11px]">PostgreSQL</span>
            </div>
          </div>

        </div>

        <div class="mt-8 pt-6 border-t border-[#262b48] flex flex-col sm:flex-row items-center justify-between gap-3 text-xs text-slate-500">
          <div>
            &copy; 2025–2026 ADHP. Universidade do Extremo Sul Catarinense.
          </div>
          <div class="text-slate-500">
            Trabalho de Conclusão de Curso
          </div>
        </div>
      </div>
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