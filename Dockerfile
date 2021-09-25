FROM localhost:32000/dqlite-base:1.0.0
ENV CGO_LDFLAGS_ALLOW="-Wl,-z,now"
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY cmd/ cmd/
COPY *.go ./

RUN go build .
