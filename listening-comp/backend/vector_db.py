import os
from typing import List, Dict, Any, Optional
import chromadb
from chromadb.config import Settings
from sentence_transformers import SentenceTransformer
import json
import glob
import random
from backend.chat import OpenRouterChat

class TranscriptVectorDB:
    def __init__(self, collection_name: str = "transcripts"):
        """
        Initialize the vector database for transcripts.
        
        Args:
            collection_name (str): Name of the Chroma collection to use.
        """
        # Ensure the chroma_db directory exists
        db_path = os.path.join(os.path.dirname(__file__), "chroma_db")
        os.makedirs(db_path, exist_ok=True)

        # Create a persistent Chroma client with explicit configuration
        try:
            self.client = chromadb.PersistentClient(
                path=db_path,
                settings=Settings(
                    allow_reset=True,  # Allow resetting the database if needed
                    anonymized_telemetry=False  # Disable telemetry
                )
            )
            
            # Explicitly create or get the collection
            self.collection = self.client.get_or_create_collection(
                name=collection_name, 
                metadata={"description": "Transcript embeddings for semantic search"}
            )
            
        except Exception as e:
            print(f"Error initializing Chroma client: {e}")
            # Fallback to in-memory client if persistent client fails
            self.client = chromadb.Client()
            self.collection = self.client.get_or_create_collection(
                name=collection_name, 
                metadata={"description": "Transcript embeddings for semantic search"}
            )
        
        # Initialize embedding model
        self.embedding_model = SentenceTransformer('all-MiniLM-L6-v2')

    def reset_collection(self, collection_name: str = "transcripts"):
        """
        Reset the entire collection, useful for clearing out old data.
        
        Args:
            collection_name (str): Name of the collection to reset
        """
        try:
            # Delete existing collection if it exists
            self.client.delete_collection(name=collection_name)
            
            # Recreate the collection
            self.collection = self.client.get_or_create_collection(
                name=collection_name, 
                metadata={"description": "Transcript embeddings for semantic search"}
            )
            return True
        except Exception as e:
            print(f"Error resetting collection: {e}")
            return False

    def add_transcripts(self, transcripts: List[Dict[str, Any]]):
        """
        Add transcripts to the vector database.
        
        Args:
            transcripts (List[Dict]): List of transcript dictionaries
        """
        # Prepare data for Chroma
        ids = []
        embeddings = []
        metadatas = []
        documents = []

        for idx, transcript in enumerate(transcripts):
            # Generate a unique ID
            doc_id = f"transcript_{idx}"
            ids.append(doc_id)
            
            # Create document text (combine relevant fields)
            doc_text = f"{transcript.get('title', '')} {transcript.get('content', '')}"
            documents.append(doc_text)
            
            # Generate embedding
            embedding = self.embedding_model.encode(doc_text).tolist()
            embeddings.append(embedding)
            
            # Store metadata
            metadatas.append({
                "title": transcript.get('title', ''),
                "source": transcript.get('source', ''),
                "timestamp": transcript.get('timestamp', '')
            })

        # Add to Chroma collection
        self.collection.add(
            ids=ids,
            embeddings=embeddings,
            metadatas=metadatas,
            documents=documents
        )

    def search_transcripts(self, query: str, n_results: int = 5) -> List[Dict[str, Any]]:
        """
        Search transcripts using semantic search.
        
        Args:
            query (str): Search query
            n_results (int): Number of results to return
        
        Returns:
            List of matching transcript results
        """
        # Generate query embedding
        query_embedding = self.embedding_model.encode(query).tolist()
        
        # Perform search
        results = self.collection.query(
            query_embeddings=[query_embedding],
            n_results=n_results
        )
        
        # Format results
        formatted_results = []
        for i in range(len(results['ids'][0])):
            result = {
                'id': results['ids'][0][i],
                'document': results['documents'][0][i],
                'metadata': results['metadatas'][0][i],
                'distance': results['distances'][0][i]
            }
            formatted_results.append(result)
        
        return formatted_results

    def generate_learning_exercise(
            self, 
            difficulty: str = 'मध्यम'
        ) -> Dict[str, Any]:
        """
        Generate a learning exercise based on difficulty level.
        
        Args:
            difficulty (str): Difficulty level of the exercise
        
        Returns:
            Dict containing structured learning exercise
        """
        # Initialize chat client
        chat_client = OpenRouterChat()
        
        # Comprehensive list of context topics for diverse exercise generation
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
        
        # Select a random context topic
        selected_topic = random.choice(context_topics)
        
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
            response = chat_client.generate_response(prompt)
            
            # Parse the JSON response
            from backend.structured_data import parse_json_response
            exercise = parse_json_response(response)
            
            # Validate exercise structure
            if not exercise or not isinstance(exercise, dict):
                raise ValueError("Invalid exercise structure generated")
            
            # Ensure options are generated if not present
            if not exercise.get('options') or len(exercise['options']) < 4:
                from backend.structured_data import generate_multiple_options
                exercise['options'] = generate_multiple_options(
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

    @classmethod
    def read_structured_transcripts(cls, directory: str = "./backend/structured_transcripts") -> List[Dict[str, Any]]:
        """
        Read structured transcripts from JSON files in a specified directory.
        
        Args:
            directory (str): Path to the directory containing structured transcript JSON files
        
        Returns:
            List of transcript dictionaries
        """
        transcripts = []
        
        # Find all JSON files in the specified directory
        transcript_files = glob.glob(os.path.join(directory, "*.json"))
        
        for file_path in transcript_files:
            try:
                with open(file_path, 'r', encoding='utf-8') as f:
                    transcript_data = json.load(f)
                    
                    # Standardize transcript format
                    transcript = {
                        "title": transcript_data.get("introduction", os.path.basename(file_path)),
                        "content": " ".join([
                            transcript_data.get("dialogue", ""),
                            transcript_data.get("answer", "")
                        ]),
                        "source": "listening_comprehension",
                        "timestamp": "",
                        "original_file": file_path,
                        "question": transcript_data.get("question", ""),
                        "help_clues": transcript_data.get("help_clues", "")
                    }
                    
                    transcripts.append(transcript)
            
            except json.JSONDecodeError as e:
                print(f"Error reading {file_path}: {e}")
            except Exception as e:
                print(f"Unexpected error processing {file_path}: {e}")
        
        return transcripts

# Example usage
def main():
    # Read structured transcripts
    transcripts = TranscriptVectorDB.read_structured_transcripts()
    
    # Initialize and populate vector DB
    vector_db = TranscriptVectorDB()
    vector_db.add_transcripts(transcripts)

    # Perform a search
    results = vector_db.search_transcripts("artificial intelligence")
    for result in results:
        print(f"Result: {result['metadata']['title']}")
        print(f"Excerpt: {result['document'][:100]}...")
        print(f"Distance: {result['distance']}\n")

if __name__ == "__main__":
    main()
