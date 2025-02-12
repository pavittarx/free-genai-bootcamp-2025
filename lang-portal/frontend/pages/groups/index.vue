<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    <main class="flex-1 px-6 py-8 bg-gray-50 overflow-hidden">
      <!-- Groups Table -->
      <div class="bg-white rounded-xl shadow-md overflow-hidden">
        <div class="p-4 flex justify-between items-center border-b">
          <h2 class="text-lg font-bold text-gray-800">Groups</h2>
          <div class="flex space-x-3">
            <input 
              v-model="searchQuery"
              type="text" 
              placeholder="Search groups..."
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
                  @click="sortBy('name')"
                  class="table-header"
                >
                  Group Name
                  <span v-if="sortColumn === 'name'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th class="table-header">
                  Description
                </th>
                <th class="table-header w-24 text-center">
                  Word Count
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
                v-for="(group, index) in paginatedGroups" 
                :key="group.id"
                class="table-row cursor-pointer hover:bg-gray-50"
                @click="navigateToGroup(group.id)"
              >
                <td class="table-cell text-gray-500 w-16">
                  {{ (currentPage - 1) * pageSize + index + 1 }}
                </td>
                <td class="table-cell font-medium">{{ group.name }}</td>
                <td class="table-cell text-gray-500">{{ group.description }}</td>
                <td class="table-cell text-center font-medium">{{ group.wordCount }}</td>
                <td class="table-cell text-gray-500 w-32">
                  {{ formatDate(group.createdAt) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="px-4 py-3 bg-gray-50 border-t flex justify-between items-center">
          <span class="text-sm text-gray-600">
            Showing {{ (currentPage - 1) * pageSize + 1 }} to 
            {{ Math.min(currentPage * pageSize, totalGroups) }} of 
            {{ totalGroups }} groups
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
import { useRouter } from 'vue-router'
import Sidebar from '~/components/Sidebar.vue'

// Enhanced group interface with more details
interface Group {
  id: number
  name: string
  description: string
  wordCount: number
  createdAt: Date | string
}

// Mock data (to be replaced with actual API call)
const groups = ref<Group[]>([
  { 
    id: 1, 
    name: 'Basic Vocabulary', 
    description: 'Essential words for beginners',
    wordCount: 50,
    createdAt: new Date('2024-01-15') 
  },
  { 
    id: 2, 
    name: 'Conversational Hindi', 
    description: 'Phrases for everyday conversations',
    wordCount: 75,
    createdAt: new Date('2024-01-20') 
  },
  { 
    id: 3, 
    name: 'Advanced Phrases', 
    description: 'Complex expressions and idioms',
    wordCount: 30,
    createdAt: new Date('2024-02-01') 
  }
])

// Pagination
const pageSize = 20
const currentPage = ref(1)
const searchQuery = ref('')
const sortColumn = ref<keyof Group>('createdAt')
const sortDirection = ref<'asc' | 'desc'>('desc')

const router = useRouter()

// Computed properties for filtering, sorting, and pagination
const filteredGroups = computed(() => {
  return groups.value.filter(group => 
    !searchQuery.value || 
    ['name', 'description'].some(field => 
      group[field as keyof Group]
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

const totalGroups = computed(() => filteredGroups.value.length)
const totalPages = computed(() => Math.ceil(totalGroups.value / pageSize))

const paginatedGroups = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  const end = start + pageSize
  return filteredGroups.value.slice(start, end)
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

function sortBy(column: keyof Group): void {
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

function navigateToGroup(groupId: number): void {
  router.push(`/groups/${groupId}`)
}
</script>

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
