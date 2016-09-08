var Mango = function(port, commands){
    console.log("Hello mango world",port);
    this.sock = new WebSocket("ws://localhost:"+port+"/");
    var self = this;
    this.sock.onmessage = function(e) {
        console.log('got: ' + e.data);
	payload = e.data.substring(e.data.indexOf('\n')+1);
	data = JSON.parse(payload);
	self.m_recv(data.header, data.args);
    };
    this.commands = commands;
}

Mango.prototype.m_send = function(command,args){
    console.log("A",JSON.stringify(command,args));
    var msg = "MANGO0.1 json\n"+JSON.stringify({"header":{"command":command},"args":args});
    this.sock.send(msg);
}

Mango.prototype.m_recv = function(header,args){
    if('command' in header && header['command'] in this.commands){
	commands[args['command']](header,args);
    }
    else{
	console.log('Bad command',JSON.stringify(header,args));
    }
}

console.log('asd');
