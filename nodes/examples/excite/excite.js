MNode = require('libmango')

function Excite(){
    this.node = new MNode();
    this.node.iface.add_interface('/home/zoom/suit/mango/nodes/examples/excite/excite.yaml',
				  {'excite':this.excite,'print':this.print});
    this.node.ready();
}

Excite.prototype.excite = function(header,args){
    return {'excited':args['str']+'!'}
}

Excite.prototype.print = function(header,args){
    console.log("PRINT",header,args);
}

var ex = new Excite();
