<template>
  <div class="bg-rust-700 min-h-screen py-12">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
    <h1 class="text-4xl font-bold mb-6 font-serif text-white text-center">
      Sitemap
    </h1>

    <ul class="list-none pl-0 text-lg">
      <li v-for="route in routes" :key="route" class="mb-3">
        <router-link
          :to="route"
          class="flex items-center transition transform hover:translate-x-2"
        >
          <span class="mr-2">&#9658;</span> {{ route }}
        </router-link>
      </li>
    </ul>  
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRouter } from 'vue-router';

const router = useRouter();
const routes = computed(() => router.getRoutes().map(route => route.path));

useSeoMeta({
  title: 'Sitemap',
})

</script>