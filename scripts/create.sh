#!/usr/bin/env bash

HOST=${HOST:-"localhost:32000"}
URL="${HOST}/api/customer"

for i in $(seq 1 40)
do
  FIRST_NAME="Tommey"
  MIDDLE_NAME="Lee"
  LAST_NAME="Jones"
  EMAIL="tjones${i}@gmail.com"
  PHONE="12345${i}"
  echo "creating user ${i} ${EMAIL}"
  curl \
    -F name_first=${FIRST_NAME} \
    -F name_middle=${MIDDLE_NAME} \
    -F name_last=${LAST_NAME} \
    -F email=${EMAIL} \
    -F phone_number=${PHONE} \
    ${URL}
done

echo ""
