// Mock service for group-related API calls
// In a real application, replace with actual API calls

export async function fetchGroupDetails(groupId: number) {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 500))

  return {
    id: groupId,
    name: 'Basic Vocabulary',
    description: 'Essential words for beginners to learn Hindi and Hinglish',
    wordCount: 50,
    createdAt: new Date('2024-01-15')
  }
}

export async function fetchGroupWords(groupId: number) {
  // Simulate API delay
  await new Promise(resolve => setTimeout(resolve, 500))

  return [
    { 
      id: 1, 
      hindi: 'नमस्ते', 
      hinglish: 'Namaste', 
      english: 'Hello', 
      createdAt: new Date('2024-01-15') 
    },
    { 
      id: 2, 
      hindi: 'धन्यवाद', 
      hinglish: 'Dhanyavaad', 
      english: 'Thank you', 
      createdAt: new Date('2024-01-20') 
    },
    // Add more mock words here
  ]
}
