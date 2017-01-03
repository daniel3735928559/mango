zmq = require('zmq');
jsyaml = require('js-yaml');
fs = require('fs');

function MNode(debug){
    var self = this;
    this.version = "0.1";
    this.debug = debug;
    this.serialiser = new Serialiser(this.version);
    this.iface = new Interface();
    this.node_id = process.env['MANGO_ID'];
    this.ports = [];
    var server = process.env['MC_ADDR'];
    
    this.run = function(){
	var ifce = {};
	for(var i in this.iface.iface){
	    ifce[i] = JSON.parse(JSON.stringify(this.iface.iface[i]));
	    delete ifce[i]['handler'];
	}
	console.log("IF",ifce)
	self.m_send('hello',{'id':self.node_id,'if':ifce,'ports':self.ports},"mc")
    }
    
    this.dispatch = function(header,args){
	console.log("DISPATCH",header,args,self.iface.iface,self.iface.iface[header['command']]);
	try{
            result = self.iface.iface[header['command']]['handler'](header,args);
            if(result){
		self.m_send("reply", result, header['port'])
	    }
	} catch(e) {
	    console.log(e);
            self.handle_error(header['src_node'],e+"")
	}
    }
    
    this.handle_error = function(src,err){
	console.log('OOPS',src,err);
	self.m_send('error',{'source':src,'message':err},"mc");
    }

    this.reg = function(header,args){
	if(header.src_node == "mc") self.node_id = args["id"];
    }

    this.reply = function(header,args){
	console.log("REPLY",header,args);
    }

    this.heartbeat = function(header,args){
	self.m_send("alive",{},"mc");
    }

    this.make_header = function(command,src_port){
	if(!src_port) src_port = "stdio";
	return {'src_port':src_port, 'command':command};
    }

    this.m_send = function(command,msg,port){
	console.log('sending',command,msg,port)
	header = self.make_header(command,port)
	self.dataflow.send(header,msg)
    }
    
    this.iface.add_interface(process.env['NODE_PATH']+'/../node.yaml',{
        'reg':self.reg, 'reply':self.reply, 'heartbeat':self.heartbeat
    });
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
    
    this.add_interface = function(if_file,handlers){
	try {
	    var spec = jsyaml.safeLoad(fs.readFileSync(if_file, 'utf8'));
	    console.log(JSON.stringify(spec));
	    var name = spec.name;
	    var missing = [], extra = [];
	    for(var i in handlers)
		if(!(i in spec.inputs)) extra.push(i);
	    for(var i in spec.inputs)
		if(!(i in handlers)) missing.push(i);
	    if(extra.length > 0) throw new MError("Functions not in interface: "+extra.join(", "));
	    if(missing.length > 0) throw new MError("Functions not implemented: "+missing.join(", "));

	    self.iface[name] = {};
	    for(var i in handlers){
		self.iface[name][i] = spec.inputs[i] ? spec.inputs[i] : {};
		self.iface[name][i]['handler'] = handlers[i];
	    }
	} catch (e) {
	    console.log(e);
	    throw new MError("Failed to load interface");
	}
    }

    this.get_function = function(fn,ns){
	if(ns) return self.iface[ns][fn];
	for(var n in self.iface) if(fn in self.iface[n]) return self.iface[n][fn];
	throw new MError("Functions not implemented: "+ns+"."+fn);	
    }
    
    this.validate = function(fn){
	return fn in self.iface;
    }
}

function Transport(target){
    var self = this;
    this.target = target
    this.socket = zmq.socket("dealer");
    this.socket.connect(target);
    
    this.tx = function(data){
	self.socket.send(data);
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
	    if(!self.iface.validate(m[0]["command"])) throw new MError("Unknown function");
	    console.log(self.dispatch_cb);
	    self.dispatch_cb(m[0],m[1]);
	} catch(e) {
	    console.log(e);
	    self.error_cb(e);
	}
    }

}

function Serialiser(version){
    var self = this;
    this.version = version;
    this.method = "json";
    
    this.make_preamble = function(){
	return "MANGO"+self.version+" json\n";
    }

    this.parse_preamble = function(data){
	console.log("GOT",data);
	var nl1 = data.indexOf('\n');
	var m = data.substring(0,nl1).match(/^MANGO([0-9.]*) ([^ ]*)$/);
	if(!m || m.length < 1) throw new MError("Preamble failed to parse");
	return [m[1],m[2],data.substring(nl1)];
    }

    this.serialise = function(header,args){
	return self.make_preamble()+JSON.stringify({"header":header,"args":args});
    }
    
    this.deserialise = function(data){
	data = data.toString('utf8');
	p = self.parse_preamble(data);
	var ver = p[0]; var method = p[1]; var message = p[2];
	if(ver != self.version) throw new MError("Version mismatch");
	if(method != self.method) throw new MError("Unsupported method");
	try{
	    var d = JSON.parse(message);
	    return [d['header'],d['args']]
	} catch (e) {
	    console.log(e);
	    throw new MError("Failed to parse message");
	}
    }
}

module.exports = MNode;
