#!/usr/bin/env bash

HOST=${HOST:-"localhost:32000"}
URL="${HOST}/api/customer"
USER=1
if [[ -n "$1" ]]
then
  USER="${1}"
fi


echo "getting user ${USER}"
curl  "${URL}/${USER}"
echo ""

