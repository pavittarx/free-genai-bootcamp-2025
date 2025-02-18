<template>
  <PageLayout 
    :is-loading="isLoading" 
    :error="error"
    :on-retry="refetch"
  >
    <template #title>
      <div class="p-4 flex justify-between items-center border-b">
        <h2 class="text-lg font-bold text-gray-800">{{ activity?.name }}</h2>
      </div>
    </template>

    <div v-if="activity" class="p-6 bg-white rounded-xl shadow-md">
      <div class="flex items-center space-x-6">
        <div class="w-1/3">
          <img 
            :src="getActivityImage(activity)" 
            :alt="activity.name" 
            class="w-full h-auto rounded-lg shadow-md"
          />
        </div>
        <div class="w-2/3">
          <h3 class="text-2xl font-bold mb-4">{{ activity.name }}</h3>
          <p class="text-gray-600 mb-6">{{ activity.description }}</p>
          
          <div class="grid grid-cols-2 gap-4 mb-6">
            <div>
              <span class="text-sm text-gray-500">Created At</span>
              <p class="font-semibold">{{ formatDate(activity.created_at) }}</p>
            </div>
            <div>
              <span class="text-sm text-gray-500">Current Score</span>
              <p class="font-semibold">{{ formatScore(activity.score) }}</p>
            </div>
          </div>
          
          <button 
            @click="startActivity" 
            class="w-full bg-blue-500 text-white py-3 rounded-lg hover:bg-blue-600 transition"
          >
            Start Activity
          </button>
        </div>
      </div>
    </div>

    <div v-else-if="isLoading" class="text-center py-16">
      <div class="mb-4 text-6xl animate-pulse">🏁</div>
      <h2 class="text-xl font-semibold text-gray-700 mb-2">Loading Activity</h2>
      <p class="text-gray-500">Fetching activity details...</p>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { useQuery } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'

interface StudyActivity {
  id: string
  name: string
  description: string
  score: number
  created_at: string
}

const route = useRoute()
const router = useRouter()

const ACTIVITY_IMAGES: Record<string, string> = {
  'Unscramble Words': '/game-1.jpg',
  'Word Match': '/game-2.jpg',
  'Sentence Constructor': '/game-3.jpg',
  'Default': '/default-activity.jpg'
}

const { 
  data: activity, 
  isLoading, 
  error, 
  refetch 
} = useQuery({
  queryKey: ['studyActivity', route.params.id],
  queryFn: async () => {
    try {
      const response = await fetch(`/api/study-activities/${route.params.id}`)
      if (!response.ok) {
        throw new Error('Failed to fetch activity')
      }
      return response.json()
    } catch (err) {
      console.error('Error fetching activity:', err)
      throw err
    }
  },
  retry: 2,
  retryDelay: 1000
})

const getActivityImage = (activity?: StudyActivity): string => {
  return activity?.name ? ACTIVITY_IMAGES[activity.name] || ACTIVITY_IMAGES['Default'] : ACTIVITY_IMAGES['Default']
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const formatScore = (score: number): string => {
  return score.toFixed(2)
}

const startActivity = () => {
  if (activity.value) {
    router.push(`/study-activities/${activity.value.id}/launch`)
  }
}
</script>

<style scoped>
/* Inherit styles from the study-activities index page */
</style>
