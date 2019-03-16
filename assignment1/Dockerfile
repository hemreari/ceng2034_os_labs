FROM golang:1.9.7
WORKDIR /home/hemreari
ADD ./main.go .
RUN apt-get update && apt-get install sudo -y
RUN go get golang.org/x/sys/unix
RUN go build -o assignment1
CMD ["./assignment1"]