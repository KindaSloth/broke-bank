# script to test concurrent/parallel requests
function perform_transfer {
  curl --cookie "sessionId=0191068b-631b-76b7-85b0-979ffd32a0c5; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly" \
       --request POST \
       --url http://localhost:5000/transaction/transfer \
       --header 'Content-Type: application/json' \
       --data '{
        "amount": "1.00",
        "to_account_id": "0191068b-1b09-7c70-8c61-d9dc4078f322",
        "from_account_id": "0191068a-f34c-7516-b6a3-4abeca213131"
       }'
}

export -f perform_transfer

num_requests=100

parallelism=10

seq $num_requests | xargs -n1 -P$parallelism -I{} bash -c 'perform_transfer'
