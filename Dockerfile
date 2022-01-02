FROM golang:1.16 AS base
ENV CGO_ENABLED=0

WORKDIR /app
COPY go.mod /app
RUN go mod download

COPY . /app
RUN go build -o server ./cmd

FROM scratch
ENV APP_ENV=PRODUCTION

WORKDIR /app
COPY --from=base /app/server /app/server
ENTRYPOINT [ "./server" ]