from vector_db import TranscriptVectorDB

def main():
    # Read structured transcripts
    transcripts = TranscriptVectorDB.read_structured_transcripts()
    
    print(f"Found {len(transcripts)} transcripts")
    
    # Print details of first transcript
    if transcripts:
        print("\nFirst Transcript Details:")
        for key, value in transcripts[0].items():
            print(f"{key}: {value}")
    
    # Initialize vector database
    vector_db = TranscriptVectorDB()
    vector_db.add_transcripts(transcripts)
    
    # Perform sample searches
    print("\nSample Searches:")
    searches = ["मौसम", "धूप", "बारिश"]
    for query in searches:
        print(f"\nSearching for: {query}")
        results = vector_db.search_transcripts(query)
        for result in results:
            print(f"Title: {result['metadata']['title']}")
            print(f"Excerpt: {result['document'][:200]}...")
            print(f"Distance: {result['distance']}\n")

if __name__ == "__main__":
    main()
