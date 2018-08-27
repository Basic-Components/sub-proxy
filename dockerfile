FROM hsz1273327/golang-zmq
ENV GOPATH="/app"
RUN go get github.com/pebbe/zmq4
RUN go get github.com/sirupsen/logrus