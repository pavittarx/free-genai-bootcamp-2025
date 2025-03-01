# Listening Learning APP

Language: Hindi
Level: 300

## Business Goal
You have been tasked with building a language learning comprehension app as an Applied AI Engineer. 
The goal is to build exam-style listening scenarios using code. 
This would include everyday conversations in target language, such as ordering food, asking for directions making plans.

## Technical Uncertainty
- Will my system be able to support the infrastructure required to run this project? 
- Can I leverage free-tier cloud services to run the app?
- Will appropriate data be available for training the model?
- I do not know how to work with vector databases.
- There are no official exams for Hindi Language as JLPT5, there are exams like CSAT, and school exams, but they do not follow the pattern as JLPT5.
- Structured Language Comprehension Videos are not available for Hindi Language as is available for Japanese
- ASR / TTS support might not be available for hindi language.
- Videos might be missing transctipts needed.
- The technical expertise with technologies used is missing. 

## Technical Requirements 
- Chat w/ LLM Model
- Generate and Save Transcripts
- Generate structured data from transcripts
- Vector DB (ChromaDB)
- OpenRouter API, Google Flash Model
- ASR / TTS (Optional)


## Technical Restrictions
- Youtube Video url should be provied to pull the transcript
- OpenRouter will be used as chat agent
- The app should be able to handle a conversation with the user.

## Transcript Sources 
- https://www.youtube.com/watch?v=ijOqyASjgoM
- https://www.youtube.com/watch?v=KXXp0fhHZys

- The structure is as follows
    1. Introduction
    2. Question
    3. Dialogue
    4. Options for answer
    5. Answer
