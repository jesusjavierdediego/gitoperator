FROM golang:alpine as golang
RUN apk add --no-cache git
WORKDIR $GOPATH/src/xqledger/gitoperator
COPY . ./
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /usr/bin/dep
ADD resources/application.yml ./
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
RUN zip -r -0 /zoneinfo.zip .

FROM scratch
COPY --from=golang /go/bin/gitoperator /app
COPY --from=golang /go/src/xqledger/gitoperator/resources/application.yml ./
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
RUN git config --global user.name 'TestOrchestrator'
RUN git config --global user.email 'TestOrchestrator@gmail.com'
RUN git config --global user.signingkey '21BE044483C263AC!'
ENTRYPOINT ["/app"]