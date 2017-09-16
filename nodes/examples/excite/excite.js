MNode = require('libmango')

function Excite(){
    this.node = new MNode();
    this.node.iface.add_interface('./excite.yaml',{'excite':this.excite});
    this.node.run();
}

Excite.prototype.excite = function(header,args){
    return ["excited",{'str':args['str']+'!'}]
}

new Excite();
