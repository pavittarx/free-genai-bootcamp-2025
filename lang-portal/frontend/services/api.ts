import axios from 'axios'

const API_BASE_URL = 'http://localhost:3000/api'

export interface Word {
  id: number
  word: string
  translation: string
  hinglish: string
  group_id?: number
  difficulty?: string
}

export interface Session {
  id: number
  activity_id: string
  start_time: string
  end_time?: string
  score: number
  status: 'in_progress' | 'completed'
}

export interface SessionActivity {
  session_id: number
  activity_id: string
  challenge: string
  answer: string
  input: string
  score: number
}

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const apiService = {
  // Session Management
  async createSession(activity_id: string): Promise<Session> {
    console.log('Creating session for activity:', activity_id)
    console.log('API Base URL:', API_BASE_URL)
    try {
      const response = await api.post('/sessions', { activity_id })
      console.log('Session creation response:', response.data)
      return response.data
    } catch (error) {
      console.error('Session creation error:', error)
      throw error
    }
  },

  async closeSession(session_id: number, score: number): Promise<Session> {
    const response = await api.put(`/sessions/${session_id}`, { 
      score,
      end_time: new Date().toISOString(),
      status: 'completed'
    })
    return response.data
  },

  // Activity Submission
  async submitActivity(activity: SessionActivity): Promise<void> {
    await api.post('/session-activity', activity)
  },

  // Word Fetching
  async getRandomWords(count: number = 10): Promise<Word[]> {
    try {
      const response = await api.get('/words/random', { 
        params: { count } 
      })

      // Handle different possible response structures
      const data = response.data
      
      // If data is an array, return it directly
      if (Array.isArray(data)) {
        return data
      }
      
      // If data has a 'words' or 'data' property that is an array, return that
      if (data.words && Array.isArray(data.words)) {
        return data.words
      }
      
      if (data.data && Array.isArray(data.data)) {
        return data.data
      }

      // Fallback to default words if no valid array found
      console.warn('Unexpected API response structure:', data)
      return [
        { id: 1, word: 'नमस्ते', translation: 'Hello', hinglish: 'Namaste' },
        { id: 2, word: 'भारत', translation: 'India', hinglish: 'Bharat' },
        { id: 3, word: 'प्यार', translation: 'Love', hinglish: 'Pyaar' },
        { id: 4, word: 'खुशी', translation: 'Happiness', hinglish: 'Khushi' },
        { id: 5, word: 'शांति', translation: 'Peace', hinglish: 'Shanti' },
        { id: 6, word: 'दोस्त', translation: 'Friend', hinglish: 'Dost' },
        { id: 7, word: 'सपना', translation: 'Dream', hinglish: 'Sapna' },
        { id: 8, word: 'जीवन', translation: 'Life', hinglish: 'Jeevan' },
        { id: 9, word: 'शक्ति', translation: 'Power', hinglish: 'Shakti' },
        { id: 10, word: 'आशा', translation: 'Hope', hinglish: 'Asha' }
      ].slice(0, count)
    } catch (error) {
      console.error('Error fetching random words:', error)
      
      // Fallback to default words if API call fails
      return [
        { id: 1, word: 'नमस्ते', translation: 'Hello', hinglish: 'Namaste' },
        { id: 2, word: 'भारत', translation: 'India', hinglish: 'Bharat' },
        { id: 3, word: 'प्यार', translation: 'Love', hinglish: 'Pyaar' },
        { id: 4, word: 'खुशी', translation: 'Happiness', hinglish: 'Khushi' },
        { id: 5, word: 'शांति', translation: 'Peace', hinglish: 'Shanti' },
        { id: 6, word: 'दोस्त', translation: 'Friend', hinglish: 'Dost' },
        { id: 7, word: 'सपना', translation: 'Dream', hinglish: 'Sapna' },
        { id: 8, word: 'जीवन', translation: 'Life', hinglish: 'Jeevan' },
        { id: 9, word: 'शक्ति', translation: 'Power', hinglish: 'Shakti' },
        { id: 10, word: 'आशा', translation: 'Hope', hinglish: 'Asha' }
      ].slice(0, count)
    }
  }
}
