var Mango = function(commands,default_handler){
    var xhr = new XMLHttpRequest();
    var self = this;
    self.sock = new ReconnectingWebSocket("ws://" + location.host + "/mangorx");
    self.txsock = new ReconnectingWebSocket("ws://"  + location.host + '/mangotx')

    self.sock.onmessage = function(e) {
	var reader = new FileReader();
	reader.addEventListener("loadend", function() {
	    console.log("stuff: ",reader.result)
	    var byte_array = new Uint8Array(reader.result);
	    payload_bin = byte_array.slice(byte_array.indexOf(10)+1);
	    payload = String.fromCharCode.apply(null, payload_bin);
	    data = JSON.parse(payload);
	    console.log("RECVD",JSON.stringify(data));
	    self.m_recv(data.header, data.args)
	});
	reader.readAsArrayBuffer(e.data);
    };
    self.commands = commands;
    self.default_handler = default_handler;
}

Mango.prototype.m_send = function(command,args){
    console.log("TX",JSON.stringify(command,args));
    var msg = "MANGO0.1 json\n"+JSON.stringify({"header":{"name":command},"args":args});
    this.txsock.send(msg);
}

Mango.prototype.m_recv = function(header,args){
    if('name' in header && header['name'] in this.commands){
	this.commands[header['name']](header,args);
    }
    else if(this.default_handler){
	this.default_handler(header, args)
    }
    else{
	console.log('Unexpected command',JSON.stringify(header,args));
    }
}
