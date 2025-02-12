
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from '~/components/Sidebar.vue'
import { 
  Book, 
  Calendar, 
  FileText 
} from 'lucide-vue-next'

// Enhanced group interface with more details
interface Group {
  id: number
  name: string
  description: string
  wordCount: number
  createdAt: Date | string
}

// Enhanced word interface
interface Word {
  id: number
  hindi: string
  hinglish: string
  english: string
  createdAt: Date | string
}

const route = useRoute()
const groupId = Number(route.params.id)

const group = ref<Group>({
  id: 1, 
  name: 'Basic Vocabulary', 
  description: 'Essential words for beginners',
  wordCount: 50,
  createdAt: new Date('2024-01-15')
})

const words = ref<Word[]>([
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
const sortColumn = ref<keyof Word>('createdAt')
const sortDirection = ref<'asc' | 'desc'>('desc')

// Computed properties for filtering, sorting, and pagination
const filteredWords = computed(() => {
  return words.value.filter(word => 
    !searchQuery.value || 
    ['hindi', 'hinglish', 'english'].some(field => 
      word[field as keyof Word]
        .toString()
        .toLowerCase()
        .includes(searchQuery.value.toLowerCase())
    )
  ).sort((a, b) => {
    const modifier = sortDirection.value === 'asc' ? 1 : -1
    
    const createdAtA = a.createdAt instanceof Date 
      ? a.createdAt 
      : new Date(a.createdAt)
    const createdAtB = b.createdAt instanceof Date 
      ? b.createdAt 
      : new Date(b.createdAt)
    
    if (sortColumn.value === 'createdAt') {
      return (createdAtA.getTime() > createdAtB.getTime() ? 1 : -1) * modifier
    }
    
    return a[sortColumn.value].toString().localeCompare(
      b[sortColumn.value].toString()
    ) * modifier
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
function nextPage(): void {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

function prevPage(): void {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

function sortBy(column: keyof Word): void {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = column
    sortDirection.value = 'asc'
  }
}

function formatDate(date: Date | string): string {
  const dateObj = date instanceof Date 
    ? date 
    : new Date(date)
  
  return dateObj.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  })
}

// Ensure date conversion on mounted
onMounted(() => {
  // Convert createdAt to Date if it's a string
  if (!(group.value.createdAt instanceof Date)) {
    group.value.createdAt = new Date(group.value.createdAt)
  }

  words.value = words.value.map(word => ({
    ...word,
    createdAt: word.createdAt instanceof Date 
      ? word.createdAt 
      : new Date(word.createdAt)
  }))
})
</script>

<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    <main class="flex-1 px-6 py-8 bg-gray-50 overflow-hidden">
      <!-- Group Detail Section -->
      <div class="bg-white rounded-xl shadow-md p-6 mb-6">
        <div class="flex items-center space-x-6">
          <div class="bg-blue-100 rounded-full w-16 h-16 flex items-center justify-center">
            <FileText 
              v-if="group"
              class="h-8 w-8 text-blue-600" 
            />
          </div>
          <div class="flex-1">
            <h2 class="text-2xl font-bold text-gray-800 mb-1">{{ group.name }}</h2>
            <p class="text-gray-600 mb-2">{{ group.description }}</p>
            <div class="flex items-center space-x-4 text-sm text-gray-500">
              <div class="flex items-center space-x-1">
                <Book 
                  v-if="group"
                  class="h-4 w-4 text-gray-400"
                />
                <span>{{ group.wordCount }} Words</span>
              </div>
              <div class="flex items-center space-x-1">
                <Calendar 
                  v-if="group"
                  class="h-4 w-4 text-gray-400"
                />
                <span>Created {{ formatDate(group.createdAt) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Words Table -->
      <div class="bg-white rounded-xl shadow-md overflow-hidden">
        <div class="p-4 flex justify-between items-center border-b">
          <h3 class="text-lg font-bold text-gray-800">Group Words</h3>
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
                <th class="table-header w-16">#</th>
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

<style scoped>
.table-header {
  @apply px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors;
}

.table-cell {
  @apply px-4 py-4 whitespace-nowrap text-sm;
}

.table-row {
  @apply hover:bg-gray-50 transition-colors;
}
</style>