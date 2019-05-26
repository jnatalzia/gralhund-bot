docker build -t dnd_bot .

docker stop gralhund-bot

docker run -d --rm --network="host" --name gralhund-bot -e AUTH_KEY=$AUTH_KEY -e GIPHY_KEY=$GIPHY_KEY dnd_bot