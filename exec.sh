docker build -t dnd_bot .

docker run -d --rm --network="host" --name gralhund-bot -e AUTH_KEY=$AUTH_KEY dnd_bot