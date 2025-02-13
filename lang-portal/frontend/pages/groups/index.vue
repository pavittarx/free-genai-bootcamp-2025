<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <Sidebar />
    <main class="flex-1 px-6 py-8 bg-gray-50 overflow-hidden">
      <!-- Groups Table -->
      <div v-if="loading" class="text-center py-4">Loading groups...</div>
      <div v-else-if="error" class="text-red-500 text-center py-4">
        {{ error }}
      </div>
      <div v-else class="bg-white rounded-xl shadow-md overflow-hidden">
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
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors w-16">#</th>
                <th 
                  @click="sortBy('name')"
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors"
                >
                  Group Name
                  <span v-if="sortColumn === 'name'" class="ml-1">
                    {{ sortDirection === 'asc' ? '↑' : '↓' }}
                  </span>
                </th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Description
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
                v-for="(group, index) in groups" 
                :key="group.id"
                class="hover:bg-gray-50 transition-colors cursor-pointer"
                @click="navigateToGroup(group.id)"
              >
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-16">
                  {{ (currentPage - 1) * pageSize + index + 1 }}
                </td>
                <td class="px-4 py-4 whitespace-nowrap text-sm font-medium">{{ group.name }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500">{{ group.description }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-32">
                  {{ formatDate(group.created_at) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div class="px-4 py-3 bg-gray-50 border-t flex justify-between items-center">
          <span class="text-sm text-gray-600">
            Showing {{ (currentPage - 1) * pageSize + 1 }} to 
            {{ Math.min(currentPage * pageSize, totalCount) }} of 
            {{ totalCount }} groups
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '~/components/Sidebar.vue'
import { useGroupService, type Group } from '~/services/groupService'

const router = useRouter()
const { groups, loading, error, fetchGroups } = useGroupService()

// Pagination
const pageSize = 20
const currentPage = ref(1)
const searchQuery = ref('')
const sortColumn = ref<keyof Group>('name')
const sortDirection = ref<'asc' | 'desc'>('asc')
const totalCount = ref(0)

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

// Fetch groups on component mount
onMounted(async () => {
  await fetchGroupsWithParams()
})

// Helper function to fetch groups with current parameters
async function fetchGroupsWithParams() {
  const result = await fetchGroups({ 
    page: currentPage.value, 
    pageSize: pageSize,
    search: searchQuery.value,
    sortBy: sortColumn.value,
    sortDirection: sortDirection.value
  })
  totalCount.value = result.totalCount
  
  // Debug logging
  console.log('Fetched groups:', groups.value)
  groups.value.forEach(group => {
    console.log('Group details:', {
      id: group.id,
      name: group.name,
      description: group.description,
      created_at: group.created_at
    })
  })
}

// Computed properties
const totalPages = computed(() => Math.ceil(totalCount.value / pageSize))

// Methods
function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    fetchGroupsWithParams()
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchGroupsWithParams()
  }
}

function sortBy(column: keyof Group) {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = column
    sortDirection.value = 'asc'
  }
  currentPage.value = 1
  fetchGroupsWithParams()
}

function navigateToGroup(groupId: number) {
  router.push(`/groups/${groupId}`)
}

// Watch for search query changes
watch(searchQuery, () => {
  currentPage.value = 1
  fetchGroupsWithParams()
})
</script>
