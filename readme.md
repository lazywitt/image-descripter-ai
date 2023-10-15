IMAGE TAGGING SERVICE

API Functionalities - 
    * "localhost:8080/image/register" - Add images to the service which will be saved in the local file system.
    * "localhost:8080/image/search" - get the file location of an image based on a given description, may return empty response


Working -
    * Spawns two processes
        Server - listens to the incoming api request.
        Cron - asynchronously runs a pipeline to tag images with AI generated description and marks them processed.

BootStrap Process - 
    * Triggers a python virtual environment in a parallel directory ../imagedescriptor and installs AI functionalities
    * Setup PGDB and entities, create a PGDB client
    * Triggers TCP server process
    * Triggers cron on a seperate GOROUTINE

Pre-requisite - 
    * postgres, golang, Python, virtual env should be installed
    * change ./configs/pgdb-dev.yml config to connect to your local pgdb

Running Instruction - 
    * If pre requistes are installed, run "bash ./shell.sh" from root

TODO -
    *   Dockerise
    *   Use better means of string/description matching for search api. Currently using PGDB GIN index
