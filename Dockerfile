# Build image
FROM golang:1.17-alpine AS build-env
WORKDIR /root/
COPY ./ ./
RUN apk update \
  && apk add git \
  && CGO_ENABLED=0 GOOS=linux go build -a

# Runtime image
FROM alpine:3.15 AS runtime-env
WORKDIR /root/
COPY --from=build-env /root/ican-tv ./
COPY ./config.json ./config.json
COPY ./index.gohtml ./index.gohtml
RUN apk update \
  && apk add ca-certificates curl bash tzdata
CMD ["./ican-tv"]
