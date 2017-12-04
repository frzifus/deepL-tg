FROM golang

COPY . /go/src/github.com/frzifus/deepL-tg/

WorkDir /go/src/github.com/frzifus/deepL-tg/

RUN go get -v
RUN go install

VOLUME ["./data", "/data"]

EXPOSE 8433

CMD ["deepL-tg", "-path=/data/"]
