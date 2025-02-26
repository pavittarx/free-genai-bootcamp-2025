import os
import requests
import streamlit as st
from typing import Optional, Dict, Any
import json


# Default Model for OpenRouter
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
            
            # Extract response content
            if 'choices' in response_data and len(response_data['choices']) > 0:
                return response_data['choices'][0]['message']['content']
            elif 'error' in response_data:
                st.error(f"OpenRouter API Error: {response_data['error']}")
                return None
            else:
                st.error("Unexpected response format from OpenRouter API")
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