#!/bin/sh

dockerd > /dev/null 2>&1 &
for i in {0..30}
do
  docker info > /dev/null 2>&1
  if [ $? == 0 ]; then
    drone-image-copy-plugin $*
    exit $?
  fi
  sleep 1s
done

echo Unabled to start docker daemon. 1>&2
