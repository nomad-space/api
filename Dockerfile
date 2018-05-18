FROM alpine
WORKDIR /
COPY ./bin/nomadapi /app/nomadapi
RUN apk --update add ca-certificates
ENV PORT 7784
EXPOSE $PORT
ENTRYPOINT /app/nomadapi