<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    <main class="flex-1 px-6 py-8 bg-gray-50 overflow-hidden">
      <!-- Words Table -->
      <div class="bg-white rounded-xl shadow-md overflow-hidden">
        <div class="p-4 flex justify-between items-center border-b">
          <h2 class="text-lg font-bold text-gray-800">Words</h2>
          <div class="flex space-x-3">
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Search words..."
              class="px-4 py-2 border rounded-lg text-sm w-64 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
            <select 
              v-model="languageFilter"
              class="px-4 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
            >
              <option value="">All Languages</option>
              <option value="hindi">Hindi</option>
              <option value="hinglish">Hinglish</option>
              <option value="english">English</option>
            </select>
          </div>
        </div>
      
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="table-header w-16">
                  #
                </th>
                <th 
                  @click="sortBy('hindi')"
                  class="table-header"
                >
                  Hindi
                  <span v-if="sortColumn === 'hindi'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th 
                  @click="sortBy('hinglish')"
                  class="table-header"
                >
                  Hinglish
                  <span v-if="sortColumn === 'hinglish'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th 
                  @click="sortBy('english')"
                  class="table-header"
                >
                  English
                  <span v-if="sortColumn === 'english'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th 
                  @click="sortBy('createdAt')"
                  class="table-header w-32"
                >
                  Created
                  <span v-if="sortColumn === 'createdAt'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr 
                v-for="(word, index) in paginatedWords" 
                :key="word.id"
                class="table-row"
              >
                <td class="table-cell text-gray-500 w-16">
                  {{ (currentPage - 1) * pageSize + index + 1 }}
                </td>
                <td class="table-cell font-medium">{{ word.hindi }}</td>
                <td class="table-cell">{{ word.hinglish }}</td>
                <td class="table-cell">{{ word.english }}</td>
                <td class="table-cell text-gray-500 w-32">
                  {{ formatDate(word.createdAt) }}
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

<script setup lang="ts">
import { ref, computed } from 'vue'
import Sidebar from '~/components/Sidebar.vue'

// Mock data (to be replaced with actual API call)
const words = ref([
  { 
    id: 1, 
    hindi: 'नमस्ते', 
    hinglish: 'Namaste', 
    english: 'Hello', 
    createdAt: new Date('2024-01-15') 
  },
  { 
    id: 2, 
    hindi: 'धन्यवाद', 
    hinglish: 'Dhanyavaad', 
    english: 'Thank you', 
    createdAt: new Date('2024-01-20') 
  }
])

// Pagination
const pageSize = 20
const currentPage = ref(1)
const searchQuery = ref('')
const languageFilter = ref('')
const sortColumn = ref('createdAt')
const sortDirection = ref('desc')

// Computed properties for filtering, sorting, and pagination
const filteredWords = computed(() => {
  return words.value.filter(word => {
    const matchesSearch = !searchQuery.value || 
      word.hindi.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      word.hinglish.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      word.english.toLowerCase().includes(searchQuery.value.toLowerCase())
    
    const matchesLanguageFilter = !languageFilter.value || 
      (languageFilter.value === 'hindi' && word.hindi) ||
      (languageFilter.value === 'hinglish' && word.hinglish) ||
      (languageFilter.value === 'english' && word.english)
    
    return matchesSearch && matchesLanguageFilter
  }).sort((a, b) => {
    const modifier = sortDirection.value === 'asc' ? 1 : -1
    if (sortColumn.value === 'createdAt') {
      return a[sortColumn.value].getTime() > b[sortColumn.value].getTime()
        ? modifier 
        : -modifier
    }
    return a[sortColumn.value].localeCompare(b[sortColumn.value]) * modifier
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
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

function sortBy(column: string) {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = column
    sortDirection.value = 'asc'
  }
}

function formatDate(date: Date) {
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  }).format(date)
}
</script>

<style scoped>
.table-header {
  @apply px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors;
}

.table-cell {
  @apply px-4 py-3 text-sm whitespace-nowrap;
}

.table-row {
  @apply hover:bg-gray-50 transition-colors;
}
</style>
