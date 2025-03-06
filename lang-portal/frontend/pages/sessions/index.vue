<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    <main class="flex-1 px-6 py-8 bg-gray-50 overflow-hidden">
      <div v-if="apiError" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" role="alert">
        <strong class="font-bold">Error: </strong>
        <span class="block sm:inline">{{ apiError }}</span>
      </div>

      <div v-if="loading" class="text-center py-4">Loading sessions...</div>
      <div v-else-if="paginatedSessions.length === 0 && !apiError" class="text-center py-4">
        No sessions found.
      </div>
      <div v-else-if="!apiError" class="bg-white rounded-xl shadow-md overflow-hidden">
        <div class="p-4 flex justify-between items-center border-b">
          <h2 class="text-lg font-bold text-gray-800">Sessions</h2>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Session ID</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Activity</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Start Time</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Score</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="(session) in paginatedSessions" :key="session.id" class="hover:bg-gray-50 transition-colors">
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500">{{ session.id }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm font-medium">{{ session.activityName || 'Unknown Activity' }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm">{{ formatDate(session.start_time) }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm">{{ session.score }}</td>
                <td class="px-4 py-4 whitespace-nowrap">
                  <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full" :class="getStatusClass(session.status)">{{ session.status }}</span>
                </td>
                <td class="px-4 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <NuxtLink :to="`/sessions/${session.id}`" class="text-blue-600 hover:text-blue-900">View Details</NuxtLink>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="px-4 py-3 bg-gray-50 border-t flex justify-between items-center">
          <span class="text-sm text-gray-600">
            Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, totalCount) }} of {{ totalCount }} sessions
          </span>
          <div class="flex space-x-2">
            <button @click="prevPage" :disabled="currentPage === 1" class="px-4 py-2 border rounded-lg text-sm bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">Previous</button>
            <button @click="nextPage" :disabled="currentPage >= totalPages" class="px-4 py-2 border rounded-lg text-sm bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors">Next</button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Sidebar from '~/components/Sidebar.vue'
import { useQuery } from '@tanstack/vue-query'

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

type SessionType = EnrichedSession

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

  return sessions.value.map((session: SessionType) => {
    const activity = studyActivities.value.find(
      (activity: StudyActivity) => activity.id === session.activity_id
    )
    return {
      ...session,
      activityName: activity ? activity.name : undefined
    }
  })
})

const loading = computed(() => isSessionsLoading.value || isActivitiesLoading.value)
const apiError = computed(() => sessionsError.value || activitiesError.value)

const itemsPerPage = 10
const currentPage = ref(1)

const totalCount = computed(() => enrichedSessions.value.length)
const totalPages = computed(() => Math.ceil(totalCount.value / itemsPerPage))

const paginatedSessions = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  return enrichedSessions.value.slice(start, start + itemsPerPage)
})

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
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
