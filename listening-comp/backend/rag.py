import json
import random
from typing import Dict, Any, Optional, List

from .vector_db import TranscriptVectorDB
from .chat import OpenRouterChat

class RAGAssistant:
    """
    Retrieval Augmented Generation Assistant for Language Learning
    Combines vector store retrieval with generative AI capabilities
    """
    
    def __init__(self):
        """
        Initialize RAG components
        """
        self.vector_db = TranscriptVectorDB()
        self.chat = OpenRouterChat()
    
    def generate_learning_exercise(
        self, 
        topic: Optional[str] = None, 
        difficulty: str = 'Intermediate'
    ) -> Dict[str, Any]:
        """
        Generate a learning exercise using generative AI
        
        Args:
            topic (Optional[str]): Specific topic for exercise generation
            difficulty (str): Difficulty level of the exercise
        
        Returns:
            Dict containing structured learning exercise
        """
        # Difficulty-based complexity mapping
        difficulty_levels = {
            'शुरुआती': {
                'complexity': 'Simple, basic vocabulary',
                'sentence_structure': 'Short, straightforward sentences',
                'language_level': 'Beginner level Hindi'
            },
            'मध्यम': {
                'complexity': 'Moderate vocabulary, some idiomatic expressions',
                'sentence_structure': 'Mixed sentence lengths, some complex structures',
                'language_level': 'Intermediate level Hindi'
            },
            'उन्नत': {
                'complexity': 'Advanced vocabulary, nuanced expressions',
                'sentence_structure': 'Complex, varied sentence structures',
                'language_level': 'Advanced level Hindi'
            }
        }
        
        # Comprehensive list of context topics
        context_topics = [
            "Everyday conversations", 
            "Office dialogues", 
            "Social situations", 
            "Educational discussions", 
            "Travel and tourism", 
            "Family and relationships", 
            "Technology and innovation", 
            "Environment and nature", 
            "Health and wellness", 
            "Art and culture"
        ]
        
        # Select topic if not provided
        selected_topic = topic or random.choice(context_topics)
        
        # Prepare the prompt for exercise generation
        prompt = f"""
        Generate a structured Hindi language learning exercise:

        Exercise Generation Guidelines:
        - Context Topic: {selected_topic}
        - Difficulty Level: {difficulty}
        - Language Complexity: {difficulty_levels.get(difficulty, difficulty_levels['मध्यम'])['complexity']}
        - Sentence Structure: {difficulty_levels.get(difficulty, difficulty_levels['मध्यम'])['sentence_structure']}
        - Language Level: {difficulty_levels.get(difficulty, difficulty_levels['मध्यम'])['language_level']}

        MANDATORY OUTPUT FORMAT:
        {{
            "introduction": "Brief context-setting phrase (15-20 words)",
            "dialogue": "Conversation text in natural Hindi",
            "question": "A specific comprehension or language question",
            "options": ["Option 1", "Option 2", "Option 3", "Option 4"],
            "answer": "Correct answer from the options"
        }}

        CRITICAL INSTRUCTIONS:
        - Use natural, conversational Hindi
        - Ensure cultural authenticity
        - Create engaging, contextually relevant content
        - Avoid generic or placeholder text
        """
        
        # Generate exercise using chat model
        try:
            # Generate response
            response = self.chat.generate_response(prompt)
            
            # Parse the JSON response
            try:
                exercise = json.loads(response)
            except json.JSONDecodeError:
                # Fallback parsing
                exercise = self._parse_exercise_response(response)
            
            # Validate exercise structure
            if not exercise or not isinstance(exercise, dict):
                raise ValueError("Invalid exercise structure generated")
            
            # Ensure options are generated if not present
            if not exercise.get('options') or len(exercise['options']) < 4:
                exercise['options'] = self._generate_multiple_options(
                    exercise.get('answer', '')
                )
            
            return exercise
        
        except Exception as e:
            # Comprehensive error logging
            print(f"Error generating learning exercise: {e}")
            
            # Fallback exercise generation
            return {
                "introduction": "एक रोचक संवाद",
                "dialogue": "यह एक सामान्य संवाद है जो हिंदी सीखने में मदद करेगा।",
                "question": "इस संवाद का मुख्य विषय क्या है?",
                "options": [
                    "शिक्षा", 
                    "यात्रा", 
                    "परिवार", 
                    "तकनीक"
                ],
                "answer": "शिक्षा"
            }
    
    def _parse_exercise_response(self, response: str) -> Dict[str, Any]:
        """
        Parse exercise response when JSON parsing fails
        
        Args:
            response (str): Raw response from chat model
        
        Returns:
            Dict with exercise structure
        """
        # Extract JSON-like content
        import re
        
        json_match = re.search(r'\{.*\}', response, re.DOTALL | re.MULTILINE)
        
        if json_match:
            try:
                return json.loads(json_match.group(0))
            except json.JSONDecodeError:
                pass
        
        # Fallback parsing
        return {
            "introduction": response[:100],
            "dialogue": response,
            "question": "",
            "options": [],
            "answer": response
        }
    
    def _generate_multiple_options(self, answer: str) -> List[str]:
        """
        Generate multiple options based on the correct answer
        
        Args:
            answer (str): Correct answer
        
        Returns:
            List of multiple choice options
        """
        prompt = f"""
        Generate 3 plausible but incorrect options for the following answer in Hindi:
        
        Correct Answer: {answer}
        
        Guidelines:
        - Create options that sound similar but are incorrect
        - Ensure options are in Hindi
        - Make sure the correct answer is not repeated
        - Options should be concise
        """
        
        try:
            response = self.chat.generate_response(prompt)
            # Split response into options
            options = response.split('\n')[:3]
            
            # Ensure we have 4 options total (3 incorrect + 1 correct)
            full_options = options + [answer]
            
            # Shuffle options to randomize correct answer position
            random.shuffle(full_options)
            
            return full_options
        except Exception:
            # Fallback options if generation fails
            return [
                f"Incorrect version 1 of {answer}",
                f"Incorrect version 2 of {answer}",
                f"Incorrect version 3 of {answer}",
                answer
            ]
    
    def retrieve_context(
        self, 
        query: str, 
        top_k: int = 3
    ) -> Dict[str, Any]:
        """
        Retrieve contextual information from vector store
        
        Args:
            query (str): Search query
            top_k (int): Number of top results to retrieve
        
        Returns:
            Dict with retrieved context and metadata
        """
        # Perform semantic search
        search_results = self.vector_db.search_transcripts(query, top_k=top_k)
        
        # Format and return results
        return {
            "query": query,
            "results": search_results
        }
    
    def generate_response_with_context(
        self, 
        query: str, 
        context: Optional[Dict[str, Any]] = None
    ) -> str:
        """
        Generate a response using retrieved context
        
        Args:
            query (str): User's query
            context (Optional[Dict]): Optional pre-retrieved context
        
        Returns:
            Generated response as a string
        """
        # If no context provided, retrieve it
        if context is None:
            context = self.retrieve_context(query)
        
        # Prepare prompt with context
        prompt = f"""
        Context: {json.dumps(context['results'], indent=2)}
        
        Query: {query}
        
        Using the provided context, generate a comprehensive and 
        informative response that directly addresses the query.
        """
        
        # Generate response using OpenRouterChat
        return self.chat.generate_response(prompt)
    
    def create_interactive_quiz(
        self, 
        topic: Optional[str] = None, 
        num_questions: int = 5
    ) -> Dict[str, Any]:
        """
        Create an interactive quiz using vector store transcripts
        
        Args:
            topic (Optional[str]): Specific topic for quiz
            num_questions (int): Number of questions to generate
        
        Returns:
            Dict containing quiz structure
        """
        quiz = {
            "topic": topic or "General Language Learning",
            "questions": []
        }
        
        # Generate multiple questions
        for _ in range(num_questions):
            # Generate individual exercise
            question = self.generate_learning_exercise(
                topic=topic, 
                difficulty='Intermediate'
            )
            quiz["questions"].append(question)
        
        return quiz
