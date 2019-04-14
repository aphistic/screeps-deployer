# Build Image
FROM golang:1.12-alpine3.9 AS builder
WORKDIR /project
COPY go.mod .
COPY go.sum .
RUN apk add --no-cache git
RUN go mod download
COPY . .
RUN cd cmd/screeps-deployer && \
    go build

# Action Image
FROM alpine:3.9

LABEL "com.github.actions.name"="Screeps Deployer"
LABEL "com.github.actions.description"="Deploy code to the game Screeps"
LABEL "com.github.actions.icon"="upload"
LABEL "com.github.actions.color"="blue"

LABEL "repository"="https://github.com/aphistic/screeps-deployer"
LABEL "homepage"="https://github.com/aphistic/screeps-deployer"
LABEL "maintainer"="Erik Davidson <erik@erikd.org>"

COPY --from=builder /project/cmd/screeps-deployer /screeps-deployer
RUN chmod 755 /screeps-deployer
ENTRYPOINT ["/screeps-deployer"]