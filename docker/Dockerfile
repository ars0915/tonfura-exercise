FROM golang:1.21-alpine3.18 as builder
ARG APP_NAME
RUN set -eux; \
	apk update && \
	apk add --no-cache sqlite
ENV GOPATH /go
ENV GO_WORKDIR $GOPATH/src/github.com/ars0915/${APP_NAME}/
WORKDIR $GO_WORKDIR
ADD . $GO_WORKDIR
RUN go build -o ${APP_NAME} -tags=jsoniter .

FROM alpine:3.18
ARG APP_NAME
COPY --from=builder /go/src/github.com/ars0915/${APP_NAME}/${APP_NAME} .
ENV APP_NAME=${APP_NAME}
CMD ./$APP_NAME
EXPOSE 8080