FROM ubuntu:20.04
WORKDIR /app
COPY Ginstart /app/ginstart
CMD ["/app/ginstart"]