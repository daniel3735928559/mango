window.onload = function(){
    
    Vue.component('simple-shell', {
	template: `<div class="command_entry">
			<div class="command_history">{{cmd_result}}</div>
			<div class="command_modeline">{{connection}}: </div>
			<input type="text" class="cmd_input" v-bind:value="value" v-on:input="$emit('input', $event.target.value)" v-on:keyup.up="cmd_up" v-on:keyup.down="cmd_down" v-on:keyup.enter="runcommand" />
		   </div>`,
	data() {
	    return {
		cmd: "",
		connection: "",
		cmd_history_index: -1,
		cmd_result: "",
		cmd_history: []
	    }
	},
	props: ["commands", "run_cb", "value"],
	methods: {
	    populate: function(c, args, run){
		var arglist = [c];
		for(var a in args){
		    arglist.push('-'+a);
		    arglist.push(args[a]);
		}
		this.cmd = shellquote.quote(arglist);
		if(run) this.runcmd();
		else Vue.nextTick(function() {
		    document.getElementById("cmd_input").select();
		    document.getElementById("cmd_input").setSelectionRange(start+1,input.length);
		});
	    },
	    addhistory: function(){
		this.cmd_history.push(this.cmd);
		this.cmd_history_index = this.cmd_history.length;
		this.cmd = "";
	    },
	    runcmd: function(){
		var tokens = shellquote.parse(this.cmd);
		var c = tokens[0];
		tokens = tokens.slice(1);
		var args = {};
		while(tokens.length > 1){
		    var arg = tokens[0];
		    if(arg[0] != '-'){
			return;
		    }
		    arg = arg.slice(1);
		    args[arg] = tokens[1];
		    tokens = tokens.slice(2);
		}
		if(tokens.length > 0){
		    return;
		}
		console.log("RUN",c,args);
		this.addhistory();
		if(c == "help") this.help(args);
		else this.run_cb(c, args);
	    },
	    cmd_up: function(){
		if(this.cmd_history_index-1 < 0) return;
		this.cmd = this.cmd_history[--this.cmd_history_index];
	    },
	    cmd_down: function(){
		if(this.cmd_history_index > this.cmd_history.length-1) return;
		else if(this.cmd_history_index == this.cmd_history.length-1){
		    this.cmd_history_index++;
		    this.cmd = "";
		}
		else this.cmd = this.cmd_history[++this.cmd_history_index];
	    }	    
	}
    });

    
    var app = new Vue({
	el: '#app',
	data() {
	    return {
		mode: 'nodes',
		detail_node: null,
		nodes: {},
		groups: {},
		types: {},
		emps: {},
		routes: {},
		id_to_name: {},
		error: '',
		network: false,
		mango: null,
		cmd: "hello",
		shell_commands: {
		    "help":["cmd"],
		    "test_cmd":["test_arg"],
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
	    //setInterval(self.query,10000);
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
	    connect: function(node){
		
		this.mango.m_send("addroute",{"group":"frontend", "spec":"system/mx_fe > fe forward {env.name = forward_cmd; del forward_cmd;} > " + node});
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
	    info_cb: function(header, args, event){
		this.nodes = {};
		this.groups = {};
		var ns = args['nodes'];
		for(var i in ns){
		    var group = ns[i]['group'];
		    var nid = ns[i]['name'];
		    console.log("G",group);
		    if(group in this.groups) this.groups[group].push(nid);
		    else this.groups[group] = [nid];
		    
		    var name = group + '/' + nid;
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
	    run_cmd: function(cmd, args){
		return this[cmd](args);
	    },
	    test_cmd: function(args){
		console.log("TEST",args);
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
