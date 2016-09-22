
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
