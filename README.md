# About
It is a telegram-bot that uses the [deepL](https://deepl.com/translate) translator.

##### Installation:
You will need openssl and docker.
```sh
make
```
"Make" uses openssl to create self-signed certificates, then it creates a docker image named deepl-tg. The default port of the container is 8443.
##### Example:
```sh
$ make
mkdir -p ./data
Enter Domain Name: example.de
Generating a 2048 bit RSA private key
.......+++
..............................................................+++
writing new private key to './data/key.pem'
-----
docker build -t deepl-tg .
```

##### Run:
- create bot.json in data folder [[example](./data/example.json)]
- create and run a container named deepl-tg
```sh
make docker-run
# or
docker run --name deepl-tg -d -p 8443:8443 deepl-tg
```

##### Usage:
Just start or stop the container.
```sh
make start
# or
make stop
```

##### Without docker:
Make sure golang is installed.
```sh
go get -v github.com/frzifus/deepL-tg
make amd64
# use:
./deepL-tg PATH_TO_CONFIG # by default ./data/
```
