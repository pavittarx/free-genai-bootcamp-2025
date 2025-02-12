<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import Sidebar from '~/components/Sidebar.vue'
import { useWordService, type Word } from '~/services/wordService'

const { 
  words, 
  loading, 
  error, 
  fetchWords 
} = useWordService()

// Pagination
const pageSize = 20
const currentPage = ref(1)
const searchQuery = ref('')
const sortColumn = ref<keyof Word>('created_at')
const sortDirection = ref<'asc' | 'desc'>('desc')

// Error handling
const apiError = ref<string | null>(null)

// Fetch words on component mount
async function loadWords() {
  try {
    apiError.value = null
    await fetchWords({ 
      page: currentPage.value, 
      limit: pageSize, 
      search: searchQuery.value 
    })
  } catch (err) {
    apiError.value = err instanceof Error 
      ? err.message 
      : 'Failed to load words from the server'
    console.error('Words page load error:', err)
  }
}

onMounted(loadWords)

// Computed properties for filtering and sorting
const filteredWords = computed(() => {
  return words.value.filter(word => 
    !searchQuery.value || 
    (['hindi', 'hinglish', 'english'] as (keyof Word)[]).some(field => {
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

const totalWords = computed(() => filteredWords.value.length)
const totalPages = computed(() => Math.ceil(totalWords.value / pageSize))

const paginatedWords = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return filteredWords.value.slice(start, end)
})

// Methods
function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadWords()
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    loadWords()
  }
}

function sortBy(column: keyof Word) {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = column
    sortDirection.value = 'asc'
  }
}

// Watch for search query changes
watch(searchQuery, () => {
  currentPage.value = 1
  loadWords()
})
</script>

<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    <main class="flex-1 px-6 py-8 bg-gray-50 overflow-hidden">
      <!-- API Error Display -->
      <div 
        v-if="apiError" 
        class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4" 
        role="alert"
      >
        <strong class="font-bold">Error: </strong>
        <span class="block sm:inline">{{ apiError }}</span>
      </div>

      <!-- Words Table -->
      <div v-if="loading" class="text-center py-4">Loading words...</div>
      <div v-else-if="words.length === 0 && !apiError" class="text-center py-4">
        No words found.
      </div>
      <div v-else-if="!apiError" class="bg-white rounded-xl shadow-md overflow-hidden">
        <div class="p-4 flex justify-between items-center border-b">
          <h2 class="text-lg font-bold text-gray-800">Words</h2>
          <div class="flex space-x-3">
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Search words..."
              class="px-4 py-2 border rounded-lg text-sm w-64 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
          </div>
        </div>
      
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors w-16">#</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors">Hindi</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors">Hinglish</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors">English</th>
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
                <td class="px-4 py-4 whitespace-nowrap text-sm font-medium">{{ word.hindi }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm">{{ word.hinglish }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm">{{ word.english }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-32">
                  {{ new Date(word.created_at).toLocaleDateString('en-US', { 
                    year: 'numeric', 
                    month: 'short', 
                    day: 'numeric' 
                  }) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="px-4 py-3 bg-gray-50 border-t flex justify-between items-center">
          <span class="text-sm text-gray-600">
            Showing {{ (currentPage - 1) * pageSize + 1 }} to 
            {{ Math.min(currentPage * pageSize, totalWords) }} of 
            {{ totalWords }} words
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
    </main>
  </div>
</template>

<style scoped>
</style>
