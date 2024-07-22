while [[ "$1" = -* ]] ; do
	case "$1" in
		-s|--sudo) SUDO=sudo ;;
		-h|--help)
			echo "Usage: $0 [options]"
			echo "Options are:"
			echo -e "-s|--sudo\\tuse sudo to startup/shutdows bank simulator (default: no sudo)"
			echo -e "-h|--help\\tprint this help and exit"
			exit ;;
		*) break ;;
	esac
	shift
done

VALID_REQUESTS=(
'{ "card_number": "2222405343248877", "expiry_year": 2025, "expiry_month": 4, "amount": 100, "currency": "GBP", "cvv": "123" }'
'{ "status":"Authorized","card_number":"8877","expiry_month":4,"expiry_year":2025,"currency":"GBP","amount":100 }'

'{ "card_number": "2222405343248112", "expiry_year": 2026, "expiry_month": 1, "amount": 60000, "currency": "USD", "cvv": "456" }'
'{ "status":"Unauthorized","card_number":"8112","expiry_month":1,"expiry_year":2026,"currency":"USD","amount":60000 }'
)

INVALID_REQUESTS=(
'{ "card_number": "123456", "expiry_year": 2026, "expiry_month": 1, "amount": 60000, "currency": "USD", "cvv": "456" }'
'{ "error": "Key: '\''PaymentRequest.CardNumber'\'' Error:Field validation for '\''CardNumber'\'' failed on the '\''min'\'' tag" }'

'{ "card_number": "12345678901234567890", "expiry_year": 2026, "expiry_month": 1, "amount": 60000, "currency": "USD", "cvv": "456" }'
'{ "error": "Key: '\''PaymentRequest.CardNumber'\'' Error:Field validation for '\''CardNumber'\'' failed on the '\''max'\'' tag" }'

'{ "card_number": "222240534.3248877", "expiry_year": 2025, "expiry_month": 4, "amount": 100, "currency": "GBP", "cvv": "123" }'
'{ "error": "Key: '\''PaymentRequest.CardNumber'\'' Error:Field validation for '\''CardNumber'\'' failed on the '\''numeric_only'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2023, "expiry_month": 4, "amount": 100, "currency": "GBP", "cvv": "123" }'
'{ "error": "Key: '\''PaymentRequest.ExpiryYear'\'' Error:Field validation for '\''ExpiryYear'\'' failed on the '\''min'\'' tag\nKey: '\''PaymentRequest.expiry_year'\'' Error:Field validation for '\''expiry_year'\'' failed on the '\''expiryyear'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2024, "expiry_month": 4, "amount": 100, "currency": "GBP", "cvv": "123" }'
'{ "error": "Key: '\''PaymentRequest.expiry_month/expiry_year'\'' Error:Field validation for '\''expiry_month/expiry_year'\'' failed on the '\''expirydate'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2025, "expiry_month": 2025, "amount": 100, "currency": "GBP", "cvv": "123" }'
'{ "error": "Key: '\''PaymentRequest.ExpiryMonth'\'' Error:Field validation for '\''ExpiryMonth'\'' failed on the '\''max'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2024, "expiry_month": 4, "amount": -100, "currency": "GBP", "cvv": "123" }'
'{ "error": "json: cannot unmarshal number -100 into Go struct field PaymentRequest.amount of type uint32" }'

'{ "card_number": "2222405343248877", "expiry_year": 2024, "expiry_month": 4, "amount": 500000000000, "currency": "GBP", "cvv": "123" }'
'{ "error": "json: cannot unmarshal number 500000000000 into Go struct field PaymentRequest.amount of type uint32" }'

'{ "card_number": "2222405343248877", "expiry_year": 2025, "expiry_month": 4, "amount": 100, "currency": "not a currency", "cvv": "123" }'
'{ "error": "Key: '\''PaymentRequest.Currency'\'' Error:Field validation for '\''Currency'\'' failed on the '\''iso4217'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2025, "expiry_month": 4, "amount": 100, "currency": "RUB", "cvv": "123" }'
'{ "error": "Key: '\''PaymentRequest.Currency'\'' Error:Field validation for '\''Currency'\'' failed on the '\''known_currency'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2025, "expiry_month": 4, "amount": 100, "currency": "GBP", "cvv": "12" }'
'{ "error": "Key: '\''PaymentRequest.CVV'\'' Error:Field validation for '\''CVV'\'' failed on the '\''min'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2025, "expiry_month": 4, "amount": 100, "currency": "GBP", "cvv": "12345" }'
'{ "error": "Key: '\''PaymentRequest.CVV'\'' Error:Field validation for '\''CVV'\'' failed on the '\''max'\'' tag" }'

'{ "card_number": "2222405343248877", "expiry_year": 2025, "expiry_month": 4, "amount": 200, "currency": "GBP", "cvv": "123" }'
'{ "error": "Invalid response from the bank" }'
)

INVALID_RECALLS=(
'hello'                                 '{ "error": "invalid UUID length: 5" }'
'0190da8b-8fea-7bc4-9252-c1e212eca33'   '{ "error": "invalid UUID length: 35" }'
'0190da8b-8fea-7bc4-9252-c1e212eca3321' '{ "error": "invalid UUID length: 37" }'
)

declare -a ASSIGNED_IDS

cd $(dirname $0)/..
$SUDO docker-compose up -d
coproc PAYMENT_PROC { go run main.go redirects.go ; }
echo "Payment processor at $PAYMENT_PROC_PID"
sleep 2
FAILS=0

function checkForFail() {
	local CASE_ID=$1
	local EXPECTED=$2
	local ACTUAL=$3

	if [[ "$EXPECTED" != "$ACTUAL" ]] ; then
		echo "FAIL at $CASE_ID"
		echo Expected:
		echo "$EXPECTED" | sed 's/^/    /'
		echo Actual:
		echo "$ACTUAL" | sed 's/^/    /'

		EXP_JSON=$(mktemp)
		ACT_JSON=$(mktemp)
		echo "$EXPECTED" > "$EXP_JSON"
		echo "$ACTUAL" > "$ACT_JSON"
		echo Diff:
		diff "$EXP_JSON" "$ACT_JSON"
		rm "$EXP_JSON" "$ACT_JSON"

		FAILS=$((FAILS + 1))
	fi
}

for ((i = 0; i < ${#INVALID_REQUESTS[@]}; i += 2)) ; do
	echo Invalid request \# $((i / 2 + 1))
	RESP=$(curl http://localhost:8090/pay -H 'Content-Type: application/json' -d "${INVALID_REQUESTS[$i]}")
	EXPECTED=$(echo ${INVALID_REQUESTS[$((i + 1))]} | jq .)
	ACTUAL=$(echo $RESP | jq 'del(.id)')
	checkForFail "invalid request case $((i / 2 + 1)):" "$EXPECTED" "$ACTUAL"
done

for ((i = 0; i < ${#INVALID_RECALLS[@]}; i += 2)) ; do
	echo Invalid recall \# $((i / 2 + 1))
	RESP=$(curl "http://localhost:8090/recall?payment_id=${INVALID_RECALLS[$i]}")
	EXPECTED=$(echo ${INVALID_RECALLS[$((i + 1))]} | jq .)
	ACTUAL=$(echo $RESP | jq .)
	checkForFail "invalid recall case $((i / 2 + 1)):" "$EXPECTED" "$ACTUAL"
done

for ((i = 0; i < ${#VALID_REQUESTS[@]}; i += 2)) ; do
	echo Valid request \# $((i / 2 + 1))
	RESP=$(curl http://localhost:8090/pay -H 'Content-Type: application/json' -d "${VALID_REQUESTS[$i]}")
	ID=$(echo "$RESP" | jq .id | sed 's/"//g')
	ASSIGNED_IDS[$i]=$ID
	EXPECTED=$(echo ${VALID_REQUESTS[$((i + 1))]} | jq .)
	ACTUAL=$(echo $RESP | jq 'del(.id)')
	checkForFail "valid request case $((i / 2 + 1)):" "$EXPECTED" "$ACTUAL"
done

for ((i = 0; i < ${#VALID_REQUESTS[@]}; i += 2)) ; do
	echo Recall request \# $((i / 2 + 1))
	RESP=$(curl "http://localhost:8090/recall?payment_id=${ASSIGNED_IDS[$i]}")
	ID=$(echo "$RESP" | jq .id | sed 's/"//g')
	if [[ "${ASSIGNED_IDS[$i]}" != "$ID" ]] ; then
		echo "FAIL at recall case $((i / 2 + 1)):"
		echo "Expected: ${ASSIGNED_IDS[$i]}"
		echo "Actual  : $ID"
	fi

	EXPECTED=$(echo ${VALID_REQUESTS[$((i + 1))]} | jq .)
	ACTUAL=$(echo $RESP | jq 'del(.id)')
	checkForFail "recall case $((i / 2 + 1)):" "$EXPECTED" "$ACTUAL"
done

MISSING_ID="${ASSIGNED_IDS[0]}"
while [[ "$MISSING_ID" = "${ASSIGNED_IDS[0]}" || "$MISSING_ID" = "${ASSIGNED_IDS[1]}" ]] ; do
	MISSING_ID=$(uuid)
done

echo Recal for unknown id "$MISSING_ID"

RESP=$(curl "http://localhost:8090/recall?payment_id=$MISSING_ID")
EXPECTED_RESP="{\"error\":\"Payment $MISSING_ID not found\"}"
if [[ "$EXPECTED_RESP" != "$RESP" ]] ; then
	echo "FAIL at unknonwn payment id case"
	echo "Expected: $EXPECTED_RESP"
	echo "Actual  : $RESP"
fi

kill $PAYMENT_PROC_PID
kill $(netstat -ltnp 2>/dev/null | grep :::8090 | gawk '{ sub(/\/.*/, "", $NF); print $NF }')
$SUDO docker-compose down
if [[ $FAILS -eq 0 ]] ; then echo All OK ; else echo $FAILS FAILs ; fi
