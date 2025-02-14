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
            <p class="text-gray-600 mb-4">{{ error }}</p>
          </div>
        </div>

        <div v-else-if="activityCompleted" class="flex-grow flex items-center justify-center">
          <div class="text-center">
            <div class="text-6xl mb-4">🏆</div>
            <h3 class="text-xl font-semibold text-gray-800 mb-2">गतिविधि पूरी हुई</h3>
            <p class="text-2xl font-bold text-green-600 mb-4">आपका कुल स्कोर: {{ score }} / 50</p>
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
                <span class="text-sm font-medium text-gray-600">चुनौती:</span>
                <span class="text-sm font-semibold text-blue-600">{{ currentChallengeIndex + 1 }} / 10</span>
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
              :disabled="userAnswer.length !== currentChallenge.scrambledWord.length"
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
import { useQuery } from '@tanstack/vue-query'
import { apiService, type Word, type Session } from '~/services/api'

const props = defineProps<{
  activityId: string
}>()

const popupControl = inject('popupControl', {
  close: () => {}
})

// Session and score management
const session = ref<Session | null>(null)
const score = ref(0)
const currentChallengeIndex = ref(0)
const activityCompleted = ref(false)

// Feedback and interaction state
const feedbackMessage = ref<{ type: 'success' | 'error'; text: string } | null>(null)
const userAnswer = ref<string[]>([])

// Initialize session on component mount
onMounted(async () => {
  try {
    session.value = await apiService.createSession(props.activityId)
    console.log('Session created:', session.value)
  } catch (error) {
    console.error('Failed to create session:', error)
    feedbackMessage.value = {
      type: 'error',
      text: 'सत्र शुरू करने में त्रुटि (Error starting session)'
    }
  }
})

// Fetch challenges
const { 
  data: challenges, 
  isLoading, 
  error,
  refetch 
} = useQuery<(Word & { scrambledWord: string })[]>({
  queryKey: ['unscrambleWordChallenges', props.activityId],
  queryFn: async () => {
    try {
      const words = await apiService.getRandomWords(10)
      
      return words.map(word => ({
        ...word,
        scrambledWord: word.word.split('').sort(() => Math.random() - 0.5).join('')
      }))
    } catch (err) {
      console.error('Challenge fetch error:', err)
      throw err
    }
  },
  staleTime: 1000 * 60 * 5, // 5 minutes
  retry: 1,
  refetchOnWindowFocus: false
})

const currentChallenge = computed(() => 
  challenges.value && challenges.value[currentChallengeIndex.value]
)

const checkAnswer = async () => {
  if (!currentChallenge.value || !session.value) return

  const userInput = userAnswer.value.join('')
  const correctAnswer = currentChallenge.value.word
  const isCorrect = userInput.toLowerCase() === correctAnswer.toLowerCase()
  const challengeScore = isCorrect ? 5 : 0
  
  try {
    // Submit activity result
    await apiService.submitActivity({
      session_id: session.value.id,
      activity_id: props.activityId,
      challenge: correctAnswer,
      answer: correctAnswer,
      input: userInput,
      score: challengeScore
    })

    // Update local score
    if (isCorrect) {
      score.value += challengeScore
      feedbackMessage.value = {
        type: 'success',
        text: 'सही उत्तर! बधाई हो! (Correct! Great job!)'
      }
    } else {
      feedbackMessage.value = {
        type: 'error',
        text: `गलत उत्तर। सही उत्तर था: ${correctAnswer} (Wrong answer. The correct word was: ${correctAnswer})`
      }
    }

    // Move to next challenge or end session
    setTimeout(async () => {
      if (currentChallengeIndex.value < 9) {
        currentChallengeIndex.value++
        resetChallenge()
      } else {
        // Final challenge completed
        await endActivity()
      }
    }, 1000)
  } catch (error) {
    console.error('Failed to submit activity:', error)
    feedbackMessage.value = {
      type: 'error',
      text: 'Error submitting answer. Please try again.'
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
      resetChallenge()
    } else {
      // Final challenge skipped
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
  if (!session.value) return

  try {
    // Close session with final score
    await apiService.closeSession(session.value.id, score.value)
    
    // Mark activity as completed
    activityCompleted.value = true
  } catch (error) {
    console.error('Failed to close session:', error)
  }
}

const resetChallenge = () => {
  userAnswer.value = []
  feedbackMessage.value = null
}

const selectLetter = (index: number) => {
  const letter = currentChallenge.value.scrambledWord[index]
  userAnswer.value.push(letter)
}

const removeLetter = (index: number) => {
  userAnswer.value.splice(index, 1)
}

const handleActivityEnd = () => {
  popupControl.close()
}
</script>

<style scoped>
/* Minimal styling to ensure readability */
</style>
