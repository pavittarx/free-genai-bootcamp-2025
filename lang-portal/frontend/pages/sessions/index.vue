<template>
  <PageLayout 
    :is-loading="isLoading" 
    :error="error"
    :on-retry="refetch"
  >
    <template #title>Learning Sessions</template>
    
    <template #actions>
      <button 
        class="bg-blue-500 text-white px-4 py-2 rounded-lg text-sm hover:bg-blue-600 transition"
        @click="startNewSession"
      >
        + Start Session
      </button>
    </template>

    <div v-if="sessions && sessions.length" class="space-y-4">
      <div 
        v-for="session in sessions" 
        :key="session.id"
        class="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-300 cursor-pointer"
      >
        <div class="p-6 flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
              <span class="text-xl">⏱️</span>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-800">{{ session.name }}</h3>
              <p class="text-sm text-gray-600">
                {{ formatDate(session.startTime) }} | 
                Duration: {{ formatDuration(session.duration) }}
              </p>
            </div>
          </div>
          <div class="flex items-center space-x-4">
            <span 
              class="text-sm px-3 py-1 rounded-full"
              :class="getStatusClass(session.status)"
            >
              {{ session.status }}
            </span>
            <NuxtLink 
              :to="`/sessions/${session.id}`"
              class="text-blue-500 hover:text-blue-700 text-sm font-medium"
            >
              View Details →
            </NuxtLink>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-center py-16 bg-white rounded-xl shadow-md">
      <div class="mb-4 text-6xl">🕰️</div>
      <h2 class="text-xl font-semibold text-gray-700 mb-2">No Learning Sessions</h2>
      <p class="text-gray-500">Start tracking your learning progress by beginning a new session!</p>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useRouter } from 'vue-router'
import PageLayout from '~/components/common/PageLayout.vue'

interface LearningSession {
  id: number
  name: string
  startTime: string
  duration: number
  status: 'Completed' | 'In Progress' | 'Not Started'
}

const router = useRouter()

const { 
  data: sessions, 
  isLoading, 
  error, 
  refetch 
} = useQuery({
  queryKey: ['sessions'],
  queryFn: async () => {
    const response = await fetch('/api/sessions')
    if (!response.ok) {
      throw new Error('Failed to fetch learning sessions')
    }
    return response.json()
  }
})

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatDuration = (minutes: number) => {
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60
  return hours > 0 
    ? `${hours}h ${remainingMinutes}m` 
    : `${minutes}m`
}

const getStatusClass = (status: string) => {
  switch(status) {
    case 'Completed': return 'bg-green-100 text-green-800'
    case 'In Progress': return 'bg-yellow-100 text-yellow-800'
    case 'Not Started': return 'bg-gray-100 text-gray-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const startNewSession = () => {
  router.push('/sessions/create')
}
</script>
