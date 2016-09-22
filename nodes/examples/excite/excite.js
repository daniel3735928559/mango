libmango = require('libmango')

function Excite(){
    process.env['MANGO_ID'] = 'exc';
    process.env['MC_ADDR'] = 'tcp://localhost:61453';
    this.node = new libmango.MNode();
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
