#!/bin/sh
cd ../imagedescriptor
python3 -m venv .env
source .env/bin/activate
pip install pillow
pip install 'transformers[torch]'