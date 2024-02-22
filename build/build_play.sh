#!/bin/sh
echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../play
npm install
npm run build
cd ../build

echo remove last package if exist
if [ -e package/web/play ]; then
  rm -rf package/web/play
fi

mv ../play/build ./package/web/play

echo play package build over.
