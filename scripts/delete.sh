#!/usr/bin/env bash

HOST=${HOST:-"localhost:32000"}
URL="${HOST}/api/customer"
USER=1
if [[ -n "$1" ]]
then
  USER="${1}"
fi


echo "deleting user ${USER}"
curl  -X DELETE "${URL}/${USER}"
echo ""
