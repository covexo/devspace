@echo off

:: Build the docker image
docker build -t %1 . -f custom/Dockerfile