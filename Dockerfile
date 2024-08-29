FROM alpine:latest

COPY entrypoint.sh /
COPY mantrae /usr/local/bin/mantrae

WORKDIR /app
EXPOSE 3000
ENTRYPOINT ["/entrypoint.sh"]
CMD ["mantrae"]
