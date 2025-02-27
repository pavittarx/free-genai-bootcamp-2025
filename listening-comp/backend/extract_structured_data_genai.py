import os
import json
from typing import Dict, List
from chat import OpenRouterChat

def load_transcript(file_path: str) -> List[Dict]:
    """Load transcript from a JSON file."""
    with open(file_path, 'r', encoding='utf-8') as f:
        return json.load(f)

def extract_structured_data_with_genai(transcript_text: str) -> Dict[str, str]:
    """
    Use GenAI to extract structured data from the transcript.
    
    Structured sections:
    1. Introduction
    2. Dialogue
    3. Question
    4. Help/Clues for question
    5. Answer
    """
    # Initialize OpenRouter chat client
    chat_client = OpenRouterChat()
    
    # Prompt for structured data extraction
    prompt = f"""
    Carefully analyze the following Hindi transcript and extract ACTUAL, SPECIFIC sections:

    CRITICAL INSTRUCTIONS:
    - Analyse the script and fix it before attempting extraction
    - DO NOT use placeholders like "..." or "…"
    - Adjust the conversation text based on the transcript
    - Fix Dialogue, Question, Help/Clues, and Answer sections
    - The introduction should be translated to hindi
    - The flow of the conversation should be in natural hindi language
    - Format the conversations properly, using proper punctuation

    Section Definitions:
    1. Introduction: The FIRST, MOST INITIAL context-setting phrase (max 15 words)
    2. Dialogue: The MAIN conversational part with KEY exchanges
    3. Question: The EXACT question phrase from the transcript
    4. Help/Clues: Background context (use null if none)
    5. Answer: The DIRECT, SPECIFIC response to the question

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
    - NO GENERIC or PLACEHOLDER text
    - Use ONLY words from the given transcript
    - Be PRECISE and CONCRETE
    """
    
    # Generate response
    response = chat_client.generate_response(prompt)
    
    # Try to parse the response as JSON
    try:
        # Remove code block markers if present
        if response and response.startswith('```json') and response.endswith('```'):
            response = response[7:-3].strip()
        
        # Attempt to parse the response
        structured_data = json.loads(response)
        
        # Strict validation of structure
        required_keys = ["introduction", "dialogue", "question", "help_clues", "answer"]
        
        # Function to clean and validate each section
        def clean_section(section):
            # Remove placeholder markers
            placeholders = ['...', '…', '""', "''"]
            for placeholder in placeholders:
                section = section.replace(placeholder, '').strip()
            return section
        
        # Validate and clean each key
        for key in required_keys:
            # Ensure key exists
            if key not in structured_data:
                raise ValueError(f"Missing required key: {key}")
            
            # If not help_clues, ensure it's a string and clean
            if key != "help_clues":
                if not isinstance(structured_data[key], str):
                    raise ValueError(f"Key {key} must be a string")
                
                # Clean the section
                structured_data[key] = clean_section(structured_data[key])
                
                # If section is empty after cleaning, use full transcript
                if not structured_data[key]:
                    if key == "introduction":
                        # Take first 50 characters for introduction
                        structured_data[key] = transcript_text[:50]
                    elif key == "dialogue":
                        structured_data[key] = transcript_text
                    elif key == "question":
                        # Try to find a question in the transcript
                        import re
                        question_match = re.search(r'\b(क्या|कैसे|क्यों|कब)\b', transcript_text)
                        structured_data[key] = question_match.group(0) if question_match else transcript_text
                    elif key == "answer":
                        structured_data[key] = transcript_text
        
        # Additional checks
        if len(structured_data["introduction"]) > 100:
            # If introduction is too long, truncate
            structured_data["introduction"] = structured_data["introduction"][:100]
        
        return structured_data
    
    except (json.JSONDecodeError, TypeError, ValueError) as e:
        # Fallback to manual parsing if JSON parsing fails
        print(f"Failed to parse GenAI response as JSON: {e}")
        print(f"Raw response: {response}")
        
        # Attempt to extract data manually if possible
        try:
            # Try to find JSON-like content within the response
            import re
            json_match = re.search(r'\{.*\}', response, re.DOTALL)
            if json_match:
                manual_parse = json.loads(json_match.group(0))
                
                # Apply same validation as above
                for key in required_keys:
                    if key not in manual_parse:
                        manual_parse[key] = "" if key != "help_clues" else None
                    
                    if key != "help_clues" and not isinstance(manual_parse[key], str):
                        manual_parse[key] = str(manual_parse[key]).strip()
                
                return manual_parse
        except Exception:
            pass
        
        # If all parsing fails, return default structure based on transcript
        return {
            "introduction": transcript_text[:50],
            "dialogue": transcript_text,
            "question": re.search(r'\b(क्या|कैसे|क्यों|कब)\b', transcript_text).group(0) if re.search(r'\b(क्या|कैसे|क्यों|कब)\b', transcript_text) else transcript_text,
            "help_clues": None,
            "answer": transcript_text
        }

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
            structured_data = extract_structured_data_with_genai(full_text)
            
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
