
<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    <main class="flex-1 px-6 py-8 bg-gray-50 overflow-hidden">
      <div v-if="loading" class="text-center py-4">Loading group details...</div>
      <div v-else-if="error" class="text-red-500 text-center py-4">
        {{ error }}
      </div>
      <div v-else class="space-y-6">
        <!-- Group Details -->
        <div class="bg-white rounded-xl shadow-md p-6">
          <h1 class="text-2xl font-bold mb-4">{{ group?.name }}</h1>
          <p class="text-gray-600 mb-4">{{ group?.description }}</p>
          <p class="text-sm text-gray-500">
            Created: {{ formatDate(group?.created_at || '') }}
          </p>
        </div>

        <!-- Words in Group -->
        <div class="bg-white rounded-xl shadow-md overflow-hidden">
          <div class="p-4 border-b flex justify-between items-center">
            <h2 class="text-lg font-bold text-gray-800">Words in this Group</h2>
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Search words..."
              class="px-4 py-2 border rounded-lg text-sm w-64 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
          
          <div v-if="paginatedWords.length === 0" class="p-4 text-center text-gray-500">
            No words in this group
          </div>
          
          <table v-else class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">#</th>
                <th 
                  @click="sortBy('hindi')"
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors"
                >
                  Hindi
                  <span v-if="sortColumn === 'hindi'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th 
                  @click="sortBy('english')"
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors"
                >
                  English
                  <span v-if="sortColumn === 'english'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th 
                  @click="sortBy('hinglish')"
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors"
                >
                  Hinglish
                  <span v-if="sortColumn === 'hinglish'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th 
                  @click="sortBy('created_at')"
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors w-32"
                >
                  Created
                  <span v-if="sortColumn === 'created_at'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr 
                v-for="(word, index) in paginatedWords" 
                :key="word.id"
                class="hover:bg-gray-50 transition-colors"
              >
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-16">
                  {{ (currentPage - 1) * pageSize + index + 1 }}
                </td>
                <td class="px-4 py-4 whitespace-nowrap text-sm">{{ word.hindi }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm">{{ word.english }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm">{{ word.hinglish }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-32">
                  {{ formatDate(word.created_at) }}
                </td>
              </tr>
            </tbody>
          </table>

          <!-- Pagination -->
          <div class="px-4 py-3 bg-gray-50 border-t flex justify-between items-center">
            <span class="text-sm text-gray-600">
              Showing {{ (currentPage - 1) * pageSize + 1 }} to 
              {{ Math.min(currentPage * pageSize, filteredWords.length) }} of 
              {{ filteredWords.length }} words
            </span>
            <div class="flex space-x-2">
              <button 
                @click="prevPage" 
                :disabled="currentPage === 1"
                class="px-4 py-2 border rounded-lg text-sm bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                Previous
              </button>
              <button 
                @click="nextPage" 
                :disabled="currentPage === totalPages"
                class="px-4 py-2 border rounded-lg text-sm bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from '~/components/Sidebar.vue'
import { useGroupService, type Group, type GroupWord } from '~/services/groupService'

const route = useRoute()
const { 
  group, 
  groupWords, 
  loading, 
  error, 
  getGroupById, 
  getGroupWords 
} = useGroupService()

const groupId = Number(route.params.id)

// Pagination and Sorting
const pageSize = 10
const currentPage = ref(1)
const searchQuery = ref('')
const sortColumn = ref<keyof GroupWord>('created_at')
const sortDirection = ref<'asc' | 'desc'>('desc')

// Helper function for consistent date formatting
function formatDate(dateString: string): string {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return 'Invalid Date'
  }
}

// Computed properties for filtering and sorting
const filteredWords = computed(() => {
  return groupWords.value.filter(word => 
    !searchQuery.value || 
    (['hindi', 'hinglish', 'english'] as (keyof GroupWord)[]).some(field => {
      const value = word[field]
      return typeof value === 'string' && 
        value.toLowerCase().includes(searchQuery.value.toLowerCase())
    })
  ).sort((a, b) => {
    const modifier = sortDirection.value === 'asc' ? 1 : -1
    const aValue = a[sortColumn.value]
    const bValue = b[sortColumn.value]
    return aValue && bValue 
      ? aValue.toString().localeCompare(bValue.toString()) * modifier 
      : 0
  })
})

const totalPages = computed(() => Math.ceil(filteredWords.value.length / pageSize))

const paginatedWords = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return filteredWords.value.slice(start, end)
})

// Methods
function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

function sortBy(column: keyof GroupWord) {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = column
    sortDirection.value = 'asc'
  }
  currentPage.value = 1
}

// Fetch group details and words on component mount
onMounted(async () => {
  try {
    await Promise.all([
      getGroupById(groupId),
      getGroupWords(groupId)
    ])
    
    // Debug logging
    console.log('Group Details:', group.value)
    console.log('Group Description:', group.value?.description)
  } catch (err) {
    console.error('Failed to load group details:', err)
  }
})

// Watch for search query changes
watch(searchQuery, () => {
  currentPage.value = 1
})
</script>

<style scoped>
/* Add any additional styling if needed */
</style>