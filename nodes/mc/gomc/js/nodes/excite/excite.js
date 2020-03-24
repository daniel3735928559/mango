MNode = require('libmango')

function Excite(){
    this.node = new MNode();
    this.node.iface.add_handlers({'excite':this.excite});
    this.node.run();
    this.node.m_send("excite",{})
}

Excite.prototype.excite = function(header,args){
    return ["excited",{'message':args['message']+'!'}]
}

console.log("How exciting")
new Excite();
console.log("Ready to excite")
