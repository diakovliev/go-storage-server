#! /bin/sh -x

SERVER_BASE_URL="http://127.0.0.1:5555/api/v1"
# CURL="curl --no-progress-meter"
CURL="curl"

SID=$(${CURL} "${SERVER_BASE_URL}/storage/create/default" | jq -r '.sid')

echo "sid: ${SID}"

${CURL} -X PUT -d "test file1 content\n" "${SERVER_BASE_URL}/storage/${SID}/test_file1?mode=0777"
${CURL} -X GET "${SERVER_BASE_URL}/storage/${SID}/test_file1"
${CURL} -X PUT -d "test file2 content\n" "${SERVER_BASE_URL}/storage/${SID}/dir/test_file2?mode=0777"
${CURL} -X GET "${SERVER_BASE_URL}/storage/${SID}/dir/test_file2"


BID=$(${CURL} "${SERVER_BASE_URL}/storage/buffer/create" | jq -r '.sid')
echo "bid: ${BID}"

${CURL} -X PUT -d "123\n" "${SERVER_BASE_URL}/storage/buffer/${BID}"
${CURL} -X PUT -d "456\n" "${SERVER_BASE_URL}/storage/buffer/${BID}"
${CURL} -X PUT -d "789\n" "${SERVER_BASE_URL}/storage/buffer/${BID}"
${CURL} -X PUT -d "0\n" "${SERVER_BASE_URL}/storage/buffer/${BID}"

${CURL} -X GET "${SERVER_BASE_URL}/storage/buffer/commit/${SID}/${BID}/test_file3?mode=0777"

${CURL} -X GET "${SERVER_BASE_URL}/storage/list/${SID}"
${CURL} -X GET "${SERVER_BASE_URL}/storage/destroy/${SID}"
