FROM golang:1.22 AS build-stage

WORKDIR /usr/local/src/snote

COPY go.mod ./
RUN go mod download

COPY . .

RUN echo "hello, world"
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/snote ./cmd/app

FROM build-stage AS release-stage

RUN useradd app
WORKDIR /usr/local/bin

COPY --from=build-stage /usr/local/bin/snote ./snote

EXPOSE 8080

USER app:app
ENTRYPOINT ["./snote", "--port", "8080"]