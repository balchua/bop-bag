FROM localhost:32000/dqlite-base:1.9.0
ENV CGO_LDFLAGS_ALLOW="-Wl,-z,now"
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY *.go ./

RUN go build .
# RUN apt-get update && apt-get install -y iputils-ping dnsutils
COPY runbopbag.sh /app
ENTRYPOINT bash runbopbag.sh