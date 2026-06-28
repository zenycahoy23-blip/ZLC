#!/bin/sh

echo "Step 1: Registering user..."
TIMESTAMP=$(date +%s)
EMAIL="test$TIMESTAMP@example.com"
PASSWORD="TestPassword123"
PHONE="+1234567890"

curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"password\":\"$PASSWORD\",\"phone_number\":\"$PHONE\"}"

echo -e "\n\nStep 2: Getting verification token from DB..."
TOKEN=$(docker exec secure-auth-db psql -U authuser -d authdb -t -c "SELECT verification_token FROM users WHERE email = '$EMAIL';")
echo "Token: '$TOKEN' (length: ${#TOKEN})"

echo -e "\n\nStep 3: Verifying email..."
curl -X POST http://localhost:8080/api/verify-email \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$EMAIL\",\"token\":\"$TOKEN\"}"

echo -e "\n\nStep 4: Getting user ID and OTP..."
USER_ID=$(docker exec secure-auth-db psql -U authuser -d authdb -t -c "SELECT id FROM users WHERE email = '$EMAIL';")
OTP=$(docker exec secure-auth-db psql -U authuser -d authdb -t -c "SELECT phone_otp FROM users WHERE id = $USER_ID;")
echo "User ID: $USER_ID"
echo "OTP: '$OTP' (length: ${#OTP})"

echo -e "\n\nBackend logs (PHONE_OTP and VERIFY events):"
docker logs secure-auth-api | grep -E "PHONE_OTP|VERIFY_PHONE|VERIFY_EMAIL|SEND_PHONE"

echo -e "\n\nStep 5: Verifying phone OTP..."
curl -X POST http://localhost:8080/api/verify-phone-otp \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":$USER_ID,\"otp\":\"$OTP\"}"

echo -e "\n\nAll backend logs:"
docker logs secure-auth-api | tail -50
