import streamlit as st
import requests

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
        # Display only the English word
        st.write(f'**English Word:** {random_word_data.get("english", "N/A")}')
        
        # Button to upload an image
        uploaded_file = st.file_uploader('Upload an image of the Hindi word', type=['jpg', 'jpeg', 'png'])
        if uploaded_file is not None:
            # Here you can add functionality to process the uploaded image
            st.image(uploaded_file, caption='Uploaded Image.', use_column_width=True)
    else:
        st.error("Failed to fetch word.")

# Button to submit the answer
if st.button('Submit'):
    pass
