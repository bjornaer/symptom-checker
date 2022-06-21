FROM golang:1.18-alpine as builder

ENV APP_NAME sympton-checker
ENV CMD_PATH cmd/SymptomChecker/

WORKDIR $GOPATH/src/$APP_NAME
COPY go.mod go.sum ./
RUN go mod download 
COPY cmd ./cmd
COPY internal ./internal
COPY ent ./ent

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

FROM alpine:3.16

ENV APP_NAME sympton-checker

COPY --from=builder /$APP_NAME .

EXPOSE 8081

CMD ./$APP_NAME
