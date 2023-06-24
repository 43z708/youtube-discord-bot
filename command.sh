#!/bin/bash

. .env

url="https://discord.com/api/v10/applications/$BOT_ID01/guilds/$GUILD_ID01/commands"

json='{
    "name": "create-channel",
    "description": "Create a channel to post youtube videos of the specified search words.",
    "options": [
        {
            "name": "channel_name",
            "description": "Write the name of the channel",
            "type": 3,
            "required": true
        },
        {
            "name": "search_words",
            "description": "Write the search words",
            "type": 3,
            "required": true
        }
    ]
}'

headers=(
    "Authorization: Bot $BOT_TOKEN01"
)

response=$(curl -X POST -H "${headers[@]}" -H "Content-Type: application/json" -d "$json" "$url")
echo $response
status_code=$(echo "$response" | awk '{print $2}')

echo "Status Code: $status_code"