./bin/goyacc -o src/mc/route/parser.go -p RouteParser -v src/mc/route/parser.output src/mc/route/parser.go.y
./bin/goyacc -o src/mc/valuetype/parser.go -p ValueTypeParser -v src/mc/valuetype/parser.output src/mc/valuetype/parser.go.y
./bin/goyacc -o src/mc/value/parser.go -p ValueParser -v src/mc/value/parser.output src/mc/value/parser.go.y
GOPATH=$(realpath .) go install mc
GOPATH=$(realpath .) go install nodes/excite
GOPATH=$(realpath .) go install nodes/log
GOPATH=$(realpath .) go install nodes/tester
GOPATH=$(realpath .) go install nodes/notify
GOPATH=$(realpath .) go install nodes/mx/agent
GOPATH=$(realpath .) go install nodes/mx
GOPATH=$(realpath .) go install nodes/imap
GOPATH=$(realpath .) go install nodes/smtp
GOPATH=$(realpath .) go install nodes/sigjam
