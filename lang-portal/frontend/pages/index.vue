<template>
  <div class="h-screen w-screen overflow-hidden flex">
    <!-- Sidebar Navigation -->
    <aside class="w-64 bg-gradient-to-b from-blue-50 to-blue-100 px-6 py-8 flex flex-col">
      <div class="text-left mb-6">
        <h2 class="text-xl font-bold text-gray-800">Language Portal</h2>
        <p class="text-xs text-gray-600 mt-1">Learning Tracker</p>
      </div>

      <nav class="flex-grow overflow-auto">
        <ul class="space-y-2">
          <li v-for="(link, index) in navLinks" :key="index">
            <NuxtLink 
              :to="link.to"
              class="flex items-center p-3 rounded-xl transition-all duration-300 group text-gray-700 hover:bg-blue-200 hover:text-blue-600"
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

    <!-- Main Content -->
    <main class="flex-1 grid grid-cols-3 grid-rows-2 gap-4 px-6 py-8 bg-gray-50 overflow-hidden">
      <!-- Dashboard Metrics -->
      <div class="col-span-2 grid grid-cols-4 gap-3">
        <div 
          v-for="metric in dashboardMetrics" 
          :key="metric.label"
          class="bg-white rounded-xl shadow-md p-4 flex flex-col justify-center items-center h-40"
        >
          <h3 class="text-[11px] text-gray-600 mb-3 uppercase tracking-[0.25em] font-bold opacity-80">
            {{ metric.label }}
          </h3>
          <p 
            class="text-4xl font-extrabold" 
            :class="metric.color"
          >
            {{ metric.value }}
          </p>
        </div>
      </div>

      <!-- Study Activities -->
      <div class="col-start-3 row-span-1 bg-white rounded-xl shadow p-4 flex flex-col overflow-hidden">
        <div class="flex justify-between items-center mb-3">
          <h2 class="text-base font-bold text-gray-800">Study Activities</h2>
          <NuxtLink 
            to="/study-activities" 
            class="text-xs text-blue-600 hover:text-blue-800"
          >
            View All
          </NuxtLink>
        </div>
        <div class="space-y-3 overflow-auto flex-grow">
          <div 
            v-for="activity in studyActivities" 
            :key="activity.id" 
            class="bg-gray-100 rounded-lg p-3 flex items-center justify-between hover:bg-blue-50 transition-colors"
          >
            <div class="flex items-center">
              <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center mr-3">
                <span class="text-base">{{ activity.icon }}</span>
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-800">{{ activity.name }}</h3>
                <p class="text-xs text-gray-500">{{ activity.description }}</p>
              </div>
            </div>
            <NuxtLink 
              :to="`/study-activities/${activity.id}`"
              class="bg-blue-500 text-white text-xs px-2 py-1 rounded hover:bg-blue-600"
            >
              Start
            </NuxtLink>
          </div>
        </div>
      </div>

      <!-- Recent Sessions -->
      <div class="col-start-3 row-start-2 bg-white rounded-xl shadow p-4 flex flex-col overflow-hidden">
        <div class="flex justify-between items-center mb-3">
          <h2 class="text-base font-bold text-gray-800">Recent Sessions</h2>
          <NuxtLink 
            to="/sessions" 
            class="text-xs text-blue-600 hover:text-blue-800"
          >
            View All
          </NuxtLink>
        </div>
        <div class="space-y-3 overflow-auto flex-grow">
          <div 
            v-for="session in recentSessions" 
            :key="session.id" 
            class="bg-gray-100 rounded-lg p-3 flex items-center justify-between hover:bg-blue-50 transition-colors"
          >
            <div>
              <h3 class="text-sm font-semibold text-gray-800">{{ session.activityName }}</h3>
              <p class="text-xs text-gray-500">{{ session.groupName }}</p>
            </div>
            <div class="text-right">
              <span 
                :class="getScoreClass(session.score)"
                class="font-bold text-xs"
              >
                {{ session.score }}%
              </span>
              <div class="text-xs text-gray-500 mt-1">
                {{ formatSessionDate(session.endTime) }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <!-- Reset Confirmation Modal -->
    <div 
      v-if="showResetConfirmation" 
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-2xl p-6 max-w-md w-full mx-4 shadow-2xl">
        <div class="text-center">
          <div class="mx-auto mb-4 w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
            <svg 
              class="w-6 h-6 text-red-600" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24" 
              xmlns="http://www.w3.org/2000/svg"
            >
              <path 
                stroke-linecap="round" 
                stroke-linejoin="round" 
                stroke-width="2" 
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              ></path>
            </svg>
          </div>
          <h2 class="text-lg font-bold text-gray-800 mb-3">Reset All Activity?</h2>
          <p class="text-xs text-gray-600 mb-4">
            This action will permanently clear all your learning progress, 
            recent sessions, and saved activities. This cannot be undone.
          </p>
          <div class="flex justify-center space-x-3">
            <button 
              @click="showResetConfirmation = false"
              class="px-4 py-2 bg-gray-200 text-gray-800 rounded-lg text-xs hover:bg-gray-300 transition-colors"
            >
              Cancel
            </button>
            <button 
              @click="resetAllActivity"
              class="px-4 py-2 bg-red-500 text-white rounded-lg text-xs hover:bg-red-600 transition-colors"
            >
              Reset Anyway
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

// Navigation Links
const navLinks = [
  { to: '/words', icon: '📖', label: 'Words' },
  { to: '/groups', icon: '🏷️', label: 'Groups' },
  { to: '/study-activities', icon: '📚', label: 'Study Activities' },
  { to: '/sessions', icon: '⏱️', label: 'Sessions' }
]

// Dashboard metrics
const dashboardMetrics = [
  { 
    label: 'Total Sessions', 
    value: 12, 
    color: 'text-blue-600' 
  },
  { 
    label: 'Study Activities', 
    value: 4, 
    color: 'text-green-600' 
  },
  { 
    label: 'Total Groups', 
    value: 3, 
    color: 'text-purple-600' 
  },
  { 
    label: 'Total Words', 
    value: 150, 
    color: 'text-orange-600' 
  }
]

// Mock data for study activities
const studyActivities = [
  { 
    id: 1, 
    name: 'Word Puzzle', 
    icon: '🧩',
    description: 'Unscramble Hindi words'
  },
  { 
    id: 2, 
    name: 'Sentence Builder', 
    icon: '📝',
    description: 'Create sentences in Hindi'
  },
  { 
    id: 3, 
    name: 'Listening', 
    icon: '👂',
    description: 'Improve listening skills'
  },
  { 
    id: 4, 
    name: 'Pronunciation', 
    icon: '🗣️',
    description: 'Practice Hindi pronunciation'
  }
]

// Mock data for recent sessions
const recentSessions = [
  { 
    id: 1, 
    activityName: 'Word Puzzle', 
    groupName: 'Everyday Words',
    score: 75, 
    endTime: '2024-02-10T14:30:00Z' 
  },
  { 
    id: 2, 
    activityName: 'Sentence Builder', 
    groupName: 'Family Terms',
    score: 60, 
    endTime: '2024-02-08T16:45:00Z' 
  }
]

// Format session date
const formatSessionDate = (dateString: string) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric', 
    hour: '2-digit', 
    minute: '2-digit' 
  })
}

// Score color classification
const getScoreClass = (score: number) => {
  if (score >= 80) return 'text-green-600'
  if (score >= 60) return 'text-yellow-600'
  return 'text-red-600'
}

const showResetConfirmation = ref(false)
const showUserPopup = ref(false)
const editedUserName = ref('')
const userName = ref('')

// Get user initial
const getUserInitial = computed(() => {
  if (userName.value) {
    return userName.value.charAt(0).toUpperCase()
  }
  return 'L'
})

// Toggle user popup
function toggleUserPopup() {
  showUserPopup.value = !showUserPopup.value
  if (showUserPopup.value) {
    editedUserName.value = userName.value || ''
  }
}

// Save user name
function saveUserName() {
  if (editedUserName.value.trim()) {
    userName.value = editedUserName.value.trim()
    toggleUserPopup()
  }
}

// Reset function to clear all data
function resetAllActivity() {
  // Reset metrics
  dashboardMetrics.forEach(metric => metric.value = 0)
  
  // Close confirmation modal
  showResetConfirmation.value = false
}

// Fetch dashboard metrics (to be implemented with actual data fetching)
async function fetchDashboardMetrics() {
  // TODO: Implement API call to fetch metrics
}

// Call on component mount
onMounted(() => {
  fetchDashboardMetrics()
})
</script>

<style scoped>
/* Prevent scrollbars from adding extra space */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
}

::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
