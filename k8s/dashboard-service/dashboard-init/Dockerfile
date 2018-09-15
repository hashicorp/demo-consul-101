FROM alpine:3.8

RUN apk add --no-cache curl

WORKDIR /app
ADD . /app
CMD ["./register-consul-service.sh"]
