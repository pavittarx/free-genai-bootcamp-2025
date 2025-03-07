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
        </div>

        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">#</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Group Name</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-32">Created</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="(group, index) in groups" :key="group.id"
                class="hover:bg-gray-50 transition-colors cursor-pointer" @click="navigateToGroup(group.id)">
                <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-500 w-16">
                  {{ index + 1 }}
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
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import Sidebar from '~/components/Sidebar.vue'
import { useGroupService } from '~/services/groupService'

const router = useRouter()
const { groups, loading, error, fetchGroups } = useGroupService()

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
  try {
    await fetchGroups()
  } catch (err) {
    console.error('Error fetching groups:', err)
  }
})

function navigateToGroup(groupId: number) {
  router.push(`/groups/${groupId}`)
}
</script>
