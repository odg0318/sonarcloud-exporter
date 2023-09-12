FROM golang:1.21-alpine AS build-stage

WORKDIR /app

COPY . .

RUN apk --no-cache add ca-certificates &&\
            CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sonarcloud-exporter ./cmd/sonarcloud-exporter

FROM alpine:3.18

RUN apk --no-cache add ca-certificates
COPY --from=build-stage /app/sonarcloud-exporter /usr/local/bin/sonarcloud-exporter

CMD ["sonarcloud-exporter"]
