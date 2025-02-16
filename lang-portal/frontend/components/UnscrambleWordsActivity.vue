<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4">
    <div class="bg-gradient-to-br from-blue-100 to-blue-300 w-full max-w-2xl rounded-xl shadow-2xl overflow-hidden">
      <div class="p-8 flex flex-col min-h-[600px]">
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

        <div v-else-if="error" class="flex-grow flex items-center justify-center">
          <div class="text-center">
            <div class="text-6xl mb-4">😕</div>
            <h3 class="text-xl font-semibold text-gray-800 mb-2">लोड करने में त्रुटि</h3>
            <p class="text-gray-600 mb-4">{{ errorMessage }}</p>
          </div>
        </div>

        <div v-else-if="activityCompleted" class="flex-grow flex items-center justify-center">
          <div class="text-center">
            <div class="text-6xl mb-4">🏆</div>
            <h3 class="text-xl font-semibold text-gray-800 mb-2">गतिविधि पूरी हुई (Activity Completed)</h3>
            <p class="text-2xl font-bold text-green-600 mb-4">
              आपका कुल स्कोर (Total Score): {{ score }} / 50
            </p>
            <div class="flex justify-center space-x-4">
              <button 
                @click="handleActivityEnd"
                class="bg-blue-500 text-white px-6 py-3 rounded-lg hover:bg-blue-600 transition"
              >
                समाप्त करें (Close)
              </button>
            </div>
          </div>
        </div>

        <div v-else-if="currentChallenge" class="flex-grow flex flex-col justify-between">
          <div class="space-y-6">
            <div class="flex justify-between items-center mb-4">
              <div class="flex space-x-2">
                <span class="text-sm font-medium text-gray-600">चुनौती (Challenge):</span>
                <span class="text-sm font-semibold text-blue-600">
                  {{ currentChallengeIndex + 1 }} / 10
                </span>
              </div>
              <div class="flex items-center space-x-2">
                <span class="text-sm font-medium text-gray-600">स्कोर (Score):</span>
                <span class="text-sm font-semibold text-green-600">{{ score }}</span>
              </div>
            </div>

            <div class="bg-white rounded-xl shadow-md p-6">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div class="bg-blue-50 rounded-lg p-4 shadow-inner">
                  <h3 class="text-lg font-semibold text-gray-800 mb-2">उलझे हुए शब्द (Scrambled Word)</h3>
                  <div class="flex justify-center space-x-2 mb-4">
                    <div 
                      v-for="(letter, index) in currentChallenge.scrambledWord.split('')" 
                      :key="index"
                      @click="selectLetter(index)"
                      class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center text-2xl font-bold text-blue-800 cursor-pointer hover:bg-blue-200 transition"
                    >
                      {{ letter }}
                    </div>
                  </div>
                </div>

                <div class="bg-green-50 rounded-lg p-4 shadow-inner">
                  <h3 class="text-lg font-semibold text-gray-800 mb-2">आपका उत्तर (Your Answer)</h3>
                  <div class="flex justify-center space-x-2 mb-4">
                    <div 
                      v-for="(letter, index) in userInput" 
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
                  मूल शब्द (Original Word): {{ currentChallenge.word }}
                </p>
                <p class="text-lg text-gray-600">
                  हिंग्लिश (Hinglish): {{ currentWord?.hinglish || '' }}
                </p>
                <p class="text-lg text-gray-600">
                  अंग्रेजी (English): {{ currentWord?.english || '' }}
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
              :disabled="userInput.length !== currentChallenge.scrambledWord.length"
            >
              उत्तर जमा करें (Submit Answer)
            </button>
            <button 
              @click="skipChallenge"
              class="bg-gray-200 text-gray-700 px-6 py-3 rounded-lg hover:bg-gray-300 transition"
            >
              छोड़ें (Skip)
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, inject, onMounted } from 'vue'
import { useQuery, type UseQueryOptions } from '@tanstack/vue-query'
import { apiService, type Word, type Session } from '~/services/api'
import axios from 'axios'

// Props definition
const props = defineProps<{
  activityId: string
}>()

// Session and score management
const session = ref<Session | null>(null)
const score = ref(0)
const currentChallengeIndex = ref(0)
const activityCompleted = ref(false)

// Feedback and interaction state
const feedbackMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null)
const userInput = ref<string[]>([])

// Custom type guard to ensure Word type
function isWord(data: unknown): data is Word {
  return (
    typeof data === 'object' && 
    data !== null && 
    'id' in data &&
    'hindi' in data &&
    'scrambled' in data &&
    'hinglish' in data &&
    'english' in data
  )
}

// Query options with explicit typing
const queryOptions: UseQueryOptions<Word, Error> = {
  queryKey: ['randomWord', currentChallengeIndex.value],
  queryFn: async () => {
    console.log('Fetching random word - Full details:', {
      challengeIndex: currentChallengeIndex.value,
      timestamp: new Date().toISOString()
    })
    
    try {
      const words = await apiService.getRandomWords(1)
      
      if (!Array.isArray(words) || words.length === 0) {
        throw new Error('No words returned from server')
      }
      
      const word = words[0]
      if (!isWord(word)) {
        throw new Error('Invalid word data')
      }
      
      return word
    } catch (error) {
      console.error('Error fetching word:', error)
      throw error
    }
  },
  enabled: currentChallengeIndex.value < 10, // Limit to 10 challenges
  retry: 3,
  retryDelay: 1000,
  staleTime: 0
}

// Use query with explicit typing
const { 
  data: currentWord, 
  refetch, 
  isLoading, 
  error, 
  isError 
} = useQuery<Word, Error>(queryOptions)

// Computed property to get current challenge
const currentChallenge = computed(() => {
  const word = currentWord.value
  if (!word) return null
  return {
    word: word.hindi,
    scrambledWord: word.scrambled
  }
})

// Computed error message
const errorMessage = computed(() => {
  if (isError.value && error.value) {
    return error.value.message || 'अज्ञात त्रुटि (Unknown error)'
  }
  return null
})

interface PopupControl {
  close?: () => void
}

const popupControl = inject<PopupControl>('popupControl', {})

const selectLetter = (index: number) => {
  if (!currentChallenge.value) return
  // Safely handle letter selection
  const scrambledWord = currentChallenge.value.scrambledWord
  const letter = scrambledWord[index]
  if (letter && !userInput.value.includes(letter)) {
    userInput.value.push(letter)
  }
}

const checkAnswer = async () => {
  // Comprehensive null checks
  if (!currentChallenge.value || !session.value || !currentWord.value) {
    console.warn('Cannot check answer: missing data')
    return
  }

  const userInputValue = userInput.value.join('')
  const correctAnswer = currentChallenge.value.word
  const isCorrect = userInputValue.toLowerCase() === correctAnswer.toLowerCase()
  const challengeScore = isCorrect ? 5 : 0
  
  try {
    // Safely access properties with type guard
    const word = currentWord.value
    const hinglish = word?.hinglish || ''
    const english = word?.english || ''

    await apiService.submitActivity({
      session_id: session.value.id,
      activity_id: props.activityId,
      challenge: currentChallenge.value.word,
      answer: userInputValue,
      input: userInputValue,
      score: challengeScore
    })

    feedbackMessage.value = {
      type: isCorrect ? 'success' : 'error',
      text: isCorrect 
        ? 'बधाई हो! आपका उत्तर सही है। (Congratulations! Your answer is correct.)' 
        : 'क्षमा करें, यह उत्तर सही नहीं है। (Sorry, this answer is incorrect.)'
    }

    // Additional feedback with Hinglish and English translations
    if (isCorrect) {
      console.log(`Word Details - Hinglish: ${hinglish}, English: ${english}`)
    }

    // Increment score or handle challenge progression
    if (isCorrect) {
      score.value += challengeScore
    }

    // Update to handle final challenge
    if (currentChallengeIndex.value === 9) {
      activityCompleted.value = true
      await endActivity()
    } else {
      currentChallengeIndex.value++
      resetChallenge()
    }
  } catch (err: unknown) {
    const errorMessage = err instanceof Error ? err.message : String(err)
    console.error('Answer submission error:', errorMessage)
    feedbackMessage.value = {
      type: 'error',
      text: `उत्तर जमा करने में त्रुटि: ${errorMessage}`
    }
  }
}

const skipChallenge = async () => {
  if (!currentChallenge.value || !session.value) return

  try {
    // Submit skipped activity result with 0 score
    await apiService.submitActivity({
      session_id: session.value.id,
      activity_id: props.activityId,
      challenge: currentChallenge.value.word,
      answer: currentChallenge.value.word,
      input: '',
      score: 0
    })

    // Move to next challenge or end session
    if (currentChallengeIndex.value < 9) {
      currentChallengeIndex.value++
      // Fetch next word
      await refetch()
      resetChallenge()
    } else {
      await endActivity()
    }
  } catch (error) {
    console.error('Failed to skip challenge:', error)
    feedbackMessage.value = {
      type: 'error',
      text: 'Error skipping challenge. Please try again.'
    }
  }
}

const endActivity = async () => {
  try {
    if (session.value) {
      // Close the session with final score
      const finalSession = await apiService.closeSession(session.value.id, score.value)
      
      // Update session with final details
      session.value = finalSession
      
      // Ensure activity is marked as completed
      activityCompleted.value = true
    }
  } catch (error) {
    console.error('Error ending activity:', error)
    
    // More detailed error handling
    if (axios.isAxiosError(error)) {
      const errorMessage = error.response?.data?.message || 
                           error.response?.data?.error || 
                           'Unknown server error'
      
      feedbackMessage.value = {
        type: 'error',
        text: `Could not complete activity: ${errorMessage}`
      }
    } else {
      feedbackMessage.value = {
        type: 'error',
        text: 'Could not complete activity. Please try again.'
      }
    }
  }
}

const resetChallenge = () => {
  userInput.value = []
  feedbackMessage.value = null
}

const removeLetter = (index: number) => {
  userInput.value.splice(index, 1)
}

const handleActivityEnd = () => {
  // Safe injection with type checking
  if (typeof popupControl.close === 'function') {
    popupControl.close()
  } else {
    console.warn('Popup close method not available')
    // Fallback method to close the activity
    activityCompleted.value = true
  }
}

// Initialize session on component mount
onMounted(async () => {
  try {
    // Create session
    session.value = await apiService.createSession(props.activityId)
    
    console.log('Session created:', session.value)
    
    // Fetch first word
    await refetch()
  } catch (error) {
    console.error('Initialization error:', error)
    
    // Set a generic error message
    feedbackMessage.value = {
      type: 'error',
      text: 'सत्र शुरू करने में त्रुटि (Error starting session)'
    }
    
    // Close the popup on critical error
    if (typeof popupControl.close === 'function') {
      popupControl.close()
    }
  }
})
</script>

<style scoped>
/* Minimal styling to ensure readability */
</style>
