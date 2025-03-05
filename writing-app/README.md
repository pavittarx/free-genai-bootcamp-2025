# Hindi Writing App
The goal is to build a hindi learning writing practice app. The user will be presented with an english word and the user has to write the corresponding hindi word. 
The word can be written digitally or on paper. The image of the word can be captured and uploaded. The uploaded image will then be analysed for correctness. 

## Technical Uncertainities
- Whether or not an OCR library will be available for hindi language.
- No prior experience with vision encoder decoder models.

## Technical Spec
- The app will be using Streamlit for frontend. 
- The will use existing lang portal backend for getting groups. 
- Existing Words will be used from LangPortal, as it already supports fetching random words by group. There are around 200 words in total.  
- The app will evaluate words for the sake of simplicity.
- Poetry used for packaging the application.
- TessaractOCR is being used for Hindi Lang Extraction.

## Extras
You need to install tessaract in order to use OCR.

```shell
    # Install Tessaract
    sudo apt-get install tesseract-ocr
    # Install Hindi Lang Pack
    sudo apt-get install tesseract-ocr-hin

    export TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata
```