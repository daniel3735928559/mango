MNode = require('libmango')

function Excite(){
    this.node = new MNode();
    this.node.iface.add_handlers({'excite':this.excite});
    this.node.run();
}

Excite.prototype.excite = function(header,args){
    return ["excited",{'message':args['message']+'!'}]
}

new Excite();
