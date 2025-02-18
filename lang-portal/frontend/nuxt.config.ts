// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { 
    enabled: true 
  },

  modules: [
    '@nuxtjs/tailwindcss'
  ],

  typescript: {
    strict: true
  },

  compatibilityDate: '2025-02-13'
})