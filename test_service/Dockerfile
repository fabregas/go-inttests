FROM golang:latest 
RUN mkdir /app
ADD main.go /app/ 
WORKDIR /app 
RUN go build -o main . 
EXPOSE 5555
VOLUME /container_volume
CMD ["/app/main"]
