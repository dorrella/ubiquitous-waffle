#!/usr/bin/env bash

HOST=${HOST:-"localhost:32000"}
URL="${HOST}/api/customer/list"

if [[ -n "$1" ]]
then
  URL="${URL}?next=${1}"
fi


echo "listing users"
curl  ${URL}
echo ""

