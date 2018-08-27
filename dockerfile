FROM hsz1273327/sub-proxy:Base-v0
ADD . /app/src/github.com/Basic-Components/sub-proxy
ENV GOPATH="/app"
WORKDIR /app/src/github.com/Basic-Components/sub-proxy
RUN go build