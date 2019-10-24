FROM alpine:3.10.3

RUN mkdir -p /app/bin && mkdir -p /app/logs
ADD ./out/g2-reverse-proxy /app/bin

CMD /app/bin/g2-reverse-proxy