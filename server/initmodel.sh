#!/bin/sh
mkdir -p ../imagedescriptor
cp -f ./model.py ../imagedescriptor
cd ../imagedescriptor
python3 -m venv .env
source .env/bin/activate
pip install pillow transformers[torch]  --no-cache-dir