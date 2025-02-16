import axios from 'axios'

const API_BASE_URL = 'http://localhost:3000/api'

export interface Word {
  id: number
  hindi: string
  scrambled: string
  hinglish: string
  english: string
  created_at: string
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
    console.log('Attempting to create session', {
      activityId: activity_id,
      activityIdType: typeof activity_id,
      baseURL: API_BASE_URL,
      timestamp: new Date().toISOString()
    })

    // Validate input
    if (!activity_id) {
      const error = new Error('Activity ID is required')
      console.error('Session creation validation error:', error)
      throw error
    }

    try {
      return api.post('/sessions', { 
        activity_id: activity_id 
      }).then(response => {
        console.log('Session creation response details:', {
          status: response.status,
          headers: response.headers,
          data: response.data,
          dataType: typeof response.data
        })

        // Validate response
        if (!response.data) {
          throw new Error('No session data returned')
        }

        // Validate session structure
        const session = response.data as Session
        if (!session.id || !session.activity_id) {
          throw new Error('Invalid session structure')
        }

        return session
      })
    } catch (error: any) {
      console.error('Detailed session creation error:', {
        message: error.message,
        response: error.response?.data,
        status: error.response?.status,
        headers: error.response?.headers
      })
      throw error
    }
  },

  async closeSession(session_id: number, score: number): Promise<Session> {
    const response = await api.put('/sessions', { 
      session_id,
      score
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
      console.log('Attempting to fetch random words', { 
        baseURL: API_BASE_URL, 
        route: '/words/random', 
        count 
      })
      
      const response = await api.get('/words/random', {
        params: { count }
      })
      
      console.log('Full API response:', {
        status: response.status,
        headers: response.headers,
        dataType: typeof response.data,
        dataIsArray: Array.isArray(response.data),
        dataLength: response.data?.length,
        firstWordDetails: response.data?.[0] ? {
          id: response.data[0].id,
          hindi: response.data[0].hindi
        } : 'No first word'
      })
      
      // Validate response
      if (!response.data) {
        throw new Error('No data returned from server')
      }
      
      // Ensure response is an array
      const words = Array.isArray(response.data) 
        ? response.data 
        : [response.data]
      
      // Validate words
      if (words.length === 0) {
        throw new Error('No words returned from server')
      }
      
      // Validate each word
      const validWords = words.filter(word => 
        word && 
        typeof word === 'object' && 
        word.id !== undefined && 
        word.hindi !== undefined
      )
      
      if (validWords.length === 0) {
        throw new Error('No valid words found in server response')
      }
      
      return validWords
    } catch (error: any) {
      console.error('Detailed error fetching random words:', {
        message: error.message,
        response: error.response?.data,
        status: error.response?.status,
        headers: error.response?.headers
      })
      throw error
    }
  }
}
