# FTPBot

Interact with the filesystem of a remote computer or server from your PC or smartphone using a Telegram client.

## Getting started

1. Grab the [latest release](https://github.com/G07cha/ftpbot/releases)
2. [Create telegram bot](https://core.telegram.org/bots#3-how-do-i-create-a-bot)
3. Copy token from BotFather and run the latest binary with `--token "%YOUR_TOKEN%"` argument

You can check all available options by running `ftpbot --help`. Don't see an option that you want in the list? Submit an issue about this!

## Development

Want to fix a bug or add future? Nice! If you're working on a thing that isn't listed in issues make sure to discuss it first.

There are 2 ways of building application:

### Build locally
Requirements:
- Go 1.6+
- make

```bash
make install
make
./bin/ftpbot --token "YOUR_TOKEN"
```

### Use Docker

Insert `token` argument in Dockerfile and run next commands:

```bash
docker build .
docker run %container_id%
```

## License

MIT © [Konstantin Azizov](http://g07cha.github.io)
