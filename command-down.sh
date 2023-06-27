#!/bin/bash

. .env

url="https://discord.com/api/v10/applications/$BOT_ID01/guilds/$GUILD_ID01/commands"

headers=(
    "Authorization: Bot $BOT_TOKEN01"
)

# Get all commands
response=$(curl -H "${headers[@]}" "$url")

# Extract command IDs
command_ids=$(echo "$response" | jq -r '.[].id')

# Loop over each command ID and delete the command
for command_id in $command_ids; do
    delete_url="$url/$command_id"
    delete_response=$(curl -X DELETE -H "${headers[@]}" "$delete_url")
    echo "Deleted command $command_id: $delete_response"
done