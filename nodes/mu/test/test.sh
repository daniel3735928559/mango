mc launch -node mu -id moo -env '{"MU_WS_PORT":"9090","MU_ROOT_DIR":"/home/zoom/suit/mango/nodes/mu/test/","MU_HTTP_PORT":"9999","MU_IF":"/home/zoom/suit/mango/nodes/mu/test/test.yaml"}'
sleep 1
mc launch -node excite -id ex
sleep 1
mc route -spec 'moo > ex > +{"command":"dostuff"} > moo'
