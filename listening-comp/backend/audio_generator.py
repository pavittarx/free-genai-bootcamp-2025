import os
import uuid
from gtts import gTTS
from typing import Dict, Any, Optional

class AudioGenerator:
    """
    Audio generation utility for Hindi language exercises
    Generates a single, comprehensive audio file for each exercise
    """
    
    @staticmethod
    def generate_audio(
        exercise: Dict[str, Any], 
        audio_dir: str = "./frontend/static/audio"
    ) -> Optional[str]:
        """
        Generate a single comprehensive audio file for the entire exercise
        
        Args:
            exercise (Dict[str, Any]): Exercise dictionary containing text components
            audio_dir (str): Directory to save generated audio files
        
        Returns:
            Optional path to the generated audio file
        """
        # Ensure audio directory exists
        os.makedirs(audio_dir, exist_ok=True)
        
        # Combine all text components into a single narrative
        audio_text_components = [
            "Language Learning Exercise",
            exercise.get('introduction', ''),
            "Dialogue begins",
            exercise.get('dialogue', ''),
            "Question",
            exercise.get('question', ''),
            "Options:",
            *exercise.get('options', [])
        ]
        
        # Join components with pauses
        full_text = " . ".join(filter(bool, audio_text_components))
        
        try:
            # Create gTTS object with Hindi language
            tts = gTTS(text=full_text, lang='hi')
            
            # Generate unique filename for the exercise
            exercise_id = str(uuid.uuid4())[:8]
            filename = f"{exercise_id}_exercise.mp3"
            filepath = os.path.join(audio_dir, filename)
            
            # Save audio file
            tts.save(filepath)
            
            # Return absolute path for frontend
            return filepath
        
        except Exception as e:
            print(f"Error generating audio: {e}")
            return None

    @classmethod
    def cleanup_old_audio_files(
        cls,
        audio_dir: str = "./frontend/static/audio", 
        max_files: int = 10
    ):
        """
        Clean up old audio files to prevent disk space overflow
        
        Args:
            audio_dir (str): Directory containing audio files
            max_files (int): Maximum number of files to keep
        """
        try:
            # Get all mp3 files sorted by creation time
            files = [os.path.join(audio_dir, f) for f in os.listdir(audio_dir) if f.endswith('.mp3')]
            files.sort(key=os.path.getctime)
            
            # Remove oldest files if exceeding max_files
            while len(files) > max_files:
                oldest_file = files.pop(0)
                os.remove(oldest_file)
        
        except Exception as e:
            print(f"Error cleaning up audio files: {e}")

# Example usage remains the same
def main():
    # Sample exercise for testing
    sample_exercise = {
        "introduction": "यह एक हिंदी भाषा अभ्यास है",
        "dialogue": "राम: नमस्ते! आप कैसे हैं?\nश्याम: मैं बिल्कुल ठीक हूँ, धन्यवाद।",
        "question": "इस संवाद में किसने पहले नमस्ते कहा?",
        "options": ["राम", "श्याम", "कोई नहीं", "दोनों"],
        "answer": "राम"
    }
    
    audio_generator = AudioGenerator()
    
    # Generate single exercise audio
    exercise_audio = audio_generator.generate_audio(sample_exercise)
    print("Exercise Audio:", exercise_audio)
    
    # Cleanup old audio files
    AudioGenerator.cleanup_old_audio_files()

if __name__ == "__main__":
    main()
