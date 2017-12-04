BINARY = deepL-tg

.PHONY: clean arm amd64 docker-run docker-stop docker-build cert

# Build the project
all: cert docker-build

amd64:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY} -v
arm:
	GOOS=linux GOARCH=arm go build ${LDFLAGS} -o ${BINARY} -v
docker-build:
	docker build -t deepl-tg .
docker-run:
	@echo "Run container"
	docker run --name deepl-tg -d -p 8443:8443 deepl-tg
start:
	@echo "Start container"
	docker start deepl-tg
stop:
	@echo "Stop container"
	docker stop deepl-tg
cert:
	mkdir -p ./data
	@read -p "Enter Domain Name:" domain; \
	openssl req -newkey rsa:2048 -sha256 -nodes -keyout ./data/key.pem -x509 -days 365 \
	-out ./data/cert.pem -subj "/C=DE/ST=Aa/L=Brooklyn/O=A Company/CN=$$domain"
clean:
	rm -rf ./data/*


