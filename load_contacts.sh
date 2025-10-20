#!/bin/bash

# Ensure script runs with bash
if [ -z "$BASH_VERSION" ]; then
  echo "This script requires Bash. Run it with: bash generate_contacts.sh"
  exit 1
fi

# Arrays for generating contact data
first_names=("John" "Jane" "Michael" "Emily" "David" "Sarah" "James" "Laura" "Robert" "Lisa")
last_names=("Smith" "Johnson" "Brown" "Taylor" "Wilson" "Davis" "Clark" "Harris" "Lewis" "Walker")
streets=("Main St" "Park Ave" "Oak Rd" "Cedar Ln" "Maple Dr" "Elm St" "Pine Rd" "Birch Ave" "Spruce Ln" "Walnut Dr")
cities=("Anytown" "Springfield" "Riverside" "Greenville" "Fairview" "Lakewood" "Hillcrest" "Brookside")
states=("CA" "NY" "TX" "FL" "PA" "IL" "OH" "GA")

# Output file
output_file="contacts_output.txt"

# Initialize output file
> "$output_file"

# Generate 100 contacts
for i in {1..100}
do
  # Randomly select data
  first_name=${first_names[$((RANDOM % ${#first_names[@]}))]}
  last_name=${last_names[$((RANDOM % ${#last_names[@]}))]}
  street="${streets[$((RANDOM % ${#streets[@]}))]} $((RANDOM % 1000 + 1))"
  city=${cities[$((RANDOM % ${#cities[@]}))]}
  state=${states[$((RANDOM % ${#states[@]}))]}
  zip=$((RANDOM % 90000 + 10000))  # 5-digit zip code
  phone="555$((RANDOM % 900 + 100))$((RANDOM % 9000 + 1000))"  # Format: 555XXXYYYY
  email=$(echo "${first_name}.${last_name}@example.com" | tr '[:upper:]' '[:lower:]')

  # Create JSON payload
  json_payload=$(cat << EOF
{
  "first_name": "$first_name",
  "last_name": "$last_name",
  "phone": "$phone",
  "email": "$email",
  "street": "$street",
  "city": "$city",
  "state": "$state",
  "zip": "$zip"
}
EOF
)

  # Send POST request and capture response
  response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST -H "Content-Type: application/json" -d "$json_payload" http://localhost:8080/records 2>&1)
  status_code=$(echo "$response" | grep "HTTP_STATUS" | cut -d':' -f2)
  response_body=$(echo "$response" | grep -v "HTTP_STATUS")

  # Check if request was successful
  if [ "$status_code" -eq 201 ]; then
    echo "Contact $i created: $first_name $last_name" | tee -a "$output_file"
  else
    echo "Failed to create contact $i: Status $status_code, Response: $response_body" | tee -a "$output_file"
  fi
done

echo "Generated 100 contacts. Output logged to $output_file" | tee -a "$output_file"
