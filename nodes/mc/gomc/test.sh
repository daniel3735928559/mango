#GOPATH=$(realpath .) go build mc
for pkg in route value valuetype; do
if [[ "$1" == "verbose" ]]; then
  GOPATH=$(realpath .) go test -coverprofile=cover.out -v mc/"$pkg"
elif [[ "$1" == "cover" ]]; then
  sed -i 's,^//.*,,' src/mc/"$pkg"/parser.go
  GOPATH=$(realpath .) go test -coverprofile=cover.out mc/"$pkg"
  GOPATH=$(realpath .) go tool cover -html cover.out -o coverage.html && firefox coverage.html
else 
  GOPATH=$(realpath .) go test -coverprofile=cover.out mc/"$pkg"
fi
done
#./bin/goyacc -o src/mc/routeparser/pt/pt.go -v y src/mc/routeparser/pt/pt.go.y
#GOPATH=$(realpath .) go test -v mc/routeparser/pt
