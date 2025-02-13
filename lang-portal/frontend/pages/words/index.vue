<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import Sidebar from '~/components/Sidebar.vue'
import { useWordService, type Word } from '~/services/wordService'

const { 
  words, 
  loading, 
  error, 
  fetchWords,
  searchWords
} = useWordService()

// Pagination
const pageSize = 10
const currentPage = ref(1)
const searchQuery = ref('')
const totalCount = ref(0)

// Error handling
const apiError = ref<string | null>(null)

// Fetch words on component mount
async function loadWords() {
  try {
    apiError.value = null
    const result = await fetchWords({ 
      page: currentPage.value, 
      pageSize: 10, 
    })
    totalCount.value = result.totalCount
  } catch (err) {
    apiError.value = err instanceof Error 
      ? err.message 
      : 'Failed to load words from the server'
    console.error('Words page load error:', err)
  }
}

// Search words
async function performSearch() {
  try {
    apiError.value = null
    const result = searchQuery.value 
      ? await searchWords(searchQuery.value)
      : await fetchWords({ 
          page: currentPage.value, 
          pageSize: 10 
        })
    
    words.value = result.words
    totalCount.value = result.totalCount
    currentPage.value = 1
  } catch (err) {
    apiError.value = err instanceof Error 
      ? err.message 
      : 'Failed to search words'
    console.error('Words search error:', err)
  }
}

onMounted(loadWords)

// Computed properties
const totalPages = computed(() => Math.ceil(totalCount.value / 10))

// Methods
async function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await loadWords()
  }
}

async function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    await loadWords()
  }
}

// Watch for search query changes
watch(searchQuery, () => {
  currentPage.value = 1
  performSearch()
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
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">#</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Hindi</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Hinglish</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">English</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-32">
                  Created
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr 
                v-for="(word, index) in words" 
                :key="word.id"
                class="hover:bg-gray-50 transition-colors"
              >
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-16">
                  {{ (currentPage - 1) * 10 + index + 1 }}
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
            Showing {{ (currentPage - 1) * 10 + 1 }} to 
            {{ Math.min(currentPage * 10, totalCount) }} of 
            {{ totalCount }} words
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
/* Add any additional styling if needed */
</style>
