#!/bin/bash
cd /home/hype/projects/cine/api_v2
export $(cat .env | grep -v '^#' | xargs)
exec ./build/main
