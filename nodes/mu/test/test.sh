DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo AAAAAAAAAAAAAAAAAA "$DIR"
mc launch -node mu -id moo -env '{"MU_WS_PORT":"9090","MU_ROOT_DIR":"'"$DIR"'","MU_HTTP_PORT":"9999","MU_IF":"'"$DIR"'/test.yaml"}'
mc launch -node excite -id ex
sleep 1
mc route -map 'moo > ex > +{"command":"dostuff"} > moo'
