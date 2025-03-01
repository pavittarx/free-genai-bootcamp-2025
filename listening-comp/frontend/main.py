import streamlit as st
import sys
import os
import json
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from backend.rag import RAGAssistant
from backend.transcript import YTTranscriptDownloader
from backend.chat import OpenRouterChat
from backend.structured_data import structured_data_with_genai
from backend.vector_db import TranscriptVectorDB

# Page configuration
st.set_page_config(
    page_title="Hindi Learning Assistant",
    page_icon="🛕",
    layout="wide"
)

class HindiLearningApp:
    def __init__(self):
        """
        Initialize the Hindi Learning Application
        """
        self.rag_assistant = RAGAssistant()
        
        # Initialize session state for exercise
        if 'current_exercise' not in st.session_state:
            st.session_state.current_exercise = None
        
        # Initialize audio playing state
        if 'audio_playing' not in st.session_state:
            st.session_state.audio_playing = False
        
        # Initialize additional states for new features
        if 'transcript' not in st.session_state:
            st.session_state.transcript = None
        if 'messages' not in st.session_state:
            st.session_state.messages = []
        if 'structured_data' not in st.session_state:
            st.session_state.structured_data = None
        if 'chat' not in st.session_state:
            st.session_state.chat = OpenRouterChat()
    
    def _render_audio_button(
        self, 
        audio_path: str, 
        label: str
    ):
        """
        Render an audio button with play/pause functionality
        
        Args:
            audio_path (str): Path to the audio file
            label (str): Button label
        """
        if not audio_path or not os.path.exists(audio_path):
            return
        
        # Determine current state
        is_playing = st.session_state.audio_playing
        
        # Audio button
        if st.button(
            f"{'🔊 Pause' if is_playing else '🔈 Play'} {label}", 
            key="audio_exercise_button"
        ):
            # Toggle audio state
            st.session_state.audio_playing = not is_playing
            
            # Audio playback logic
            st.components.v1.html(f"""
            <script>
                var audio = document.getElementById('exercise_audio');
                if (audio) {{
                    if ({str(not is_playing).lower()}) {{
                        audio.play();
                    }} else {{
                        audio.pause();
                    }}
                }}
            </script>
            <audio id="exercise_audio" src="{audio_path}"></audio>
            """, height=0)
    
    def generate_exercise(self, topic=None, difficulty='मध्यम'):
        """
        Generate a new language learning exercise
        
        Args:
            topic (str, optional): Specific topic for the exercise
            difficulty (str, optional): Difficulty level
        """
        # Generate exercise
        exercise = self.rag_assistant.generate_learning_exercise(
            topic=topic, 
            difficulty=difficulty
        )
        
        # Store in session state
        st.session_state.current_exercise = exercise
        
        # Reset audio playing state
        st.session_state.audio_playing = False
    
    def render_exercise(self):
        """
        Render the current exercise with audio support
        """
        from backend.audio_generator import AudioGenerator
        import os

        if not st.session_state.current_exercise:
            st.info("Click 'Generate Exercise' to start learning!")
            return
        
        exercise = st.session_state.current_exercise
        
        # Exercise Introduction Section
        st.subheader("Exercise Introduction")
        st.write(exercise.get('introduction', ''))
        
        # Dialogue Section
        st.subheader("Dialogue")
        st.write(exercise.get('dialogue', ''))
        
        # Question Section
        st.subheader("Question")
        st.write(exercise.get('question', ''))
        
        # Audio Generation Button
        if st.button("🔊 Generate Audio"):
            # Generate single audio file for the entire exercise
            audio_path = AudioGenerator.generate_audio(exercise)
            
            # Store audio path in session state
            st.session_state.exercise_audio_path = audio_path
        
        # Display Audio Player if audio is generated
        if hasattr(st.session_state, 'exercise_audio_path') and st.session_state.exercise_audio_path:
            # Check if file exists before playing
            if os.path.exists(st.session_state.exercise_audio_path):
                with open(st.session_state.exercise_audio_path, 'rb') as audio_file:
                    st.audio(audio_file, format='audio/mp3')
            else:
                st.error("Audio file not found.")
        
        # Options Section
        st.subheader("Options")
        options = exercise.get('options', [])
        
        with st.form("exercise_submission"):
            selected_option = st.radio(
                "Select your answer:", 
                options, 
                index=None
            )
            
            submit_button = st.form_submit_button("Check Answer")
        
        # Handle answer submission
        if submit_button and selected_option:
            correct_answer = exercise.get('answer', '')
            if selected_option == correct_answer:
                st.success("🎉 सही उत्तर! (Correct Answer!)")
            else:
                st.error(f"क्षमा करें, यह सही उत्तर नहीं है। सही उत्तर है: {correct_answer}")
    
    def render_chat_stage(self):
        """Render an improved chat interface"""
        st.header("Chat with Lexia")

        # Introduction text
        st.markdown("""
        Start by exploring Nova's base Hindi language capabilities. Try asking questions about Hindi grammar, 
        vocabulary, or cultural aspects.
        """)

        # Display chat messages
        for message in st.session_state.messages:
            with st.chat_message(message["role"], avatar="🧑‍💻" if message["role"] == "user" else "🤖"):
                st.markdown(message["content"])

        # Chat input area
        if prompt := st.chat_input("Ask about Hindi language..."):
            # Process the user input
            self.process_message(prompt)

        # Example questions in sidebar
        with st.sidebar:
            st.markdown("### Try These Examples")
            example_questions = [
                "How do I say 'Which station is this?' in Hindi?",
                "Explain the difference between namaste and shubh?",
                "What's the polite form of tu?",
                "How do I count in Hindi?",
                "What is dusk and dawn in hindi?",
                "How do I ask for directions politely?"
            ]
            
            for q in example_questions:
                if st.button(q, use_container_width=True, type="secondary"):
                    # Process the example question
                    self.process_message(q)
                    st.rerun()

        # Add a clear chat button
        if st.session_state.messages:
            if st.button("Clear Chat", type="primary"):
                st.session_state.messages = []
                st.rerun()

    def process_message(self, message: str):
        """Process a message and generate a response"""
        # Add user message to state and display
        st.session_state.messages.append({"role": "user", "content": message})
        with st.chat_message("user", avatar="🧑‍💻"):
            st.markdown(message)

        # Generate and display assistant's response
        with st.chat_message("assistant", avatar="🤖"):
            response = st.session_state.chat.generate_response(message)
            if response:
                st.markdown(response)
                st.session_state.messages.append({"role": "assistant", "content": response})

    def render_transcript_processing_stage(self):
        """
        Unified stage for transcript download, processing, and structured data generation
        with a clear, chronological workflow
        """
        st.header("Transcript Processing & Learning")
        
        # Create steps for a clear workflow
        steps = [
            "1. Enter YouTube URL",
            "2. Download Transcript",
            "3. View Transcript Details",
            "4. Generate Learning Exercise"
        ]
        
        # Progress tracking
        current_step = 0
        
        # URL input - Step 1
        st.markdown("## 🔗 Enter YouTube URL")
        url = st.text_input(
            "YouTube URL", 
            placeholder="Paste a Hindi lesson YouTube URL",
            key="transcript_url_input"
        )
        
        # Ensure URL is valid before proceeding
        if not url:
            st.info("Please enter a valid YouTube URL to begin.")
            return
        
        # Transcript Download - Step 2
        st.markdown("## 📥 Download Transcript")
        if st.button("Download Transcript", key="download_transcript_btn"):
            try:
                # Download transcript
                downloader = YTTranscriptDownloader()
                transcript = downloader.get_transcript(url)
                
                if transcript:
                    # Store the raw transcript text in session state
                    transcript_text = "\n".join([entry['text'] for entry in transcript])
                    st.session_state.transcript = transcript_text
                    st.success("Transcript downloaded successfully!")
                else:
                    st.error("No transcript found for this video.")
            except Exception as e:
                st.error(f"Error downloading transcript: {str(e)}")
        
        # Check if transcript exists
        if not st.session_state.get('transcript'):
            st.warning("Download the transcript before proceeding.")
            return
        
        # Transcript Details - Step 3
        st.markdown("## 📄 Transcript Analysis")
        
        # Character and Language Analysis
        hi_chars, total_chars = self.count_characters(st.session_state.transcript)
        
        col1, col2, col3 = st.columns(3)
        
        with col1:
            st.metric("Total Characters", total_chars)
        
        with col2:
            st.metric("Hindi Characters", hi_chars)
        
        with col3:
            st.metric("Hindi Character %", 
                      f"{(hi_chars/total_chars * 100):.2f}%" if total_chars > 0 else "N/A")
        
        # Transcript Display
        with st.expander("View Full Transcript"):
            st.text_area(
                label="Raw Transcript", 
                value=st.session_state.transcript, 
                height=300,
                disabled=True
            )
        
        # Structured Data Generation - Step 4
        st.markdown("## 🧩 Generate Learning Exercise")
        
        if st.button("Create Learning Exercise", key="generate_structured_data_btn"):
            try:
                # Preprocess transcript: remove newlines, normalize spaces
                processed_transcript = ' '.join(st.session_state.transcript.split())
                
                # Extract structured data directly from the current transcript
                structured_data = structured_data_with_genai(processed_transcript)
                st.session_state.structured_data = structured_data
                
                # Display structured data details
                self.render_structured_data_details(structured_data)
                
            except Exception as e:
                st.error(f"Error processing structured data: {e}")
        
        # If structured data already exists, show it
        elif st.session_state.get('structured_data'):
            st.markdown("## 🎓 Previous Learning Exercise")
            self.render_structured_data_details(st.session_state.structured_data)

    def render_structured_data_details(self, structured_data):
        """
        Render detailed view of structured data with save and vector store options
        
        Args:
            structured_data (dict): Processed structured data from transcript
        """
        col1, col2 = st.columns(2)
        
        with col1:
            st.subheader("Learning Scenario")
            
            # Context/Introduction
            st.markdown("#### 📝 Context")
            st.info(structured_data.get('introduction', 'No context available'))
            
            # Dialogue
            st.markdown("#### 💬 Dialogue")
            st.write(structured_data.get('dialogue', 'No dialogue available'))
            
            # Question
            st.markdown("#### ❓ Learning Question")
            st.warning(structured_data.get('question', 'No question generated'))
        
        with col2:
            st.subheader("Interactive Exercise")
            
            # Multiple Choice Options
            options = structured_data.get('options', [])
            answer = structured_data.get('answer', '')
            
            # Ensure we have options
            if not options:
                st.warning("No multiple-choice options available.")
                return
            
            # Radio button selection for multiple choice
            user_selection = st.radio(
                "Select the correct answer:", 
                options=options
            )
            
            # Update selected answer in session state
            st.session_state.selected_answer = user_selection
            
            # Check answer
            if st.button("Submit Answer"):
                if user_selection == answer:
                    st.success("🎉 Correct! Great job understanding the dialogue.")
                else:
                    st.error(f"Incorrect. The correct answer is: {answer}")
    
        # Action Buttons
        col1, col2 = st.columns(2)
        
        with col1:
            # Save Structured Transcript
            if st.button("💾 Save Structured Transcript"):
                # Generate unique filename
                import uuid
                filename = f"transcript_{uuid.uuid4().hex[:8]}_structured.json"
                
                # Path to structured transcripts directory
                output_dir = os.path.join(
                    os.path.dirname(__file__), 
                    '..', 
                    'backend', 
                    'structured_transcripts'
                )
                
                # Ensure directory exists
                os.makedirs(output_dir, exist_ok=True)
                
                # Full path for the new file
                output_path = os.path.join(output_dir, filename)
                
                # Save structured data
                with open(output_path, 'w', encoding='utf-8') as f:
                    json.dump(structured_data, f, ensure_ascii=False, indent=2)
                
                st.success(f"Structured transcript saved as {filename}")
    
        with col2:
            # Build Vector Store
            if st.button("🗃️ Build Vector Store"):
                try:
                    from backend.vector_db import TranscriptVectorDB
                    
                    # Read structured transcripts
                    transcripts = TranscriptVectorDB.read_structured_transcripts()
                    
                    # Check if any transcripts were found
                    if not transcripts:
                        st.warning("No structured transcripts found. Save some transcripts first.")
                        return
                    
                    # Initialize and populate vector DB
                    vector_db = TranscriptVectorDB()
                    
                    # Optional: Reset collection before adding new transcripts
                    vector_db.reset_collection()
                    
                    # Add transcripts
                    vector_db.add_transcripts(transcripts)
                    
                    st.success(f"Vector store built with {len(transcripts)} transcripts!")
                    
                    # Optional: Display some sample search results
                    with st.expander("🔍 Sample Vector Store Search"):
                        # Perform a sample search
                        results = vector_db.search_transcripts("language learning")
                        
                        if results:
                            for result in results:
                                st.markdown(f"**Title:** {result['metadata']['title']}")
                                st.markdown(f"**Excerpt:** {result['document'][:200]}...")
                                st.markdown(f"**Relevance Score:** {1 - result['distance']:.2f}")
                                st.markdown("---")
                        else:
                            st.info("No search results found.")
            
                except ImportError as e:
                    st.error(f"Missing dependencies: {e}")
                    st.info("Ensure you have installed chromadb and sentence-transformers")
            
                except Exception as e:
                    st.error(f"Error building vector store: {e}")
                    
                    # Provide more detailed troubleshooting information
                    st.markdown("### Troubleshooting Tips:")
                    st.markdown("1. Ensure Chroma DB is installed correctly")
                    st.markdown("2. Check that you have write permissions in the project directory")
                    st.markdown("3. Verify that structured transcripts exist")
                    st.markdown("4. Check for any permission or disk space issues")

    def count_characters(self, text):
        """Count Hindi and total characters in text"""
        if not text:
            return 0, 0
        
        def is_hindi(char):
            return '\u0900' <= char <= '\u097F'
        
        hi_chars = sum(1 for char in text if is_hindi(char))
        return hi_chars, len(text)

    def main(self):
        """
        Main application interface
        """
        st.title("🇮🇳 हिंदी सीखो (Learn Hindi)")
        
        # Sidebar for exercise generation
        st.sidebar.header("Exercise Generator")
        
        # Topic selection
        topics = [
            "Everyday Conversations", 
            "Office Dialogues", 
            "Travel", 
            "Family", 
            "Technology"
        ]
        selected_topic = st.sidebar.selectbox(
            "Select Topic", 
            topics
        )
        
        # Difficulty selection
        difficulties = {
            'शुरुआती (Beginner)': 'शुरुआती', 
            'मध्यम (Intermediate)': 'मध्यम', 
            'उन्नत (Advanced)': 'उन्नत'
        }
        selected_difficulty = st.sidebar.selectbox(
            "Select Difficulty", 
            list(difficulties.keys())
        )
        
        # Stage selection
        stages = [
            "Original Exercise Generator", 
            "Chat with Nova", 
            "Transcript Processing"
        ]
        selected_stage = st.sidebar.radio(
            "Select Learning Stage", 
            stages
        )
        
        # Render appropriate stage
        if selected_stage == "Original Exercise Generator":
            # Generate Exercise Button
            if st.sidebar.button("Generate Exercise"):
                self.generate_exercise(
                    topic=selected_topic, 
                    difficulty=difficulties[selected_difficulty]
                )
            
            # Render Exercise
            self.render_exercise()
        
        elif selected_stage == "Chat with Nova":
            self.render_chat_stage()
        
        elif selected_stage == "Transcript Processing":
            self.render_transcript_processing_stage()

def main():
    app = HindiLearningApp()
    app.main()

if __name__ == "__main__":
    main()