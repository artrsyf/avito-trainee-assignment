curl --location --request POST 'http://0.0.0.0:8080/api/auth' \
--data ''

400

{
    "errors": "bad request"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "password": "strong_pass"
}'

400

{
    "errors": "Username is required"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user1"
}'

400

{
    "errors": "Password is required"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "us",
  "password": "strong_pass"
}'

400

{
    "errors": "Username is too short"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user1",
  "password": "str"
}'

400

{
    "errors": "Password is too short"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": 1,
  "password": "strong_pass"
}'

400

{
    "errors": "bad request"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user1",
  "password": "strong_pass"
}'

200

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user1",
  "password": "wrong_pass"
}'

401

{
    "errors": "wrong credentials"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user1",
  "password": "strong_pass"
}'

200

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user2",
  "password": "strong_pass2"
}'

200

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQzOjQ0LjIxMDc5NjAyMloiLCJpYXQiOjE3Mzk2NDE0MjQsInVzZXIiOnsiaWQiOiIyIiwidXNlcm5hbWUiOiJ1c2VyMiJ9fQ.q0QT4W5TIn_hwyjpe8-z9g0_JTYXVuGf9EaprhmiOxw"
}

curl --location 'http://0.0.0.0:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{
  "username": "user3",
  "password": "strong_pass3"
}'

200

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQzOjU1LjUxNjM5Mzk3WiIsImlhdCI6MTczOTY0MTQzNSwidXNlciI6eyJpZCI6IjMiLCJ1c2VybmFtZSI6InVzZXIzIn19.A8efQFlzeSzLUhuKd1RbM4vmtJvorkQmdsHTc3L86YE"
}

curl --location 'http://0.0.0.0:8080/api/info'

401

{"errors":"missing token"}

curl --location 'http://0.0.0.0:8080/api/buy/pink-hoody'

401

{"errors":"missing token"}

curl --location --request GET 'http://0.0.0.0:8080/api/buy/sendCoin' \
--header 'Content-Type: application/json' \
--data '{
  "toUser": "user2",
  "amount": 100
}'

401

{"errors":"missing token"}

curl --location 'http://0.0.0.0:8080/api/info' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

{
    "coins": 1000,
    "inventory": [],
    "coinHistory": {
        "received": [],
        "sent": []
    }
}

curl --location 'http://0.0.0.0:8080/api/buy/' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

404

404 page not found

curl --location 'http://0.0.0.0:8080/api/buy/not_available' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

404

{
    "errors": "item not found"
}

curl --location 'http://0.0.0.0:8080/api/buy/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

404

{
    "errors": "item not found"
}

curl --location 'http://0.0.0.0:8080/api/buy/socks' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

curl --location 'http://0.0.0.0:8080/api/buy/socks' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

curl --location 'http://0.0.0.0:8080/api/buy/socks' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

curl --location 'http://0.0.0.0:8080/api/buy/wallet' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

curl --location 'http://0.0.0.0:8080/api/buy/wallet' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

curl --location 'http://0.0.0.0:8080/api/buy/hoody' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

curl --location --request POST 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data ''

400

{
    "errors": "bad request"
}

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user2",
  "amount": 0
}'

400

{
    "errors": "Amount is required"
}

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user2",
  "amount": -1000
}'

400

{
    "errors": "bad request"
}

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user_not_exist",
  "amount": 100
}'

400

{
    "errors": "can't find such user"
}

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user2",
  "amount": 100
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user2",
  "amount": 38
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user3",
  "amount": 120
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user3",
  "amount": 50
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user3",
  "amount": 87
}'

200

curl --location 'http://0.0.0.0:8080/api/info' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

200

{
    "coins": 175,
    "inventory": [
        {
            "type": "hoody",
            "quantity": 1
        },
        {
            "type": "socks",
            "quantity": 3
        },
        {
            "type": "wallet",
            "quantity": 2
        }
    ],
    "coinHistory": {
        "received": [],
        "sent": [
            {
                "toUser": "user2",
                "amount": 138
            },
            {
                "toUser": "user3",
                "amount": 257
            }
        ]
    }
}

curl --location 'http://0.0.0.0:8080/api/buy/hoody' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE'

400

{
    "errors": "not enough balance"
}

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user2",
  "amount": 300
}'

400

{
    "errors": "not enough balance"
}

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQzOjU1LjUxNjM5Mzk3WiIsImlhdCI6MTczOTY0MTQzNSwidXNlciI6eyJpZCI6IjMiLCJ1c2VybmFtZSI6InVzZXIzIn19.A8efQFlzeSzLUhuKd1RbM4vmtJvorkQmdsHTc3L86YE' \
--data '{
  "toUser": "user1",
  "amount": 120
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQzOjU1LjUxNjM5Mzk3WiIsImlhdCI6MTczOTY0MTQzNSwidXNlciI6eyJpZCI6IjMiLCJ1c2VybmFtZSI6InVzZXIzIn19.A8efQFlzeSzLUhuKd1RbM4vmtJvorkQmdsHTc3L86YE' \
--data '{
  "toUser": "user1",
  "amount": 370
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQzOjQ0LjIxMDc5NjAyMloiLCJpYXQiOjE3Mzk2NDE0MjQsInVzZXIiOnsiaWQiOiIyIiwidXNlcm5hbWUiOiJ1c2VyMiJ9fQ.q0QT4W5TIn_hwyjpe8-z9g0_JTYXVuGf9EaprhmiOxw' \
--data '{
  "toUser": "user1",
  "amount": 120
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQzOjQ0LjIxMDc5NjAyMloiLCJpYXQiOjE3Mzk2NDE0MjQsInVzZXIiOnsiaWQiOiIyIiwidXNlcm5hbWUiOiJ1c2VyMiJ9fQ.q0QT4W5TIn_hwyjpe8-z9g0_JTYXVuGf9EaprhmiOxw' \
--data '{
  "toUser": "user1",
  "amount": 28
}'

200

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQzOjQ0LjIxMDc5NjAyMloiLCJpYXQiOjE3Mzk2NDE0MjQsInVzZXIiOnsiaWQiOiIyIiwidXNlcm5hbWUiOiJ1c2VyMiJ9fQ.q0QT4W5TIn_hwyjpe8-z9g0_JTYXVuGf9EaprhmiOxw' \
--data '{
  "toUser": "user1",
  "amount": 89
}'

200

curl --location --request GET 'http://0.0.0.0:8080/api/info' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user1",
  "amount": 89
}'

200

{
    "coins": 902,
    "inventory": [
        {
            "type": "hoody",
            "quantity": 1
        },
        {
            "type": "socks",
            "quantity": 3
        },
        {
            "type": "wallet",
            "quantity": 2
        }
    ],
    "coinHistory": {
        "received": [
            {
                "fromUser": "user2",
                "amount": 237
            },
            {
                "fromUser": "user3",
                "amount": 490
            }
        ],
        "sent": [
            {
                "toUser": "user2",
                "amount": 138
            },
            {
                "toUser": "user3",
                "amount": 257
            }
        ]
    }
}

curl --location 'http://0.0.0.0:8080/api/sendCoin' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAyLTE1VDE5OjQyOjIzLjQ2NTg1MDc2NVoiLCJpYXQiOjE3Mzk2NDEzNDMsInVzZXIiOnsiaWQiOiIxIiwidXNlcm5hbWUiOiJ1c2VyMSJ9fQ.1aFncYmPh_0nvCGYvifOH5MR4PBPgkV2-fZLlUAHHFE' \
--data '{
  "toUser": "user1",
  "amount": 250
}'

400

{
    "errors": "money transfer to yourself is not allowed"
}