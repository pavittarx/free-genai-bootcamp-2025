import { ref } from 'vue'

// Group interface matching backend schema
export interface Group {
  id: number
  name: string
  description: string
  created_at: string
}

// Word interface to match backend schema
export interface GroupWord {
  id: number
  hindi: string
  hinglish: string
  english: string
  created_at: string
}

// Pagination interface
interface PaginationParams {
  page?: number
  pageSize?: number
  search?: string
  sortBy?: string
  sortDirection?: 'asc' | 'desc'
}

// Base API URL
const BASE_URL = 'http://localhost:3000/api'

export function useGroupService() {
  const groups = ref<Group[]>([])
  const group = ref<Group | null>(null)
  const groupWords = ref<GroupWord[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchGroups(params: PaginationParams = {}) {
    loading.value = true
    error.value = null

    try {
      const queryParams = new URLSearchParams()
      if (params.page) queryParams.append('page', params.page.toString())
      if (params.pageSize) queryParams.append('pageSize', params.pageSize.toString())
      if (params.search) queryParams.append('search', params.search)
      if (params.sortBy) queryParams.append('sortBy', params.sortBy)
      if (params.sortDirection) queryParams.append('sortDirection', params.sortDirection)

      const response = await fetch(`${BASE_URL}/groups?${queryParams}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch groups: ${response.status} - ${errorText}`)
      }

      const data = await response.json()
      
      // Debug logging
      console.log('Raw groups response:', data)

      // Correctly parse the nested group object
      const parsedGroups = data.groups.map((item: any) => ({
        id: item.group.id,
        name: item.group.name,
        description: item.group.description,
        created_at: item.group.created_at
      }))

      groups.value = parsedGroups
      
      return {
        groups: groups.value,
        totalCount: data.total || 0
      }
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while fetching groups'
      
      error.value = errorMessage
      groups.value = []
      
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getGroupById(id: number) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/groups/${id}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch group: ${response.status} - ${errorText}`)
      }

      const data = await response.json()
      
      // Correctly parse the nested group object
      group.value = {
        id: data.group.id,
        name: data.group.name,
        description: data.group.description,
        created_at: data.group.created_at
      }

      return group.value
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while fetching group'
      
      error.value = errorMessage
      group.value = null
      
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createGroup(groupData: { name: string; description: string }) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/groups`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(groupData)
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to create group: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while creating a group'
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateGroup(id: number, groupData: { name: string; description: string }) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/groups/${id}`, {
        method: 'PUT',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(groupData)
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to update group: ${response.status} - ${errorText}`)
      }

      return await response.json()
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while updating a group'
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteGroup(id: number) {
    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${BASE_URL}/groups/${id}`, {
        method: 'DELETE',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to delete group: ${response.status} - ${errorText}`)
      }

      return true
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : 'An unknown error occurred while deleting a group'
      
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  async function getGroupWords(groupId: number) {
    loading.value = true
    error.value = null
    groupWords.value = []

    try {
      const response = await fetch(`${BASE_URL}/words/groups/${groupId}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        }
      })
      
      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to fetch words for group ${groupId}: ${response.status} - ${errorText}`)
      }

      const data = await response.json()
      groupWords.value = data || []
      
      return groupWords.value
    } catch (err) {
      const errorMessage = err instanceof Error 
        ? err.message 
        : `An unknown error occurred while fetching words for group ${groupId}`
      
      error.value = errorMessage
      groupWords.value = []
      
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    groups,
    group,
    groupWords,
    loading,
    error,
    fetchGroups,
    getGroupById,
    createGroup,
    updateGroup,
    deleteGroup,
    getGroupWords
  }
}
