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
                  @click="sortBy('group')"
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors"
                >
                  Group Name
                  <span v-if="sortColumn === 'group'" class="ml-1">
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
                v-for="(group, index) in paginatedGroups" 
                :key="group.id"
                class="hover:bg-gray-50 transition-colors cursor-pointer"
                @click="navigateToGroup(group.id)"
              >
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-16">
                  {{ (currentPage - 1) * pageSize + index + 1 }}
                </td>
                <td class="px-4 py-4 whitespace-nowrap text-sm font-medium">{{ group.group }}</td>
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-32">
                  {{ new Date(group.created_at).toLocaleDateString('en-US', { 
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
const sortColumn = ref<keyof Group>('group')
const sortDirection = ref<'asc' | 'desc'>('asc')

// Fetch groups on component mount
onMounted(async () => {
  await fetchGroups({ 
    page: currentPage.value, 
    limit: pageSize, 
    search: searchQuery.value 
  })
})

// Computed properties for filtering and sorting
const filteredGroups = computed(() => {
  return groups.value.filter(group => 
    !searchQuery.value || 
    group.group.toLowerCase().includes(searchQuery.value.toLowerCase())
  ).sort((a, b) => {
    const modifier = sortDirection.value === 'asc' ? 1 : -1
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
function nextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    fetchGroups({ 
      page: currentPage.value, 
      limit: pageSize, 
      search: searchQuery.value 
    })
  }
}

function prevPage() {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchGroups({ 
      page: currentPage.value, 
      limit: pageSize, 
      search: searchQuery.value 
    })
  }
}

function sortBy(column: keyof Group) {
  if (sortColumn.value === column) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortColumn.value = column
    sortDirection.value = 'asc'
  }
}

function navigateToGroup(groupId: number) {
  router.push(`/groups/${groupId}`)
}

// Watch for search query changes
watch(searchQuery, () => {
  currentPage.value = 1
  fetchGroups({ 
    page: currentPage.value, 
    limit: pageSize, 
    search: searchQuery.value 
  })
})
</script>
