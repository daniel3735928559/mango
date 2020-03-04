./bin/goyacc -o src/mc/route/parser.go -p RouteParser -v src/mc/route/parser.output src/mc/route/parser.go.y && \
./bin/goyacc -o src/mc/valuetype/parser.go -p ValueTypeParser -v src/mc/valuetype/parser.output src/mc/valuetype/parser.go.y && \
GOPATH=$(realpath .) go build mc
#GOPATH=$(realpath .) go build mc/router
#./bin/goyacc -o src/mc/routeparser/pt/pt.go -v y src/mc/routeparser/pt/pt.go.y
#GOPATH=$(realpath .) go test -v mc/routeparser/pt
