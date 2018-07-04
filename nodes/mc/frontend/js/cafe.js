window.onload = function(){
    var app = new Vue({
	el: '#app',
	data() {
	    return {
		mode: 'nodes',
		detail_node: null,
		nodes: {},
		types: {},
		emps: {},
		routes: {},
		id_to_name: {},
		error: '',
		network: false,
		mango: null,
		cmd: "",
		cmd_history_index: -1,
		cmd_result: "",
		cmd_history: [],
		commands: {
		    "help":["cmd"],
		    "addroute":["spec"],
		    "addnode":["type", "name", "group"],
		    "delnode":["name"],
		    "delroute":["src","dst","group"],
		    "startemp":["name", "group"]
		}
	    };
	},
	mounted: function(){
	    var self = this;
            self.mango = new Mango({"info":self.info_cb,"success":self.success});
	    console.log(self.mango);
	    var x = 0;
	    var init_group = function(){
		self.mango.m_send("addgroup",{"name":"frontend"});
		self.query();
		clearInterval(x);
	    }
	    x = setInterval(init_group,1000);
	    setInterval(self.query,10000);
	},
	methods: {
	    sorted_nodes: function(){
		ans = Object.keys(this.nodes);
		ans.sort();
		return ans;
	    },
	    graph_data: function(event) {
		var ns = [], es = [], id = 0;
		var name_to_id = {};
		this.id_to_name = {};
		for(var n in this.nodes){
		    ns.push({"id":id, "label":n, "mass":2});
		    name_to_id[n] = id;
		    this.id_to_name[id] = n;
		    id++;
		}
		console.log("N2I",name_to_id);
		for(var n in this.nodes){
		    for(var r in this.nodes[n].routes){
			var rt = this.nodes[n].routes[r];
			if(rt.edits) {
			    console.log("EDIT NODE",rt.edits,id);
			    ns.push({"id":id, "label":rt.edits, "shape":"dot","size":"5","mass":2});
			    es.push({"from":name_to_id[rt.src],"to":id});
			    es.push({"from":id,"to":name_to_id[rt.dst]});
			    id++;
			}
			else{
			    es.push({"from":name_to_id[rt.src],"to":name_to_id[rt.dst]});
			}
		    }
		}
		return {"nodes":new vis.DataSet(ns), "edges": new vis.DataSet(es)};
	    },
	    set_mode: function(mode, event){
		this.mode = mode;
		if(mode == 'graph'){
		    this.update_graph();
		}
		else {
		    if(this.network){
			this.network.destroy();
			this.network = false;
		    }
		}
		if(mode == 'nodes'){
		    this.query();
		}
	    },
	    update_graph: function(){
		if(this.mode != 'graph') return;
		var self = this;
		Vue.nextTick(function() {
		    var container = document.getElementById('graph');
		    self.network = new vis.Network(container,
						   self.graph_data(),
						   {
						       height: '100%',
						       width: '100%',
						       edges: { arrows: "to", font: { size: 12 } },
						       nodes: { shape: 'box', margin: 10 },
						       layout: { hierarchical: { direction: "UD", sortMethod: "directed", parentCentralization: false } },
						       physics: { enabled: false }
						   });
		    
		    self.network.on("doubleClick", function(params){
		    	if(params.nodes.length > 0 && params.nodes[0] in self.id_to_name){
		    	    self.details(self.id_to_name[params.nodes[0]]);
		    	}
		    });
		});
	    },
	    query: function(){
		this.mango.m_send("query",{"nodes":"","routes":"","types":"","emps":""});
	    },
	    details: function(name, event){
		this.detail_node = name;
		this.set_mode('detail');
	    },
	    startemp: function(args){
		this.mango.m_send("startemp",{"name":args.name,"group":args.group});
	    },
	    addnode: function(args){
		this.mango.m_send("addnode",{"name":args.name,"node_type":args.type,"group":args.group});
	    },
	    addroute: function(args){
		this.mango.m_send("addroute",{"group":args.group, "spec":args.spec});
	    },
	    delroute: function(args){
		this.mango.m_send("delroute",{"src":args.src,"dst":args.dst,"group":args.group});
	    },
	    delnode: function(args){
		this.mango.m_send("delnode",{"node":args.name});
	    },
	    editroute: function(code){
		this.cmd = "addroute " + code;
	    },
	    populatecmd: function(input){
		var tokens = input.split(" ");
		var c = tokens[0].trim()
		if(c in this.commands){
		    var start = input.length;
		    var runnable = tokens.length - 1 >= this.commands[c].length;
		    for(var i = tokens.length - 1; i < this.commands[c].length; i++){
			input += " <" + this.commands[c][i] + ">";
		    }
		    this.cmd = input;
		    if(runnable) this.runcommand();
		    else{
			Vue.nextTick(function() {
			    document.getElementById("cmd_input").select();
			    document.getElementById("cmd_input").setSelectionRange(start+1,input.length);
			});
		    }
		}
	    },
	    help: function(args){
		if(args && args.cmd in this.commands){
		    this.cmd_result = args.cmd;
		    for(var i = 0; i < this.commands[args.cmd].length; i++){
			this.cmd_result += " <" + this.commands[args.cmd][i] + ">";
		    }
		}
		else{
		    this.cmd_result = "";
		    for(var cmd in this.commands){
			this.cmd_result += cmd;
			for(var i = 0; i < this.commands[cmd].length; i++){
			    this.cmd_result += " <" + this.commands[cmd][i] + ">";
			}
			this.cmd_result += "\n";
		    }
		}
	    },
	    addhistory: function(){
		this.cmd_history.push(this.cmd);
		this.cmd_history_index = this.cmd_history.length;
		this.cmd = "";
	    },
	    runcommand: function(){
		var tokens = this.cmd.split(" ");
		console.log("RC",this.cmd,tokens.length);
		var c = tokens[0].trim();
		tokens = tokens.slice(1);
		
		if(c in this.commands) {
		    var nargs = this.commands[c].length
		    console.log(nargs,tokens.length);
		    if(nargs > tokens.length) {
			this.addhistory();
			this.help(c);
			return;
		    }
		    if(nargs < tokens.length){
			tokens[nargs-1] = tokens.slice(nargs-1).join(" ")
			tokens = tokens.slice(0,nargs)
		    }
		    var args = {}
		    for(var i = 0; i < nargs; i++){
			args[this.commands[c][i]] = tokens[i].trim();
		    }
		    console.log("RUN",c,args);
		    this.addhistory();
		    this[c](args);
		}
		else{
		    this.addhistory();
		    this.help();
		}
	    },
	    cmd_up: function(){
		console.log("CHI",this.cmd_history_index);
		if(this.cmd_history_index-1 < 0) return;
		this.cmd = this.cmd_history[--this.cmd_history_index];
	    },
	    cmd_down: function(){
		console.log("CHI",this.cmd_history_index);
		if(this.cmd_history_index > this.cmd_history.length-1) return;
		else if(this.cmd_history_index == this.cmd_history.length-1){
		    this.cmd_history_index++;
		    this.cmd = "";
		}
		else this.cmd = this.cmd_history[++this.cmd_history_index];
	    },
	    info_cb: function(header, args, event){
		this.nodes = {};
		var ns = args['nodes'];
		for(var i in ns){
		    var name = ns[i]['group'] + '/' + ns[i]['name']
		    this.nodes[name] = ns[i];
		    this.nodes[name].routes = [];
		    this.nodes[name]['interface'] = JSON.parse(ns[i]['interface']) || {};
		    this.detail_node = this.detail_node || name;
		}

		var rs = args['routes']
		for(var i in rs){
		    var src_name = rs[i]['src']['group']+'/'+rs[i]['src']['name'];
		    var dst_name = rs[i]['dst']['group']+'/'+rs[i]['dst']['name'];
		    this.nodes[src_name].routes.push({"src":src_name,"dst":dst_name,"group":rs[i].group,"edits":rs[i].edits,"name":rs[i].name});
		}

		this.types = args['types'];
		this.emps = args['emps'];
		
		for(var n in this.nodes){
		    console.log("N",n,this.nodes[n]);
		}
	    },
	    success: function(header,args,event){
		if(!args.success) this.error = args.message;
		else {
		    this.error = "success";
		    this.query();
		}
	    }
	}
    });
}
