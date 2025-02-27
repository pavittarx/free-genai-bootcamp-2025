# Structured Data Extraction for Hindi Listening Comprehension

## Overview
The script `extract_structured_data_genai.py` uses a GenAI model to process YouTube transcripts and extract structured data with the following sections:

1. **Introduction**: Context-setting part of the transcript
2. **Dialogue**: Main conversational content
3. **Question**: Specific question asked in the transcript
4. **Help/Clues**: Hints or additional context for the question
5. **Answer**: Response or solution to the question

## Extraction Method
- Utilizes OpenRouter's Google Gemini Flash Lite model
- Sends a detailed prompt to extract structured sections
- Attempts to parse the response as JSON
- Fallback mechanism for non-JSON responses

## Key Components
- `OpenRouterChat` class from `chat.py`
- Dynamic prompt generation
- Flexible error handling

## Advantages
- Leverages AI's contextual understanding
- More nuanced extraction compared to rule-based methods
- Adaptable to various transcript structures

## Limitations
- Depends on the AI model's interpretation
- Potential for inconsistent results
- Requires careful prompt engineering

## Future Improvements
- Fine-tune prompts for specific transcript types
- Implement validation mechanisms
- Explore multiple AI model options
- Add confidence scoring for extracted sections