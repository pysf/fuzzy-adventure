FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go test -v
RUN go build -o /rsc

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /rsc /rsc

ENTRYPOINT ["/rsc"]