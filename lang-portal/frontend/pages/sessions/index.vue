<template>
  <PageLayout :is-loading="isLoading" :error="error" :on-retry="refetch">
    <template #title>Learning Sessions</template>

    <div class="bg-white rounded-lg shadow-md overflow-hidden">
      <table class="w-full">
        <thead class="bg-gray-50 border-b">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Session ID
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Activity
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Start Time
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Score
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200">
          <tr 
            v-for="session in enrichedSessions" 
            :key="session.id" 
            class="hover:bg-gray-50 transition-colors"
          >
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ session.id }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ session.activityName || 'Unknown Activity' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(session.start_time) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ session.score }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span 
                class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                :class="getStatusClass(session.status)"
              >
                {{ session.status }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <NuxtLink 
                :to="`/sessions/${session.id}`" 
                class="text-blue-600 hover:text-blue-900"
              >
                View Details
              </NuxtLink>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="!enrichedSessions || enrichedSessions.length === 0" class="text-center py-12">
        <div class="text-6xl mb-4">🕰️</div>
        <h2 class="text-xl font-semibold text-gray-700 mb-2">No Learning Sessions</h2>
        <p class="text-gray-500">Start tracking your learning progress by beginning a new session!</p>
      </div>
    </div>
  </PageLayout>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import PageLayout from '~/components/common/PageLayout.vue'

interface StudyActivity {
  id: number
  name: string
}

interface Session {
  id: number
  activity_id: number
  start_time: string
  score: number
  status: 'Completed' | 'In Progress' | 'Not Started'
}

interface EnrichedSession extends Session {
  activityName?: string
}

const baseUrl = import.meta.env.VITE_BACKEND_BASE_URL || 'http://localhost:3000'

const { 
  data: sessions, 
  isLoading: isSessionsLoading, 
  error: sessionsError 
} = useQuery({
  queryKey: ['sessions'],
  queryFn: async () => {
    const response = await fetch(`${baseUrl}/api/sessions`)
    if (!response.ok) {
      throw new Error('Failed to fetch learning sessions')
    }
    return response.json()
  }
})

const { 
  data: studyActivities, 
  isLoading: isActivitiesLoading, 
  error: activitiesError 
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

const enrichedSessions = computed<EnrichedSession[]>(() => {
  if (!sessions.value || !studyActivities.value) return []

  return sessions.value.map(session => {
    const activity = studyActivities.value.find(
      (activity: StudyActivity) => activity.id === session.activity_id
    )
    return {
      ...session,
      activityName: activity ? activity.name : undefined
    }
  })
})

const isLoading = computed(() => isSessionsLoading.value || isActivitiesLoading.value)
const error = computed(() => sessionsError.value || activitiesError.value)

const refetch = () => {
  sessions.value
  studyActivities.value
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusClass = (status: string) => {
  switch (status) {
    case 'Completed': return 'bg-green-100 text-green-800'
    case 'In Progress': return 'bg-yellow-100 text-yellow-800'
    case 'Not Started': return 'bg-gray-100 text-gray-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}
</script>
