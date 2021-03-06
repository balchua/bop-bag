FROM balchu/dqlite-base:1.9.0
ENV CGO_LDFLAGS_ALLOW="-Wl,-z,now"
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY cmd/ cmd/
COPY pkg/ pkg/
COPY *.go ./

RUN go build .
#RUN apt-get update && apt-get install -y iputils-ping dnsutils
COPY runbopbag.sh /app
RUN chmod +x runbopbag.sh

ENTRYPOINT ["/app/runbopbag.sh"]