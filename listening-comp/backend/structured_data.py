import os
import json
import re
import random
from typing import Dict, List, Any
from backend.chat import OpenRouterChat

def load_transcript(file_path: str) -> List[Dict]:
    """Load transcript from a JSON file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        return json.load(f)

def clean_section(section: str) -> str:
    section = section.replace('```json', '').replace('```', '').strip()
    
    placeholders = ['...', '…', '""', "''"]
    for placeholder in placeholders:
        section = section.replace(placeholder, '').strip()
    
    return section

def parse_json_response(response: str) -> Dict[str, Any]:
    # Try direct JSON parsing first
    try:
        # Remove any leading/trailing whitespace
        response = response.strip()
        
        # Try parsing as-is
        return json.loads(response)
    except json.JSONDecodeError:
        # Try extracting JSON from code block
        json_match = re.search(r'\{.*\}', response, re.DOTALL | re.MULTILINE)
        
        if json_match:
            try:
                return json.loads(json_match.group(0))
            except json.JSONDecodeError:
                pass
        
        # Last resort: manual parsing
        try:
            # Remove code block markers and try parsing
            cleaned_response = clean_section(response)
            return json.loads(cleaned_response)
        except Exception:
            # If all parsing fails, return a default structure
            return {
                "introduction": response[:100],
                "dialogue": response,
                "question": "",
                "options": [],
                "answer": response
            }

def generate_multiple_options(answer: str) -> List[str]:
    """
    Generate multiple options based on the correct answer
    
    Args:
        answer (str): Correct answer
    
    Returns:
        List of multiple choice options
    """
    chat_client = OpenRouterChat()
    
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
        response = chat_client.generate_response(prompt)
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

def structured_data_with_genai(transcript_text: str) -> Dict[str, str]:
    # Initialize OpenRouter chat client
    chat_client = OpenRouterChat()
    
    # Prompt for structured data extraction
    prompt = f"""
    Carefully analyze the following Hindi transcript and extract ACTUAL, SPECIFIC sections:

    CRITICAL INSTRUCTIONS:
    - Analyse the script and fix it before attempting extraction
    - DO NOT use placeholders like "..." or "…"
    - Adjust the conversation text based on the transcript
    - Introduction should include the setting, environment or scenario of what follows.
    - The flow of the conversation should be in natural hindi language
    - Use proper punctuation to format the sections.

    Section Definitions:
    1. Introduction: The FIRST, MOST INITIAL context-setting phrase (max 15-20 words)
    2. Dialogue: The conversational part of the transcript, on which the question is asked.
    3. Question: The question or questions being asked in the transcript.
    4. Multiple Options: Generate up to four options for the question
    5. Answer: The correct answer to the questions being asked.

    MANDATORY FORMAT:
    {{
        "introduction": "first words",
        "dialogue": "conversation text",
        "question": "question from transcript",
        "options": ["option1", "option2", "option3", "option4"],
        "answer": "correct option"
    }}

    Transcript Text:
    {transcript_text}

    FINAL WARNING: 
    - Do not use Generic or placeholder text
    """
    
    # Generate structured response
    response = chat_client.generate_response(prompt)
    
    # Parse the JSON response
    structured_data = parse_json_response(response)
    
    # Ensure options are generated if not present
    if not structured_data.get('options') or len(structured_data['options']) < 4:
        structured_data['options'] = generate_multiple_options(
            structured_data.get('answer', '')
        )
    
    return structured_data

def process_transcripts(transcript_dir: str, output_dir: str):
    """Process all transcripts in the given directory."""
    # Ensure output directory exists
    os.makedirs(output_dir, exist_ok=True)
    
    # Process each transcript
    for filename in os.listdir(transcript_dir):
        if filename.endswith('.json'):
            file_path = os.path.join(transcript_dir, filename)
            
            # Load transcript
            transcript = load_transcript(file_path)
            
            # Convert transcript to full text
            full_text = " ".join([entry['text'] for entry in transcript])
            
            # Extract structured data using GenAI
            structured_data = structured_data_with_genai(full_text)
            
            # Create output filename
            output_filename = os.path.splitext(filename)[0] + '_structured.json'
            output_path = os.path.join(output_dir, output_filename)

            # Save structured data
            with open(output_path, 'w', encoding='utf-8') as f:
                json.dump(structured_data, f, ensure_ascii=False, indent=2)
            
            print(f"Processed {filename} -> {output_filename}")

def main():
    # Directories
    transcript_dir = os.path.join(os.path.dirname(__file__), 'transcripts')
    output_dir = os.path.join(os.path.dirname(__file__), 'structured_transcripts')
    
    # Process transcripts
    process_transcripts(transcript_dir, output_dir)

if __name__ == "__main__":
    main()
