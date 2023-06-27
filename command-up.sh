#!/bin/bash

. .env

# create-channelコマンドの登録
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
echo "create-channel response: $response"


# register-apikeyコマンドの登録
url="https://discord.com/api/v10/applications/$BOT_ID01/guilds/$GUILD_ID01/commands"

json='{
    "name": "register-apikey",
    "description": "Register the apikey to post youtube videos",
    "options": [
        {
            "name": "apikey",
            "description": "Write the apikey of Youtube",
            "type": 3,
            "required": true
        }
    ]
}'

headers=(
    "Authorization: Bot $BOT_TOKEN01"
)

response=$(curl -X POST -H "${headers[@]}" -H "Content-Type: application/json" -d "$json" "$url")
echo "register-apikey response: $response"


# add-blacklistコマンドの登録
url="https://discord.com/api/v10/applications/$BOT_ID01/guilds/$GUILD_ID01/commands"

json='{
    "name": "add-blacklist",
    "description": "Add the channel to the blacklist[format:youtube.com/@~~~]",
    "options": [
        {
            "name": "channel-url",
            "description": "Please post the Youtube channel link you wish to add to the blacklist",
            "type": 3,
            "required": true
        }
    ]
}'

headers=(
    "Authorization: Bot $BOT_TOKEN01"
)

response=$(curl -X POST -H "${headers[@]}" -H "Content-Type: application/json" -d "$json" "$url")
echo "add-blacklist response: $response"


# remove-blacklistコマンドの登録
url="https://discord.com/api/v10/applications/$BOT_ID01/guilds/$GUILD_ID01/commands"

json='{
    "name": "remove-blacklist",
    "description": "Remove the channel from the blacklist[format:youtube.com/@~~~]",
    "options": [
        {
            "name": "channel-url",
            "description": "Please post the Youtube channel link you wish to remove from the blacklist",
            "type": 3,
            "required": true
        }
    ]
}'

headers=(
    "Authorization: Bot $BOT_TOKEN01"
)

response=$(curl -X POST -H "${headers[@]}" -H "Content-Type: application/json" -d "$json" "$url")
echo "remove-blacklist response: $response"


# get-blacklistコマンドの登録
url="https://discord.com/api/v10/applications/$BOT_ID01/guilds/$GUILD_ID01/commands"

json='{
    "name": "get-blacklist",
    "description": "Get the blacklist"
}'

headers=(
    "Authorization: Bot $BOT_TOKEN01"
)

response=$(curl -X POST -H "${headers[@]}" -H "Content-Type: application/json" -d "$json" "$url")
echo "remove-blacklist response: $response"

