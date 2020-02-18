#GOPATH=$(realpath .) go build mc
GOPATH=$(realpath .) go test -v mc/router

#./bin/goyacc -o src/mc/routeparser/pt/pt.go -v y src/mc/routeparser/pt/pt.go.y
#GOPATH=$(realpath .) go test -v mc/routeparser/pt
