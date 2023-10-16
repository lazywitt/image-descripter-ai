from transformers import pipeline
import os


captioner = pipeline("image-to-text",model="Salesforce/blip-image-captioning-base")
print( captioner(os.environ['IMAGE_URL'])[0]['generated_text'])
