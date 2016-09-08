var Mango = function(port, commands){
    console.log("Hello mango world",port);
    this.sock = new WebSocket("ws://localhost:"+port+"/");
    var self = this;
    this.sock.onmessage = function(e) {
	var reader = new FileReader();
	reader.addEventListener("loadend", function() {
	    console.log("stuff: ",reader.result)
	    var byte_array = new Uint8Array(reader.result);
	    
	    payload_bin = byte_array.slice(byte_array.indexOf(10)+1);
	    payload = String.fromCharCode.apply(null, payload_bin);
	    data = JSON.parse(payload);
	    self.m_recv(data.header, data.args)
	});
	reader.readAsArrayBuffer(e.data);
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
	this.commands[header['command']](header,args);
    }
    else{
	console.log('Bad command',JSON.stringify(header,args));
    }
}
