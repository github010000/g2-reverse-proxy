FROM alpine:3.10.3

RUN mkdir /app
ADD ./out/g2-reverse-proxy /app/

CMD /app/g2-reverse-proxy