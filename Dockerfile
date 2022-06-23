### build frontend ###
FROM node:16.15.0-alpine as feBuilder
WORKDIR /usr/src/app
ENV PATH /app/node_modules/.bin:$PATH
COPY frontend/package.json ./
COPY frontend/package-lock.json ./
RUN npm ci --silent
RUN npm install react-scripts@5.0.1 -g --silent
COPY frontend/ ./
RUN npm run build

### build backend ###
FROM golang:1.18 as beBuilder

ENV APP_NAME sympton-checker
ENV CMD_PATH cmd/SymptomChecker/

WORKDIR $GOPATH/src/$APP_NAME
COPY go.mod go.sum ./
RUN go mod download 
COPY cmd ./cmd
COPY internal ./internal
COPY ent ./ent

RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

### run built image ###
FROM alpine:3.16
# Create a group and user
RUN addgroup -S symptom && adduser -S symptom -G symptom

ARG ROOT_DIR=/home/symptom/app
 
WORKDIR ${ROOT_DIR}
 
RUN chown symptom:symptom ${ROOT_DIR}

ENV APP_NAME sympton-checker

# copy static assets file from frontend build
COPY --from=feBuilder --chown=symptom:symptom /usr/src/app/build ./frontend/build
COPY --from=beBuilder --chown=symptom:symptom /$APP_NAME .

USER symptom

EXPOSE 8081

CMD ./$APP_NAME
