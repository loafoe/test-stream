# syntax = docker/dockerfile:1-experimental

FROM --platform=${BUILDPLATFORM} golang:1.15.6-alpine AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.* .
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /app/server -ldflags "-X main.GitCommit=${GIT_COMMIT}" .

FROM golang:1.15.6-alpine
RUN apk add --no-cache tzdata jq
LABEL maintainer="Andy Lo-A-Foe <andy.lo-a-foe@philips.com>"
COPY --from=build /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]
