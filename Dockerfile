FROM golang:1.16-buster AS build

WORKDIR /go/src/github.com/bbchallenge/bbchallenge

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY lib_bbchallenge ./lib_bbchallenge
COPY main.go ./
RUN go build -o /bbchallenge


FROM gcr.io/distroless/base-debian11 AS runtime

WORKDIR /
COPY --from=build /bbchallenge /bbchallenge
USER nonroot:nonroot

ENTRYPOINT ["/bbchallenge"]
