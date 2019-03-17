docker build -t dnd_bot .

docker run -it --rm --name gralhund-bot -e AUTH_KEY=$AUTH_KEY dnd_bot