#!/bin/sh
echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../monitor
npm install
npm run build
cd ../build

echo remove last package if exist
if [ -e package/web/monitor ]; then
  rm -rf package/web/monitor
fi

mv ../monitor/build ./package/web/monitor

echo monitor package build over.
