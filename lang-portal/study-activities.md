## Study Activities
1. Unscramble the words
You are given a list of words, which may or may not be in order. Reorder the words in any order you like, but ensure that no word is left out.
- The answers can be input out in both hindi / hinglish.
- Score for each correct answer is 5 points.
- Score for each incorrect answer is 0 points.
- Score for each empty answer is 0 points.
- The words can be picked up from /api/words/random. It has scrambled letters.
- Each session will show 10 words.
- The challenge should show unscrabled hindi words fetched from /api/words/random.

2. Group the Words
You are given a list of words, which belong to the same group. Group the words in any order you like, but ensure that no word is left out.
- Total of 6 words are displayed at a time. 
- Maximum of 4 words are correct. 
- The answers can be input out in both hindi / hinglish.
- Score for each correct answer is 10 points.
- Score for each incorrect answer is 0 points.
- Score for each empty answer is 0 points.

3. Complete the Word
You are given a word, which is not displayed correctly. You need to guess the complete word. 
- The answers can be input out in both hindi / hinglish.
- Score for each correct answer is 10 points.
- Score for each incorrect answer is 0 points.
- Score for each empty answer is 0 points.
- Total possible attemps are 3 for each entry. 

# Study Activities Execution Plan

## Overview
This plan outlines the comprehensive strategy for developing and implementing study activities in our language learning portal. The goal is to create engaging, educational, and interactive learning experiences that help users improve their language skills.

## Activity Types

### 1. Unscramble Words
- **Objective**: Improve vocabulary and word recognition
- **Mechanics**:
  - Present scrambled letters of a word
  - User must rearrange letters to form a correct word
- **Difficulty Levels**:
  - Easy: 3-4 letter words
  - Medium: 5-6 letter words
  - Hard: 7+ letter words
- **Scoring**:
  - Correct word: +10 points
  - Time bonus: Additional points for quick solving
  - Hint penalty: Points deducted for using hints

### 2. Word Match
- **Objective**: Enhance vocabulary and word-meaning associations
- **Mechanics**:
  - Present words in two columns (Hindi/English)
  - User must match corresponding words
- **Difficulty Levels**:
  - Easy: Common, basic vocabulary
  - Medium: Intermediate vocabulary
  - Hard: Advanced, context-specific words
- **Scoring**:
  - Correct match: +5 points per pair
  - Perfect match bonus: +20 points
  - Time-based scoring

### 3. Sentence Constructor
- **Objective**: Improve grammar and sentence formation skills
- **Mechanics**:
  - Provide a set of words
  - User must construct a grammatically correct sentence
  - Support for multiple languages (Hindi, English)
- **Difficulty Levels**:
  - Easy: Simple subject-verb-object sentences
  - Medium: Sentences with adjectives and adverbs
  - Hard: Complex sentences with multiple clauses
- **Scoring**:
  - Grammatically correct sentence: +15 points
  - Contextually appropriate sentence: +10 bonus points
  - Complexity bonus

## Implementation Guidelines for Study Activities

### Core Principles
- Study Activities are interactive mini-apps
- Launch as a popup on the current page
- Create a session on start
- End the session on completion
- Score based on answers given

### Popup Implementation Requirements
- Full-screen or modal overlay
- Smooth transition animations
- Responsive design
- Close button for exit
- Session tracking
- Progress preservation

### Session Management
- Generate unique session ID on activity start
- Track:
  - Start time
  - End time
  - Challenges attempted
  - Scores
  - Completion status

### Technical Constraints
- Minimal performance overhead
- Quick load times
- Seamless user experience
- Consistent UI/UX across activities

### Interaction Flow
1. User selects activity from study activities page
2. Popup launches immediately
3. Activity starts new session
4. User completes challenges
5. Popup closes, returning to study activities
6. Session data saved and processed

### Error Handling
- Graceful popup closure
- Session data preservation
- Clear error messaging
- Ability to resume or restart

## Technical Implementation

### Frontend Development
- [ ] Create individual activity pages
- [ ] Implement dynamic routing for activities
- [ ] Design responsive and interactive UI
- [ ] Create reusable activity components
- [ ] Implement state management for activities

### Backend Development
- [ ] Design activity data models
- [ ] Create API endpoints for:
  - Fetching activity details
  - Submitting activity results
  - Tracking user progress
- [ ] Implement activity generation logic
- [ ] Create scoring and progression system

### Data Management
- [ ] Store activity templates
- [ ] Track user performance per activity
- [ ] Generate adaptive difficulty levels
- [ ] Implement spaced repetition algorithms

## User Experience Considerations
- Engaging and intuitive interface
- Clear instructions for each activity
- Immediate feedback on performance
- Progressive difficulty scaling
- Motivational elements (badges, streaks)

## Performance Metrics
- Time to complete activity
- Accuracy rate
- Progression through difficulty levels
- Retention of learned words/concepts

## Future Expansions
- AI-powered personalized learning paths
- Multi-language support
- Social learning features
- Integration with external language resources

## Technical Stack
- Frontend: Vue.js, Nuxt.js
- State Management: Tanstack Query
- Styling: Tailwind CSS

Guides to Build Study Activities
- There should be option to skip the challenge and move to the next challenge.
- Study Activities need to be their own interactive mini apps.
- Study Activities should launch as a popup on the current page.
- Study Activities will create a session on start
- Study Activities will end the session on end
- Study Activities will have a score based on the answers given.
- The scoring mechanism is provided in the prompt.
- The score to be given on correct answer is also stored in the database.
- All the answers and challenges are stored in the database.
- The session for the study activity is stored in the session table.
- The challenged is stored in session_activity table.
- If any study activity session is closed in between, last answer's time should be taken as the end time of the session.
- If any study activity session is closed in between, a new session should be created for the next challenge.
- The study activity session should be interactive, mimicing a game like enviornment and UI.
- On challenge completion, the user is back on study activity page.
- The popup should come as an overlay centered on the screen, 
- The popup should close if the user navigates to other pages