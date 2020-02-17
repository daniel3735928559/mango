GOPATH=$(realpath .) go build mc
./bin/goyacc -o src/mc/routeparser/routeparser.go -v y src/mc/routeparser/routeparser.go.y
GOPATH=$(realpath .) go test -v mc/routeparser

#./bin/goyacc -o src/mc/routeparser/pt/pt.go -v y src/mc/routeparser/pt/pt.go.y
#GOPATH=$(realpath .) go test -v mc/routeparser/pt
