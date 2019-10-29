FROM alpine:3.10.3

RUN apk update && apk upgrade && apk add --update alpine-sdk && \
    apk add --no-cache bash git openssh make cmake
RUN make build

RUN mkdir /app
ADD ./out/g2-reverse-proxy /app/

CMD /app/g2-reverse-proxy