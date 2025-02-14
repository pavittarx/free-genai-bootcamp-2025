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
- Study Activities need to be their own interactive mini apps.
- Study Activities should launch as a popup on the current page.
- The popup should come as an overlay centered on the screen, 
- The popup should close if the user navigates to other pages
- The scoring mechanism is provided with the activity.

- Study Activities will create a session on start
- A session consists of a single run of the activity,
- A session activity is the challenge user answers within the acitivity. 
- Sessions and Session activities are saved via the API.

- Each challenge is scored based on the input that is provided.
- The score are saved via the means of session activity API. 
- There should be option to skip the challenge and move to the next challenge.
- Activity should end after final challenge.
- Final score should be displayed at the end of the activity.

- Session ends at the end of study activity.
- Final Score of the session is prepared based on all the scores in the session. 

- Last challenge is marked as end of study activity, if the activity is closed. 
- A new session is created each time the user launches the study activity. 
- The study activity session should be interactive, mimicing a game like enviornment and UI.
- On challenge completion, the user is back on study activity page.

## API Integration and Endpoints
- Random words are fetched from `/api/words/random`
- Session is started with a POST request to the API `/api/session`
  - include activity_id in the body
- Each challenge in the study activity is saved using POST request on `/api/session-activity`
  - the session activity has following paramters
  - current challenge word passed as a string
  - current challenge answer passed as a string
  - what did the user input as a string
  - the score of the user on this attempt
- The challenge ends with PUT request on `/api/session` 
  - the body provides the overall final score of the challenge
- All the session and session activity related data can be found via GET request on `/api/sessions/:id` with session id passes as id.