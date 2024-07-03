#!/usr/bin/env bash

HOST=${HOST:-"localhost:32000"}
URL="${HOST}/api/customer"
USER=1
if [[ -n "$1" ]]
then
  USER="${1}"
fi


echo "getting user ${USER}"
user=$(curl -s  "${URL}/${USER}")
echo "${user}"
#strip outer json and update name
user=$(echo "${user}" | \
         sed -e 's/{"customer"://' \
             -e 's/}}/}/' \
             -e 's/"name_first":"Tommey"/"name_first":"Bobby"/' \
             -e 's/"name_last":"Jones"/"name_last":"Hill"/')
echo "new user: ${user}"
echo "updating user"
curl -X PUT -H "Content-Type: application/json" -d "${user}" "${URL}/${USER}"
echo ""
