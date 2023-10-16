#!/bin/sh
cd ../imagedescriptor
source .env/bin/activate
export IMAGE_URL=$1
python model.py