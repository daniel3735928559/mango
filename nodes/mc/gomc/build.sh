#GOPATH=$(realpath .) go build mc
./bin/goyacc -o src/mc/router/parser.go -v src/mc/router/parser.output src/mc/router/parser.go.y && \
GOPATH=$(realpath .) go build mc/router
#./bin/goyacc -o src/mc/routeparser/pt/pt.go -v y src/mc/routeparser/pt/pt.go.y
#GOPATH=$(realpath .) go test -v mc/routeparser/pt
