#!/bin/bash
docker build -t application_image:latest .

docker run --name application_container --network concurrency_process -p 8000:8000 application_image:latest