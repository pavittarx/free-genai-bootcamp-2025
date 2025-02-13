<template>
  <Teleport to="body">
    <div 
      v-if="isOpen" 
      class="fixed inset-0 z-[9999] flex items-center justify-center bg-transparent pointer-events-none"
    >
      <div 
        class="relative w-[90vw] max-w-[1200px] h-[80vh] bg-white rounded-2xl shadow-2xl border border-gray-200 overflow-hidden pointer-events-auto animate-popup-enter flex"
      >
        <button 
          @click="close"
          class="absolute top-4 right-4 z-10 text-gray-600 hover:text-gray-900 bg-white rounded-full p-2 shadow-md transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <div class="flex-grow overflow-y-auto p-6">
          <slot></slot>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, provide, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  initiallyOpen: {
    type: Boolean,
    default: false
  }
})

const route = useRoute()
const router = useRouter()

const isOpen = ref(props.initiallyOpen)

const open = () => {
  isOpen.value = true
}

const close = () => {
  isOpen.value = false
}

// Watch for route changes to close popup
watch(() => route.fullPath, () => {
  close()
})

// Provide popup methods to child components
provide('popupControl', {
  open,
  close
})

defineExpose({
  open,
  close
})
</script>

<style scoped>
@keyframes popupEnter {
  from {
    transform: scale(0.9);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.animate-popup-enter {
  animation: popupEnter 0.3s ease-out;
}
</style>
