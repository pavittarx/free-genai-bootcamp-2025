import os
import json
import re
from typing import Dict, List, Any
from backend.chat import OpenRouterChat

def load_transcript(file_path: str) -> List[Dict]:
    """Load transcript from a JSON file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        return json.load(f)

def clean_section(section: str) -> str:
    """
    Clean and normalize a text section.
    
    Args:
        section (str): Input text section
    
    Returns:
        str: Cleaned and normalized text
    """
    # Remove code block markers and extra whitespace
    section = section.replace('```json', '').replace('```', '').strip()
    
    # Remove placeholder markers
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
                "help_clues": None,
                "answer": response
            }

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
    - Introduction should include the setting, enviornment or scenary of the what follows.
    - The flow of the conversation should be in natural hindi language
    - Use proper punctuation to format the sections.

    Section Definitions:
    1. Introduction: The FIRST, MOST INITIAL context-setting phrase (max 15-20 words)
    2. Dialogue: The conversational part of the transcript, on which the question is asked.
    3. Question: The question or questions being asked in the transcript.
    4. Help/Clues: Any clues if provided to guess the answer.
    5. Answer: The answer to the questions being asked.

    MANDATORY FORMAT:
    {{
        "introduction": "first words",
        "dialogue": "conversation text",
        "question": "question from transcript",
        "help_clues": null or "SPECIFIC context",
        "answer": "answer"
    }}

    Transcript Text:
    {transcript_text}

    FINAL WARNING: 
    - Do not use Generic or placeholder text
    - Try to provide as much information as possible
    - Do not use placeholders like "..." or "..."
    - Keep the response format to exactly match what is described.
    - Be PRECISE and CONCRETE
    """
    
    # Generate response
    response = chat_client.generate_response(prompt)
    
    # Parse the response
    structured_data = parse_json_response(response)
    
    # Clean and validate each section
    required_keys = ["introduction", "dialogue", "question", "help_clues", "answer"]
    
    for key in required_keys:
        # Ensure key exists
        if key not in structured_data:
            structured_data[key] = "" if key != "help_clues" else None
        
        # Clean non-help_clues sections
        if key != "help_clues":
            structured_data[key] = clean_section(str(structured_data[key]))
    
    # Truncate introduction if too long
    if len(structured_data["introduction"]) > 100:
        structured_data["introduction"] = structured_data["introduction"][:100]
    
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
