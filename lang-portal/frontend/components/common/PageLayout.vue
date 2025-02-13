<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    
    <main class="flex-grow bg-gray-50 p-8 overflow-auto">
      <div class="container mx-auto">
        <header class="mb-6 flex justify-between items-center">
          <h1 class="text-3xl font-bold text-gray-800">
            <slot name="title">Page Title</slot>
          </h1>
          <div class="flex items-center space-x-4">
            <slot name="actions"></slot>
          </div>
        </header>

        <div class="space-y-6">
          <LoadingSpinner v-if="isLoading" :full-screen="false" />
          
          <ErrorMessage 
            v-else-if="error" 
            :message="error.message" 
            :on-retry="onRetry"
          />
          
          <slot v-else></slot>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { PropType } from 'vue'

defineProps({
  isLoading: {
    type: Boolean,
    default: false
  },
  error: {
    type: Object as PropType<Error | null>,
    default: null
  },
  onRetry: {
    type: Function,
    default: () => {}
  }
})
</script>
