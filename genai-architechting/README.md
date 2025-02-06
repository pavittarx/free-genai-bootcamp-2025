# GENAI Architechting

### Diagram
![](./GenAI-Arch.png?raw=true)

### Use Case / Business Goals
Lingo is a language learning platform that people to subscribe and learn a foreign language using a language instructor. Lingo offers prerecorded courses as well as live classes with the instructors. 
The company plans to expand and allow for interactive learning using GenAI tools.

### Functional Requirements
The company wants to leverage cloud services such as AWS. They are concerned about the feature adoption and upfront costs required to set up their own GenAI system. Hence, they plan use cloud GenAI services in combination of their existing systems. 

- The company has 200 active learners, and nearly 100 daily repeat learners. 
- The company wants to pilot run this project for 90 days under a budget of 3000 dollars. 

### Non-Functional Requirements
- The company plans to utilise existing on-premise user database for integrating into the GenAI workflow. 
- The company plans to filter out any user identifiable information using during RAG step and input guardrails. 
- The system must be able to serve a load of around 300 users at any given time. 
- The system will be using caches to resolve repeated queries / prompts.  
- Billing alerts should be set to monitor cloud costs and optimise accordingly.

### Assumptions
- The pilot project could be hosted on cloud for 90 days under a budget of 3000 dollars.
- The project could be replicated on-premise once the pilot project is successfull.
- The learners wants an interactive guide to help with their language learning.
- The GenAI system can adopt to the goals and language level of the learner.  

### Risks
- Cloud usage bills could negatively hamper the success of the project. 
- User identifiable and personal information should not be exposed to AI models.
- Incorrect responses could lead to learner dissatisfaction, wrong learning outcomes, and completely ignoring the tool.

### Data Strategy
- The system should use either self developed or purchased materials, no material should be used without appropriate permissions.
- Data should include all language levels of materials.
- Personal / user-identifiable information should not be used or filtered out.
- The model should use existing database for data needs. New data could be stored and provided as needed.  

### Model Selection
- SaaS Model
- Text-Text I/O
- One model 
- One call per prompt