<template>
  <div 
    class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" 
    role="alert"
  >
    <strong class="font-bold block mb-2">{{ title }}</strong>
    <p class="text-sm">{{ message }}</p>
    <button 
      v-if="onRetry" 
      @click="onRetry"
      class="mt-2 bg-red-500 text-white px-3 py-1 rounded text-xs hover:bg-red-600 transition"
    >
      Retry
    </button>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  title?: string
  message: string
  onRetry?: () => void
}>(), {
  title: 'Error',
})
</script>
