## Frontend for the Lang Portal Project
A hindi language learning school wants to build a sample project that could help the users learn appropriate language.
The portal in its current form, will do the following:
- Act as a launchpad for different study activities.
- Store Language words and sentences, that would aid in learning activities. 
- Act as a record keeping app for tracking progress and performance.

## Technical Requirements
- Frontend will be wriiten using Vue.
- The frontend with be communicated with the backend using HTTP requests.
- The API request and response will be in JSON format.
- The API will be stateless, and will not store any persistent data.
- The frontend is meant to be used by single user.
- No authentication/authorisation is required.

## General Instructions
- Words are not editable.
- Groups are not editable.
- Study activities are not editable.
- Sessions can be created.
- Sessions cannot be edited or deleted.
- Session Activity can be created.
- Sessions Activity cannot be edited or deleted. 

## Pages
1. Dashboard Page  
- The dashboard page will act as the landing page for the user. 
- The dashboard will contain menu. The menu will contain links to different pages.
    1. Words Page
    2. Groups Page
    3. Study Activities Page
    4. Sessions Page
    5. Settings Page
- The menu links will be contained in the left sidebar.
- The dashboard will contain a main content area.
- The main content area will contain cards. The cards will list the following items: 
    - Total Sessions
    - Total Study Activities
    - Total Groups
    - Total Words

2. Words Page
- The words page will contain a table of words.
- The table will contain the following columns: 
    - Serial Number
    - Hindi
    - Hinglish
    - English
    - Created At
- The table will be paginated and show 20 items at a time. 

2. Groups Page
- The groups page will contain a table of groups.
- The groups table will contain the following columns: 
    - Serial Number
    - Group Name
    - Created At
- No Pagination is required for Groups page. 
- Clicking a row on Groups page will navigate to the individual group page. 
- The Group page will list the Group name at the top.
- The group page will contain a table of words that are part of the group. The table will contain the following columns:
    - Serial Number
    - Hindi
    - Hinglish
    - English
    - Created At

3. Study Activities Page
- The study activities page will contain cards for each study activity. The cards will contain the following:
    - Image for the activity
    - Activity Name
    - Description
- The study activity will have a url to the activity page.
- Clicking a study activity card will launch study activity.

4. Sessions Page
- The sessions page will contain a table of sessions.
- The sessions table will contain the following columns: 
    - Serial Number
    - Session Name
    - Study Activity
    - Score
- On clicking a row in sessions table, a new page will be opened. 
- This new page contains sessions activity data.
- The page will contain the following:
    - Serial Number
    - Study Activity
    - Challenge
    - Answer
    - Input
    - Score
    - Start Time
    - End Time

# Code Conventions
1. Naming Conventions
- Use camelCase for function and variable names.
- Use kebab-case for HTML element IDs.
- Use kebab-case for class names.
- Use UPPER_SNAKE_CASE for environment variables.

2. Code Formatting
- Line length should be around 80-120 characters
- Use 4 spaces for indentation
- Use Prettier for code formatting
- Use Eslint for code linting

3. Testing
- Use Vitest for testing
- Use Playwright for e2e testing

4. Development
- The project will be written using Typescript.
- Dependency Management will be done using pnpm.
- The UI will be built using VueJs using Nuxt Framework.
- The project will use NuxtUI for components. 
- The project will use Vite bundler for Nuxt.
- The project will use tailwind css for styling.
- Use Vue CLI for project setup and development
- Use Eslint for code linting
- Use Prettier for code formatting
- Data Fetching should be done via Tanstack Query. 
- The latest versions of all packages should be used.