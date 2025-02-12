<template>
  <aside class="w-64 bg-gradient-to-b from-blue-50 to-blue-100 px-6 py-8 flex flex-col">
    <div class="text-left mb-6">
      <h2 class="text-xl font-bold text-gray-800">Language Portal</h2>
      <p class="text-xs text-gray-600 mt-1">Learning Tracker</p>
    </div>

    <nav class="flex-grow overflow-auto">
      <ul class="space-y-2">
        <li v-for="link in navigationLinks" :key="link.to">
          <NuxtLink 
            :to="link.to"
            class="flex items-center p-3 rounded-xl transition-all duration-300 group text-gray-700 hover:bg-blue-200 hover:text-blue-600"
            :class="{ 'bg-blue-200 text-blue-600': isCurrentRoute(link.to) }"
          >
            <span class="text-lg mr-3">{{ link.icon }}</span>
            <span class="text-sm flex-grow">{{ link.label }}</span>
          </NuxtLink>
        </li>
      </ul>
    </nav>

    <div class="mt-4 pt-4 border-t border-gray-200 relative">
      <div 
        @click="toggleUserPopup"
        class="flex items-center cursor-pointer hover:bg-blue-100 p-2 rounded-xl transition-colors"
      >
        <div class="w-10 h-10 bg-blue-400 text-white rounded-full flex items-center justify-center mr-3">
          {{ getUserInitial }}
        </div>
        <div>
          <p class="text-sm font-semibold">{{ userName || 'Set Name' }}</p>
          <p class="text-xs text-gray-500">{{ userName ? 'Edit Profile' : 'Click to set name' }}</p>
        </div>
      </div>

      <!-- User Popup -->
      <div 
        v-if="showUserPopup"
        class="absolute bottom-full left-0 right-0 mb-2 bg-white rounded-xl shadow-lg border border-gray-200 p-4 z-50"
      >
        <div class="mb-4">
          <label class="block text-xs text-gray-600 mb-1">Your Name</label>
          <input 
            v-model="editedUserName"
            type="text" 
            placeholder="Enter your name"
            class="w-full px-3 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
        </div>
        <div class="space-y-2">
          <button 
            @click="saveUserName"
            class="w-full bg-blue-500 text-white py-2 rounded-lg text-sm hover:bg-blue-600 transition-colors"
          >
            Save Name
          </button>
          <button 
            @click="showResetConfirmation = true; toggleUserPopup()"
            class="w-full bg-red-500 text-white py-2 rounded-lg text-sm hover:bg-red-600 transition-colors"
          >
            Reset / Clear
          </button>
          <button 
            @click="toggleUserPopup"
            class="w-full bg-gray-200 text-gray-700 py-2 rounded-lg text-sm hover:bg-gray-300 transition-colors"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const showUserPopup = ref(false)
const editedUserName = ref('')
const userName = ref('')
const showResetConfirmation = ref(false)

const navigationLinks = [
  { to: '/', icon: '📊', label: 'Dashboard' },
  { to: '/words', icon: '📚', label: 'Words' },
  { to: '/practice', icon: '✏️', label: 'Practice' },
  { to: '/stats', icon: '📈', label: 'Statistics' }
]

const getUserInitial = computed(() => {
  return userName.value ? userName.value.charAt(0).toUpperCase() : '?'
})

function isCurrentRoute(path: string) {
  return route.path === path
}

function toggleUserPopup() {
  showUserPopup.value = !showUserPopup.value
  if (showUserPopup.value) {
    editedUserName.value = userName.value
  }
}

function saveUserName() {
  userName.value = editedUserName.value
  showUserPopup.value = false
}

function resetAllActivity() {
  userName.value = ''
  showResetConfirmation.value = false
  // Add any other reset logic here
}
</script>
