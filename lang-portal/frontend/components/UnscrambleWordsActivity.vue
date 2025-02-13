<template>
  <div class="p-8 bg-gradient-to-br from-blue-100 to-blue-300 min-h-[600px] flex flex-col">
    <div class="text-center mb-6">
      <h2 class="text-3xl font-bold text-gray-800">शब्द उलझाओ (Unscramble Words)</h2>
      <p class="text-gray-600">अक्षरों को सही क्रम में व्यवस्थित करें</p>
    </div>

    <div v-if="isLoading" class="flex-grow flex items-center justify-center">
      <div class="text-center">
        <div class="animate-pulse text-6xl mb-4">🧩</div>
        <p class="text-gray-600">चुनौतियाँ लोड हो रही हैं... (Loading challenges)</p>
      </div>
    </div>

    <div v-else-if="error || !challenges || challenges.length === 0" class="flex-grow flex items-center justify-center">
      <div class="text-center">
        <div class="text-6xl mb-4">😕</div>
        <h3 class="text-xl font-semibold text-gray-800 mb-2">कोई चुनौती उपलब्ध नहीं</h3>
        <p class="text-gray-600 mb-4">वर्तमान में कोई चुनौती नहीं मिल रही</p>
        <button 
          @click="refetch"
          class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition"
        >
          फिर से प्रयास करें (Try Again)
        </button>
      </div>
    </div>

    <div v-else-if="currentChallenge" class="flex-grow flex flex-col justify-between">
      <div class="space-y-6">
        <div class="flex justify-between items-center mb-4">
          <div class="flex space-x-2">
            <span class="text-sm font-medium text-gray-600">चुनौती:</span>
            <span class="text-sm font-semibold text-blue-600">{{ currentChallengeIndex + 1 }} / {{ challenges.length }}</span>
          </div>
          <div class="flex items-center space-x-2">
            <span class="text-sm font-medium text-gray-600">स्कोर:</span>
            <span class="text-sm font-semibold text-green-600">{{ score }}</span>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-md p-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="bg-blue-50 rounded-lg p-4 shadow-inner">
              <h3 class="text-lg font-semibold text-gray-800 mb-2">उलझे हुए शब्द (Scrambled Word)</h3>
              <div class="flex justify-center space-x-2 mb-4">
                <div 
                  v-for="(letter, index) in scrambledLetters" 
                  :key="index"
                  @click="selectLetter(index)"
                  class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-2xl font-bold text-blue-800 cursor-pointer hover:bg-blue-200 transition"
                  :class="{ 'opacity-50': letter.selected }"
                >
                  {{ letter.value }}
                </div>
              </div>
            </div>

            <div class="bg-green-50 rounded-lg p-4 shadow-inner">
              <h3 class="text-lg font-semibold text-gray-800 mb-2">आपका उत्तर (Your Answer)</h3>
              <div class="flex justify-center space-x-2 mb-4">
                <div 
                  v-for="(letter, index) in userAnswer" 
                  :key="index"
                  @click="removeLetter(index)"
                  class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center text-2xl font-bold text-green-800 cursor-pointer hover:bg-green-200 transition"
                >
                  {{ letter }}
                </div>
              </div>
            </div>
          </div>

          <div class="mt-4 text-center">
            <p class="text-xl font-semibold text-gray-800">
              मूल शब्द (Original Word): {{ currentChallenge.translation }}
            </p>
          </div>
        </div>

        <div v-if="feedbackMessage" class="text-center">
          <p 
            :class="{
              'text-green-600': feedbackMessage.type === 'success',
              'text-red-600': feedbackMessage.type === 'error'
            }"
          >
            {{ feedbackMessage.text }}
          </p>
        </div>
      </div>

      <div class="mt-6 flex justify-center space-x-4">
        <button 
          @click="checkAnswer"
          class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition disabled:opacity-50"
          :disabled="userAnswer.length !== scrambledLetters.length"
        >
          उत्तर जमा करें (Submit Answer)
        </button>
        <button 
          @click="resetChallenge"
          class="bg-gray-200 text-gray-700 px-6 py-3 rounded-lg hover:bg-gray-300 transition"
        >
          रीसेट करें (Reset)
        </button>
      </div>
    </div>

    <div v-else class="flex-grow flex items-center justify-center">
      <div class="text-center">
        <div class="text-6xl mb-4">🏁</div>
        <h3 class="text-xl font-semibold text-gray-800 mb-2">सत्र समाप्त (Session Completed)</h3>
        <p class="text-gray-600 mb-4">आपने सभी चुनौतियाँ पूरी कर ली हैं (You've completed all challenges)</p>
        <button 
          @click="popupControl.close"
          class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition"
        >
          सत्र बंद करें (Close Session)
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'

interface Challenge {
  id: string
  word: string
  scrambledWord: string
  translation: string
}

interface FeedbackMessage {
  type: 'success' | 'error'
  text: string
}

const props = defineProps({
  activityId: {
    type: String,
    required: true
  }
})

const popupControl = inject('popupControl', {
  close: () => {}
})

const score = ref(0)
const currentChallengeIndex = ref(0)
const feedbackMessage = ref<FeedbackMessage | null>(null)
const userAnswer = ref<string[]>([])

const { 
  data: challenges, 
  isLoading, 
  error,
  refetch 
} = useQuery({
  queryKey: ['unscrambleWordChallenges', props.activityId],
  queryFn: async () => {
    try {
      const response = await fetch('/api/words/random', {
        method: 'GET',
        headers: {
          'Accept': 'application/json'
        }
      })

      if (!response.ok) {
        console.error('Failed to fetch challenges:', response.status, response.statusText)
        throw new Error('Failed to fetch challenges')
      }

      const data = await response.json()
      console.log('Received challenges:', data)
      
      // Ensure we have at least 10 challenges
      const processedChallenges = data.slice(0, 10).map(challenge => ({
        ...challenge,
        scrambledWord: challenge.word.split('').sort(() => Math.random() - 0.5).join('')
      }))

      // If less than 10 challenges, repeat the challenges
      while (processedChallenges.length < 10) {
        processedChallenges.push(...processedChallenges.slice(0, 10 - processedChallenges.length))
      }

      return processedChallenges
    } catch (err) {
      console.error('Error in challenge fetch:', err)
      
      // Fallback challenges if API fails
      const fallbackChallenges = [
        { id: '1', word: 'नमस्ते', scrambledWord: 'तेमनास', translation: 'Hello' },
        { id: '2', word: 'भारत', scrambledWord: 'तभारा', translation: 'India' },
        { id: '3', word: 'प्यार', scrambledWord: 'रापय', translation: 'Love' },
        { id: '4', word: 'खुशी', scrambledWord: 'शीखु', translation: 'Happiness' },
        { id: '5', word: 'शांति', scrambledWord: 'तिशाां', translation: 'Peace' },
        { id: '6', word: 'दोस्त', scrambledWord: 'तसदो', translation: 'Friend' },
        { id: '7', word: 'सपना', scrambledWord: 'नापास', translation: 'Dream' },
        { id: '8', word: 'जीवन', scrambledWord: 'वनजी', translation: 'Life' },
        { id: '9', word: 'शक्ति', scrambledWord: 'तिकश', translation: 'Power' },
        { id: '10', word: 'आशा', scrambledWord: 'शाआ', translation: 'Hope' }
      ]

      return fallbackChallenges
    }
  },
  staleTime: 1000 * 60 * 5, // 5 minutes
  retry: 1,
  refetchOnWindowFocus: false
})

// Watch for challenges to be loaded and reset index
watch(challenges, (newChallenges) => {
  if (newChallenges && newChallenges.length > 0) {
    currentChallengeIndex.value = 0
  }
})

const currentChallenge = computed(() => 
  challenges.value && challenges.value[currentChallengeIndex.value]
)

const scrambledLetters = computed(() => 
  currentChallenge.value?.scrambledWord.split('').map((letter, index) => ({
    value: letter,
    selected: false,
    originalIndex: index
  })) || []
)

const selectLetter = (index: number) => {
  if (scrambledLetters.value[index].selected) return

  userAnswer.value.push(scrambledLetters.value[index].value)
  scrambledLetters.value[index].selected = true
}

const removeLetter = (index: number) => {
  const letter = userAnswer.value[index]
  userAnswer.value.splice(index, 1)
  
  const originalIndex = scrambledLetters.value.findIndex(
    l => l.value === letter && l.selected
  )
  if (originalIndex !== -1) {
    scrambledLetters.value[originalIndex].selected = false
  }
}

const resetChallenge = () => {
  userAnswer.value = []
  scrambledLetters.value.forEach(letter => letter.selected = false)
  feedbackMessage.value = null
}

const checkAnswer = () => {
  const isCorrect = userAnswer.value.join('').toLowerCase() === currentChallenge.value?.word.toLowerCase()
  
  if (isCorrect) {
    score.value += 5
    feedbackMessage.value = {
      type: 'success',
      text: 'सही उत्तर! बधाई हो! (Correct! Great job!)'
    }
  } else {
    feedbackMessage.value = {
      type: 'error',
      text: 'गलत उत्तर। पुनः प्रयास करें! (Incorrect answer. Try again!)'
    }
  }

  // Move to next challenge or end session
  setTimeout(() => {
    if (currentChallengeIndex.value < 9) {
      currentChallengeIndex.value++
      resetChallenge()
    } else {
      // Session completed
      popupControl.close()
    }
  }, 1000)
}
</script>

<style scoped>
/* Support for Devanagari script */
body {
  font-family: 'Noto Sans Devanagari', Arial, sans-serif;
}
</style>
