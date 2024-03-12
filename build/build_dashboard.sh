#!/bin/sh
echo create folder for build package ...
if [ ! -e package ]; then
  mkdir package
fi

if [ ! -e package/web ]; then
  mkdir package/web
fi

echo build the code ...
cd ../dashboard
npm install
npm run build
cd ../build

echo remove last package if exist
if [ -e package/web/dashboard ]; then
  rm -rf package/web/dashboard
fi

mv ../dashboard/build ./package/web/dashboard

echo dashboard package build over.
