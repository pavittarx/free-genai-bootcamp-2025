import streamlit as st
import requests
import pytesseract
from PIL import Image

# Title of the app
st.title('Hindi Writing Practice App')

# Fetch word groups from the API
response = requests.get('http://localhost:3000/api/groups')
if response.status_code == 200:
    data = response.json()
    if 'groups' in data:
        word_groups = data['groups']
        group_names = [group['name'] for group in word_groups]
    else:
        st.error("'groups' key not found in the response.")
        group_names = []
else:
    st.error("Failed to fetch data from API.")
    group_names = []

# Dropdown for selecting word groups
selected_group = st.selectbox('Select a word group:', group_names)

# Button to fetch a random word based on the selected group
if st.button('Generate Word'):
    group_id = word_groups[group_names.index(selected_group)]['id']
    random_word_response = requests.get(f'http://localhost:3000/api/words/random?group_id={group_id}', headers={'accept': 'application/json'})
    if random_word_response.status_code == 200:
        random_word_data = random_word_response.json()
        # Store the fetched word in session state
        st.session_state.fetched_word = random_word_data.get("english", "N/A")
        st.write(f'**English Word:** {st.session_state.fetched_word}')

# Check if an image has been uploaded previously
if 'uploaded_image' in st.session_state:
    uploaded_file = st.session_state.uploaded_image
    # Load and display the image using Pillow
    image = Image.open(uploaded_file)
    st.image(image, caption='Uploaded Image.', use_column_width=True)
    # Use Tesseract to extract text
    extracted_text = pytesseract.image_to_string(image, lang='hin')  # Specify Hindi language
    st.write(f'**Extracted Hindi Text:** {extracted_text}')
else:
    # File uploader for image upload
    uploaded_file = st.file_uploader('Upload an image of the Hindi word', type=['jpg', 'jpeg', 'png'], key='image_uploader')
    if uploaded_file is not None:
        # Store the uploaded image in session state
        st.session_state.uploaded_image = uploaded_file
        # Load and display the image using Pillow
        image = Image.open(uploaded_file)
        st.image(image, caption='Uploaded Image.', use_column_width=True)
        # Use Tesseract to extract text
        extracted_text = pytesseract.image_to_string(image, lang='hin')  # Specify Hindi language
        st.write(f'**Extracted Hindi Text:** {extracted_text}')

# Check if the fetched word exists in session state
if 'fetched_word' in st.session_state:
    st.write(f'**Fetched Word:** {st.session_state.fetched_word}')  

# Button to submit the answer
if st.button('Submit'):
    pass

# # Debugging: Print session state to check if the image is stored
# st.write(st.session_state)
