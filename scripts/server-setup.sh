#!/bin/bash

valid_ip=false

while ! $valid_ip ; do
  echo "Input control-plane server ip:"
  read -r server_ip
  if [[ $server_ip =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    valid_ip=true
  fi
done

