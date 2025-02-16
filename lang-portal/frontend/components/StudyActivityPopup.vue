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
          @click="endSession"
          class="absolute top-4 right-4 z-10 text-gray-600 hover:text-gray-900 bg-white rounded-full p-2 shadow-md transition"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <div v-if="!isFinished" class="flex-grow overflow-y-auto p-6">
          <slot></slot>
        </div>

        <div v-else class="flex-grow flex items-center justify-center flex-col p-6">
          <h2 class="text-3xl font-bold mb-4">Session Completed</h2>
          <p class="text-xl mb-6">Your final score: {{ sessionScore }}</p>
          <button 
            @click="close" 
            class="px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, provide, watch, computed, defineEmits, defineExpose } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  initiallyOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'finish'])

const route = useRoute()
const router = useRouter()

const isOpen = ref(props.initiallyOpen)
const isFinished = ref(false)
const sessionScore = ref(0)

const open = () => {
  isOpen.value = true
  isFinished.value = false
}

const close = () => {
  isOpen.value = false
}

const finish = () => {
  isFinished.value = true
  sessionScore.value = calculateFinalScore()
  emit('finish', sessionScore.value)
}

const endSession = () => {
  finish()
  close()
}

const calculateFinalScore = () => {
  // Placeholder for final score calculation
  // This should be replaced with actual scoring logic
  return Math.floor(Math.random() * 100)
}

// Watch for route changes to close popup
watch(() => route.fullPath, () => {
  close()
})

// Provide popup methods to child components
provide('popupControl', {
  open,
  close,
  finish,
  endSession
})

defineExpose({
  open,
  close,
  finish,
  endSession
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
