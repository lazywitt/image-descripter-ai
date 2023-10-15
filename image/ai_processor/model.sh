#!/bin/sh
cd ../imagedescriptor
source .env/bin/activate
export IMAGE_URL=$1
echo $IMAGE_URL
python model.py