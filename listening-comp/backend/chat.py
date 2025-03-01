import os
import requests
import json
import streamlit as st
from typing import Optional, Dict, Any
from dotenv import load_dotenv

load_dotenv()

MODEL_ID = "google/gemini-2.0-flash-lite-001"

class OpenRouterChat:
    def __init__(self, model_id: str = MODEL_ID):
        """Initialize OpenRouter chat client"""
        # Retrieve API key from environment variable
        self.api_key = os.getenv('OPENROUTER_API_KEY')
        if not self.api_key:
            raise ValueError("OpenRouter API key not found. Please set OPENROUTER_API_KEY environment variable.")
        
        # OpenRouter API endpoint
        self.api_url = "https://openrouter.ai/api/v1/chat/completions"
        self.model_id = model_id

    def generate_response(self, 
                          message: str, 
                          inference_config: Optional[Dict[str, Any]] = None) -> Optional[str]:
        """Generate a response using OpenRouter"""
        if inference_config is None:
            inference_config = {"temperature": 0.7}
        
        # Prepare headers and payload
        headers = {
            "Authorization": f"Bearer {self.api_key}",
            "Content-Type": "application/json"
        }

        payload = {
            "model": self.model_id,
            "messages": [
                {"role": "system", "content": "You are a helpful AI assistant focused on Hindi language learning."},
                {"role": "user", "content": message}
            ],
            "temperature": inference_config.get('temperature', 0.7)
        }

        try:
            # Send request to OpenRouter
            response = requests.post(
                self.api_url, 
                headers=headers, 
                json=payload
            )
            
            # Check for successful response
            response.raise_for_status()
            
            # Parse response
            response_data = response.json()
            
            # Debug: Print full response for investigation
            print("Full OpenRouter API Response:")
            print(json.dumps(response_data, indent=2))
            
            # Enhanced error handling for JSON parsing
            try:
                # Check if the response is a valid JSON with expected structure
                if not isinstance(response_data, dict):
                    st.error(f"Unexpected response type: {type(response_data)}")
                    return None

                # Extract response content with more robust checks
                if 'choices' in response_data and response_data['choices']:
                    first_choice = response_data['choices'][0]
                    
                    # Check if 'message' and 'content' exist
                    if 'message' in first_choice and 'content' in first_choice['message']:
                        content = first_choice['message']['content']
                        
                        # Additional check to ensure content is not empty or None
                        if content and isinstance(content, str):
                            # Try to parse content if it looks like it might be a JSON string
                            if content.strip().startswith('{') and content.strip().endswith('}'):
                                try:
                                    parsed_content = json.loads(content)
                                    return json.dumps(parsed_content)  # Return as formatted JSON
                                except json.JSONDecodeError:
                                    # If JSON parsing fails, return original content
                                    return content
                            return content
                        
                    st.error("Response content is empty or invalid")
                    return None

                elif 'error' in response_data:
                    st.error(f"OpenRouter API Error: {response_data['error']}")
                    return None
                else:
                    st.error("Unexpected response format from OpenRouter API")
                    # Log the full response for debugging
                    print("Unexpected response structure:")
                    print(json.dumps(response_data, indent=2))
                    return None

            except Exception as parsing_error:
                st.error(f"Error parsing OpenRouter response: {str(parsing_error)}")
                # Log the full response and error for debugging
                print("Response that caused parsing error:")
                print(json.dumps(response_data, indent=2))
                print(f"Parsing error details: {str(parsing_error)}")
                return None
            
        except requests.RequestException as e:
            st.error(f"Request Error: {str(e)}")
            # If request fails, print full error details
            if hasattr(e, 'response'):
                st.error(f"Response content: {e.response.text}")
            return None
        except Exception as e:
            st.error(f"Unexpected Error: {str(e)}")
            return None


if __name__ == "__main__":
    chat = OpenRouterChat()
    while True:
        user_input = input("You: ")
        if user_input.lower() == '/exit':
            break
        response = chat.generate_response(user_input)
        print("Bot:", response)