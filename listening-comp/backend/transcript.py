from youtube_transcript_api import YouTubeTranscriptApi
from typing import Optional, Dict, List
import re
import os
import json
import csv
from datetime import datetime

class YTTranscriptDownloader():
    def __init__(self, save_dir: str = None):
        self.languages = ['hi', 'en']
        # Set default save directory if not provided
        self.save_dir = save_dir or os.path.join(os.path.dirname(__file__), 'transcripts')
        
        # Create save directory if it doesn't exist
        os.makedirs(self.save_dir, exist_ok=True)
    
    def get_video_id(self, url: str) -> Optional[str]:
        """Extracts the video ID from a YouTube URL.
        
        Supports various YouTube URL formats:
        - https://www.youtube.com/watch?v=VIDEO_ID
        - https://youtu.be/VIDEO_ID
        - https://youtube.com/embed/VIDEO_ID
        - https://www.youtube.com/v/VIDEO_ID
        """

        # Regex patterns for different YouTube URL formats
        patterns = [
            r'(?:https?:\/\/)?(?:www\.)?youtube\.com\/watch\?v=([^&\s]+)',  # Standard watch URL
            r'(?:https?:\/\/)?youtu\.be\/([^&\s]+)',  # Shortened URL
            r'(?:https?:\/\/)?(?:www\.)?youtube\.com\/embed\/([^&\s]+)',  # Embed URL
            r'(?:https?:\/\/)?(?:www\.)?youtube\.com\/v\/([^&\s]+)'  # V URL
        ]

        for pattern in patterns:
            match = re.search(pattern, url)
            if match:
                return match.group(1)
        
        return None
    
    def get_transcript(self, url: str) -> List[Dict]:
        
        video_id = self.get_video_id(url)
        
        if not video_id:
            raise ValueError("Invalid YouTube URL")
        
        try: 
            return YouTubeTranscriptApi.get_transcript(video_id, languages=self.languages)
        except Exception as e:
            print(f"Error occured while getting transcript for {url}: {str(e)}")
            return None
    
    def save_transcript(self, 
                        transcript: List[Dict], 
                        url: str = None) -> str:
        if not transcript:
            raise ValueError("No transcript provided to save")
        
        # Determine the most appropriate format
        format = self._detect_best_format(transcript)
        
        # Generate filename using video ID or timestamp
        if url:
            video_id = self.get_video_id(url)
        else:
            video_id = datetime.now().strftime("%Y%m%d_%H%M%S")
        
        filename = f"transcript_{video_id}"
        filename = re.sub(r'[^\w\-_\.]', '_', filename)
        full_path = os.path.join(self.save_dir, f"{filename}.{format}")
        
        os.makedirs(self.save_dir, exist_ok=True)
        
        # Save based on detected format
        try:
            if format == 'json':
                with open(full_path, 'w', encoding='utf-8') as f:
                    json.dump(transcript, f, ensure_ascii=False, indent=2)
            
            elif format == 'csv':
                with open(full_path, 'w', encoding='utf-8', newline='') as f:
                    writer = csv.DictWriter(f, fieldnames=['start', 'duration', 'text'])
                    writer.writeheader()
                    for entry in transcript:
                        writer.writerow({
                            'start': entry.get('start', ''), 
                            'duration': entry.get('duration', ''), 
                            'text': entry.get('text', '')
                        })
            
            elif format == 'txt':
                with open(full_path, 'w', encoding='utf-8') as f:
                    for entry in transcript:
                        f.write(f"[{entry.get('start', 0):.2f}s] {entry.get('text', '')}\n")
            
            else:
                raise ValueError(f"Unsupported format: {format}")
            
            print(f"Transcript saved to {full_path}")
            return full_path
        
        except Exception as e:
            print(f"Error saving transcript: {str(e)}")
            raise

    def _detect_best_format(self, transcript: List[Dict]) -> str:
        if not transcript:
            return 'txt'
        
        first_entry = transcript[0]
        entry_keys = set(first_entry.keys())
        
        if 'start' in entry_keys and 'duration' in entry_keys and 'text' in entry_keys:
            return 'json'
        
        if 'start' in entry_keys and 'text' in entry_keys:
            return 'csv'
        
        return 'txt'
    
def main():
    # Create an instance of YTTranscriptDownloader
    downloader = YTTranscriptDownloader()
    
    # Prompt the user for a YouTube video URL
    print("YouTube Transcript Downloader")
    print("-----------------------------")
    url = input("Enter a YouTube video URL: ").strip()
    
    try:
        transcript = downloader.get_transcript(url)
        
        if transcript:
            saved_file = downloader.save_transcript(transcript, url)
            
            print("\nTranscript Details:")
            print(f"Number of transcript entries: {len(transcript)}")
            print(f"Saved transcript to: {saved_file}")
            
            for entry in transcript[:5]:
                print(f"Time: {entry['start']:.2f}s, Duration: {entry['duration']:.2f}s")
                print(f"Text: {entry['text']}\n")
        else:
            print("No transcript could be retrieved for the given URL.")
    
    except ValueError as ve:
        print(f"Error: {ve}")
    except Exception as e:
        print(f"An unexpected error occurred: {e}")

if __name__ == "__main__":
    main()