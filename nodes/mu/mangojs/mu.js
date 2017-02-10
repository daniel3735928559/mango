var Mango = function(commands){
    var xhr = new XMLHttpRequest();
    var self = this;
    xhr.onreadystatechange = function() {
        if(xhr.readyState == XMLHttpRequest.DONE) {
           if (xhr.status == 200) {
               var port = parseInt(xhr.responseText);
	       console.log("Hello mango world",port);
	       self.sock = new WebSocket("ws://localhost:"+port+"/");
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
           }
            else if (xhr.status == 400) alert('There was an error 400');
            else alert('something else other than 200 was returned');
        }
    };

    xhr.open("GET", "/ws_port", true);
    xhr.send();
}

Mango.prototype.m_send = function(command,args){
    console.log("A",JSON.stringify(command,args));
    var msg = "MANGO0.1 json\n"+JSON.stringify({"header":{"name":command},"args":args});
    this.sock.send(msg);
}

Mango.prototype.m_recv = function(header,args){
    if('name' in header && header['name'] in this.commands){
	this.commands[header['name']](header,args);
    }
    else{
	console.log('Bad command',JSON.stringify(header,args));
    }
}
