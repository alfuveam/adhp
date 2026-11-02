<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
    <div class="text-center mb-10 space-y-3">
      <span class="study-badge study-badge-indigo">Navegação Geral</span>
      <h1 class="text-3xl sm:text-4xl font-semibold text-slate-800 tracking-tight">Mapa do Site</h1>
      <p class="text-sm text-slate-600">
        Acesse rapidamente todas as rotas e seções disponíveis na plataforma ADHP
      </p>
    </div>

    <div class="study-card p-6 sm:p-8 border-[#dce1d9] bg-[#ffffff] shadow-none">
      <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
        <NuxtLink
          v-for="route in routes"
          :key="route"
          :to="route"
          class="p-4 rounded-lg bg-[#ffffff] border border-[#dce1d9] hover:border-[#285b75] hover:bg-[#ffffff] transition-all duration-200 flex items-center justify-between group"
        >
          <div class="flex items-center gap-3">
            <span
              class="w-8 h-8 rounded-lg bg-sky-900/5 border border-purple-500/20 text-sky-800 flex items-center justify-center text-sm group-hover:bg-sky-900/5 transition"
            >
              <i class="pi pi-link"></i>
            </span>
            <span
              class="text-sm font-semibold text-slate-700 group-hover:text-slate-800 transition font-mono"
            >
              {{ route }}
            </span>
          </div>

          <i
            class="pi pi-arrow-right text-sm text-slate-500 group-hover:text-sky-800 group-hover:translate-x-1 transition-all"
          ></i>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const routes = computed(() => {
  return router.getRoutes()
    .map(r => r.path)
    .filter(p => !p.includes(':') && p !== '/sitemap')
    .sort();
});

useSeoMeta({
  title: 'Mapa do Site • ADHP',
  description: 'Navegação e rotas do Ambiente de Desenvolvimento de Habilidades em Programação.'
});
</script>