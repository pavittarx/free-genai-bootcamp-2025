<template>
  <div class="flex justify-center items-center" :class="sizeClasses">
    <div 
      class="animate-spin rounded-full border-t-2 border-blue-500" 
      :class="[`h-${size} w-${size}`, borderWidth]"
    ></div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  size?: number
  fullScreen?: boolean
}>(), {
  size: 12,
  fullScreen: false
})

const sizeClasses = computed(() => ({
  'h-screen w-screen': fullScreen,
}))

const borderWidth = computed(() => {
  switch (size) {
    case 8: return 'border-2'
    case 12: return 'border-4'
    case 16: return 'border-6'
    default: return 'border-4'
  }
})
</script>
