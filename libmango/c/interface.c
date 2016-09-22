
function Interface(){
    var self = this;
    this.iface = {};
    
    this.add_interface = function(if_file,handlers){
	try {
	    var spec = jsyaml.safeLoad(fs.readFileSync(if_file, 'utf8'));
	    console.log(JSON.stringify(spec));
	    var missing = [], extra = [], existing = [];
	    for(var i in spec)
		if(!(i in handlers)) missing.push(i);
	    for(var i in handlers)
		if(!(i in spec)) extra.push(i);
	    for(var i in spec)
		if(i in self.iface) existing.push(i);
	    if(missing.length > 0) throw new MError("Functions not implemented: "+missing.join(", "));
	    if(extra.length > 0) throw new MError("Functions not in interface: "+extra.join(", "));
	    if(existing.length > 0) throw new MError("Functions already implemented: "+existing.join(", "));

	    for(var i in handlers){
		self.iface[i] = spec[i] ? spec[i] : {};
		self.iface[i]['handler'] = handlers[i];
	    }
	} catch (e) {
	    console.log(e);
	    throw new MError("Failed to load interface");
	}
    }

    this.validate = function(fn){
	return fn in self.iface;
    }
}
