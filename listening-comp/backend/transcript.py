from youtube_transcript_api import YouTubeTranscriptApi
from typing import Optional, Dict, List
import re

class YTTranscriptDownloader():
    
    def __init__(self):
        self.languages = ['hi', 'en']
    
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