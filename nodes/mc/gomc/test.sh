#GOPATH=$(realpath .) go build mc
if [[ "$1" == "verbose" ]]; then
  GOPATH=$(realpath .) go test -coverprofile=cover.out -v mc/router
else
  GOPATH=$(realpath .) go test -coverprofile=cover.out mc/router
fi
#./bin/goyacc -o src/mc/routeparser/pt/pt.go -v y src/mc/routeparser/pt/pt.go.y
#GOPATH=$(realpath .) go test -v mc/routeparser/pt
