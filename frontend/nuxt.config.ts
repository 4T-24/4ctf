import yaml from "@rollup/plugin-yaml";
import Aura from "@primevue/themes/aura";

export default defineNuxtConfig({
  compatibilityDate: "2024-11-01",
  devtools: { enabled: true },
  css: ['~/assets/css/main.css'],

  vite: {
    plugins: [yaml()],
  },
  primevue: {
    options: {
      theme: {
        preset: Aura,
      },
    },
  },
  modules: [
  "@primevue/nuxt-module", "@pinia/nuxt", '@nuxtjs/tailwindcss', "@nuxt/icon"],
});