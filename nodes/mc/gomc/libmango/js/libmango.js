zmq = require('zeromq');
fs = require('fs');

function MNode(debug){
    var self = this;
    this.version = "0.1";
    this.debug = debug;
    this.serialiser = new Serialiser(this.version);
    this.iface = new Interface();
    this.node_id = process.env['MANGO_COOKIE'];
    var server = process.env['MANGO_SERVER'];
    
    this.run = function(){
    }
    
    this.dispatch = function(header,args){
	console.log("H",header)
	console.log("DISPATCH",header,args,self.iface.iface,self.iface.get_function(header['command']));
	try{
            result = self.iface.get_function(header['command'])(header,args);
            if(result){
		self.m_send(result[0], result[1], header.mid)
	    }
	} catch(e) {
	    console.log("ONO",e);
            self.handle_error(header,e+"")
	}
    }
    
    this.handle_error = function(src,err){
	console.log('OOPS',src,err);
	self.m_send('error',{'source':JSON.stringify(src),'message':err},"mc");
    }

    this.heartbeat = function(header,args){
	self.m_send("alive",{},undefined,"system");
    }

    this.make_header = function(name,mid,type){
	ans = {'command':name,
	       'cookie':this.node_id,
	       'mid':mid ? mid : RandomId(),
	       'format':'json'};
	return ans
    }

    this.exit = function(header,args){
	process.exit();
    }

    this.m_send = function(name,msg,mid,type){
	console.log('sending',name,msg,mid)
	header = self.make_header(name,mid,type)
	self.dataflow.send(header,msg)
    }
    
    this.iface.add_handlers({'heartbeat':self.heartbeat,'exit':self.exit});
    this.local_gateway = new Transport(server);
    var s = this.local_gateway.socket;
    this.dataflow = new Dataflow(self.iface, self.local_gateway, self.serialiser, self.dispatch, self.handle_error);
    s.on('message',function(data){ self.dataflow.recv(data); });
}

function MError(message){
    this.message = message;
    this.name = "Mango Error";
}

function Interface(){
    var self = this;
    this.iface = {};
    
    this.add_handlers = function(handlers){
	for(var i in handlers){
	    self.iface[i] = handlers[i];
	}
    }

    this.get_function = function(fn){
	if(fn in self.iface) return self.iface[fn];
	throw new MError("Function not implemented: "+fn);
    }
    
    this.validate = function(fn){
	return self.get_function(fn);
    }
}

function RandomId() {
    var result = '';
    var cs = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_';
    var numcs = cs.length;
    for ( var i = 0; i < numcs; i++ ) {
	result += cs.charAt(Math.floor(Math.random() * (numcs-1)));
    }
    return result;
}

function Transport(target){
    var self = this;
    this.target = target
    this.route = RandomId();
    this.socket = zmq.socket("dealer");
    this.socket.identity = this.route
    this.socket.connect(target);
    
    this.tx = function(data){
	self.socket.send(['',data]);
    }

}

function Dataflow(iface,transport,serialiser,dispatch_cb,error_cb){
    var self = this;
    this.iface = iface;
    this.transport = transport;
    this.serialiser = serialiser;
    this.dispatch_cb = dispatch_cb;
    this.error_cb = error_cb;
    console.log(this.dispatch_cb);
    this.send = function(header,args){
	self.transport.tx(self.serialiser.serialise(header,args));
    }
    
    this.recv = function(data){
	try{
	    var m = self.serialiser.deserialise(data);
	    console.log("MMM",m,m[0],m[1])
	    if(!self.iface.validate(m[0]["command"])) throw new MError("Unknown function");
	    console.log(self.dispatch_cb);
	    self.dispatch_cb(m[0],m[1]);
	} catch(e) {
	    console.log("ONBO",e);
	    self.error_cb(e);
	}
    }

}

function Serialiser(version){
    var self = this;
    this.version = version;
    this.method = "json";
    
    // this.make_preamble = function(){
    // 	return "MANGO"+self.version+" json\n";
    // }

    // this.parse_preamble = function(data){
    // 	console.log("GOT",data);
    // 	var nl1 = data.indexOf('\n');
    // 	var m = data.substring(0,nl1).match(/^MANGO([0-9.]*) ([^ ]*)$/);
    // 	if(!m || m.length < 1) throw new MError("Preamble failed to parse");
    // 	return [m[1],m[2],data.substring(nl1)];
    // }

    this.serialise = function(header,args){
	return JSON.stringify(header) + "\n" + JSON.stringify(args);
	//return self.make_preamble()+JSON.stringify({"header":header,"args":args});
    }
    
    this.deserialise = function(data){
	data = data.toString('utf8');
	idx = data.indexOf("\n")
	header_str = data.substring(0, idx);
	body_str = data.substring(idx+1);
	try {
	    return [JSON.parse(header_str),JSON.parse(body_str)];
	} catch(e) {
	    console.log(e);
	    throw new MError("Failed to parse message:"+header_str+"\n"+body_str);
	}
	// p = self.parse_preamble(data);
	// var ver = p[0]; var method = p[1]; var message = p[2];
	// if(ver != self.version) throw new MError("Version mismatch");
	// if(method != self.method) throw new MError("Unsupported method");
	// try{
	//     var d = JSON.parse(message);
	//     return [d['header'],d['args']]
	// } catch (e) {
	//     console.log(e);
	//     throw new MError("Failed to parse message");
	// }
    }
}

module.exports = MNode;
