#!/usr/bin/env bash

echo ""
echo "STARTING SETUP"
echo ""

sudo apt-get update

echo "INSTALLING MONGODB"
echo ""

sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 2930ADAE8CAF5059EE73BB4B58712A2291FA4AD5
echo "deb http://repo.mongodb.org/apt/debian jessie/mongodb-org/3.6 main" | sudo tee /etc/apt/sources.list.d/mongodb-org-3.6.list

sudo apt-get update
sudo apt-get install -y mongodb-org

echo ""
echo "FINISHED INSTALLING MONGODB"

echo ""

echo "INSTALLING REDIS SERVER"
echo ""

sudo apt-get install redis-server

echo ""
echo "FINISHED INSTALLING REDIS SERVER"

echo ""
echo "DONE"
echo ""