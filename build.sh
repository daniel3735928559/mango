./bin/goyacc -o src/mc/route/parser.go -p RouteParser -v src/mc/route/parser.output src/mc/route/parser.go.y && \
    ./bin/goyacc -o src/mc/valuetype/parser.go -p ValueTypeParser -v src/mc/valuetype/parser.output src/mc/valuetype/parser.go.y && \
    ./bin/goyacc -o src/mc/value/parser.go -p ValueParser -v src/mc/value/parser.output src/mc/value/parser.go.y && \
    echo hi
GOPATH=$(realpath .) go build mc
GOPATH=$(realpath .) go build nodes/excite
GOPATH=$(realpath .) go build nodes/log
GOPATH=$(realpath .) go build nodes/tester
GOPATH=$(realpath .) go build nodes/notify
GOPATH=$(realpath .) go build nodes/mx/agent
GOPATH=$(realpath .) go build nodes/mx
GOPATH=$(realpath .) go build nodes/imap
