<template>
  <PageLayout 
    :is-loading="isLoading" 
    :error="error"
    :on-retry="refetch"
  >
    <template #title>
      <div class="p-4 flex justify-between items-center border-b">
        <h2 class="text-lg font-bold text-gray-800">Study Activities</h2>
      </div>
    </template>

    <StudyActivityPopup ref="activityPopup">
      <component 
        v-if="selectedActivity"
        :is="getActivityComponent(selectedActivity.name)"
        :activity-id="selectedActivity.id"
      />
    </StudyActivityPopup>

    <div 
      v-if="studyActivities && studyActivities.length" 
      class="w-full flex justify-center items-center gap-8 p-8"
    >
      <div 
        v-for="activity in studyActivities" 
        :key="activity.id"
        @click="launchActivity(activity)"
        class="pokemon-card group mx-4 cursor-pointer"
      >
        <div class="pokemon-card-inner">
          <div class="pokemon-card-front">
            <div class="pokemon-card-image">
              <img 
                :src="getActivityImage(activity)" 
                :alt="activity.name" 
                class="w-full h-full object-cover"
              />
            </div>
            <div class="pokemon-card-content">
              <h3 class="pokemon-card-title">{{ activity.name }}</h3>
              <p class="pokemon-card-description">{{ activity.description }}</p>
              <div class="pokemon-card-stats">
                <span class="pokemon-card-score">Score: {{ formatScore(activity.score) }}</span>
                <span class="pokemon-card-date">{{ formatDate(activity.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="isLoading" class="text-center py-16 bg-white rounded-xl shadow-md">
      <div class="mb-4 text-6xl">🏁</div>
      <h2 class="text-xl font-semibold text-gray-700 mb-2">No Study Activities</h2>
      <p class="text-gray-500">Study activities will be added soon!</p>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import PageLayout from '~/components/common/PageLayout.vue'
import StudyActivityPopup from '~/components/StudyActivityPopup.vue'
import UnscrambleWordsActivity from '~/components/UnscrambleWordsActivity.vue'

interface StudyActivity {
  id: string
  name: string
  description: string
  image: string
  score: number
  created_at: string
}

const activityPopup = ref(null)
const selectedActivity = ref<StudyActivity | null>(null)

const ACTIVITY_IMAGES: Record<string, string> = {
  'Unscramble Words': '/game-1.jpg',
  'Group Words': '/game-2.jpg',
  'Complete the Word': '/game-3.jpg'
}

const { 
  data: studyActivities, 
  isLoading, 
  error, 
  refetch 
} = useQuery<StudyActivity[]>({
  queryKey: ['studyActivities'],
  queryFn: async () => {
    try {
      const response = await fetch('http://localhost:3000/api/study-activities', {
        method: 'GET',
        headers: {
          'Accept': 'application/json'
        }
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`HTTP error! status: ${response.status}, message: ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      console.error('Failed to fetch study activities:', err)
      throw err instanceof Error 
        ? err 
        : new Error('An unknown error occurred')
    }
  },
  retry: 2,
  retryDelay: 1000
})

const getActivityImage = (activity: StudyActivity): string => {
  return ACTIVITY_IMAGES[activity.name] || '/game-default.jpg'
}

const getActivityComponent = (activityName: string) => {
  const componentMap = {
    'Unscramble Words': UnscrambleWordsActivity,
    // Add other activity components as they are created
  } as const

  return componentMap[activityName as keyof typeof componentMap] || null
}

const launchActivity = (activity: StudyActivity) => {
  selectedActivity.value = activity
  if (activityPopup.value && 'open' in activityPopup.value) {
    (activityPopup.value as { open: () => void }).open()
  }
}

const formatDate = (dateString: string): string => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const formatScore = (score: number): string => {
  return score.toString().padStart(2, '0')
}
</script>

<style scoped>
.pokemon-card {
  position: relative;
  width: 16rem;
  height: 24rem;
  background: linear-gradient(to bottom right, #dbeafe, #93c5fd);
  border-radius: 1rem;
  overflow: hidden;
  transition: all 0.3s ease;
  cursor: pointer;
  margin: 1rem;
}

.pokemon-card:hover {
  transform: scale(1.05);
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.pokemon-card:active {
  transform: scale(0.95);
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
}

.pokemon-card-inner {
  width: 100%;
  height: 100%;
  position: relative;
}

.pokemon-card-front {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.pokemon-card-image {
  width: 100%;
  height: 50%;
  overflow: hidden;
}

.pokemon-card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.pokemon-card-content {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex-grow: 1;
}

.pokemon-card-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 0.5rem;
}

.pokemon-card-description {
  font-size: 0.875rem;
  color: #4b5563;
  margin-bottom: 1rem;
  flex-grow: 1;
}

.pokemon-card-stats {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
  color: #6b7280;
}

.pokemon-card-score {
  background-color: #3b82f6;
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
}
</style>
