<template>
  <PageLayout :is-loading="isLoading" :error="error" :on-retry="refetch">
    <template #title>Session Details</template>

    <div v-if="session" class="space-y-6">
      <!-- Session Overview Card -->
      <div class="bg-white rounded-lg shadow-md p-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h2 class="text-2xl font-bold text-gray-800">
              Session #{{ session.id }}
            </h2>
            <p class="text-gray-600">
              {{ session.activity_name }}
            </p>
          </div>
          <div class="flex items-center space-x-3">
            <span 
              class="px-3 py-1 rounded-full text-sm font-medium"
              :class="getStatusClass(session.status)"
            >
              {{ session.status }}
            </span>
            <span class="text-gray-500">
              {{ formatDate(session.start_time) }}
            </span>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-4 border-t pt-4">
          <div>
            <p class="text-xs text-gray-500 uppercase">Total Score</p>
            <p class="text-lg font-semibold text-gray-800">{{ session.score }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-500 uppercase">Duration</p>
            <p class="text-lg font-semibold text-gray-800">
              {{ formatDuration(session.duration || 0) }}
            </p>
          </div>
          <div>
            <p class="text-xs text-gray-500 uppercase">Completed Challenges</p>
            <p class="text-lg font-semibold text-gray-800">
              {{ sessionActivities.length }}
            </p>
          </div>
        </div>
      </div>

      <!-- Session Activities Table -->
      <div v-if="sessionActivities && sessionActivities.length > 0" class="bg-white rounded-lg shadow-md overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-50 border-b">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Challenge
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Input
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Answer
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Score
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Timestamp
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr 
              v-for="activity in sessionActivities" 
              :key="activity.id" 
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ activity.challenge }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ activity.input }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ activity.answer }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ activity.score }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(activity.created_at) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-else-if="sessionActivities" class="text-center py-8 text-gray-500">
        No activities in this session
      </div>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import PageLayout from '~/components/common/PageLayout.vue'

interface SessionActivity {
  id: number
  challenge: string
  input: string
  answer: string
  score: number
  created_at: string
}

interface Session {
  id: number
  activity_id: number
  activity_name: string
  start_time: string
  duration?: number
  score: number
  status: 'Completed' | 'In Progress' | 'Not Started'
  activities: SessionActivity[]
}

interface StudyActivity {
  id: number
  name: string
}

const route = useRoute()
const baseUrl = import.meta.env.VITE_BACKEND_BASE_URL || 'http://localhost:3000'
const sessionId = route.params.id as string

// Fetch session details
const { 
  data: session, 
  isLoading, 
  error,
  refetch 
} = useQuery({
  queryKey: ['session', sessionId],
  queryFn: async () => {
    const response = await fetch(`${baseUrl}/api/sessions/${sessionId}`)
    if (!response.ok) {
      throw new Error('Failed to fetch session details')
    }
    return response.json()
  }
})

const sessionActivities = computed(() => session.value?.activities || [])

// Fetch study activities to get activity name
const { 
  data: studyActivities, 
  isLoading: isStudyActivitiesLoading, 
  error: studyActivitiesError 
} = useQuery({
  queryKey: ['study-activities'],
  queryFn: async () => {
    const response = await fetch(`${baseUrl}/api/study-activities`)
    if (!response.ok) {
      throw new Error('Failed to fetch study activities')
    }
    return response.json()
  }
})

const activityName = computed(() => {
  if (!session.value || !studyActivities.value) return 'Unknown Activity'
  const activity = studyActivities.value.find(
    (a: StudyActivity) => a.id === session.value.activity_id
  )
  return activity ? activity.name : 'Unknown Activity'
})

const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A'
  
  // Try parsing with multiple strategies
  let date: Date | null = null
  
  // Try standard Date constructor
  date = new Date(dateString)
  
  // If standard parsing fails, try manual parsing
  if (isNaN(date.getTime())) {
    // Remove timezone if present and try parsing
    const cleanDateString = dateString.replace(/\.\d+Z?$/, '')
    date = new Date(cleanDateString)
  }
  
  // Check if date is still invalid
  if (isNaN(date.getTime())) {
    console.warn(`Invalid date: ${dateString}`)
    return dateString  // Return original string if parsing fails
  }
  
  return date.toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
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
  switch (status) {
    case 'Completed': return 'bg-green-100 text-green-800'
    case 'In Progress': return 'bg-yellow-100 text-yellow-800'
    case 'Not Started': return 'bg-gray-100 text-gray-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

console.log('Session:', session.value)
console.log('Session Activities:', sessionActivities.value)
console.log('Study Activities:', studyActivities.value)
</script>
