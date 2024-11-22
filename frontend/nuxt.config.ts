import yaml from '@rollup/plugin-yaml'


export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: true },

  vite: {
    plugins: [
      yaml()
    ]
  },

  modules: ['@nuxt/ui']
})

