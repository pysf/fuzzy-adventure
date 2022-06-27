#!/usr/bin/env bash

docker build -t rsc .
docker run rsc -help
echo "########### HI THERE ###########"
echo ""
echo "=> To run the cli app with default parametars run:"
echo "docker run rsc"
echo ""
echo "=> To read tha app help run:"
echo "docker run rsc -help"
echo ""
echo "############## END #############"