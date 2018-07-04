window.onload = function(){
    var app = new Vue({
	el: '#app',
	data() {
	    return {
		mode: 'nodes',
		detail_node: null,
		nodes: {},
		routes: {},
		error: '',
		network: false,
		mango: null,
		new_route: ""
	    };
	},
	mounted: function(){
	    var self = this;
            self.mango = new Mango({"info":self.info_cb,"success":self.success});
	    console.log(self.mango);
	    var x = 0;
	    var init_query = function(){
		self.mango.m_send("addgroup",{"name":"frontend"});
		self.query();
		clearInterval(x);
	    }
	    x = setInterval(init_query,1000);
	},
	methods: {
	    graph_data: function(event) {
		var ns = [], es = [], id = 0;
		var name_to_id = {};
		for(var n in this.nodes){
		    ns.push({"id":id, "label":n, "mass":2});
		    name_to_id[n] = id;
		    id++;
		}
		console.log("N2I",name_to_id);
		for(var n in this.nodes){
		    for(var r in this.nodes[n].routes){
			var rt = this.nodes[n].routes[r];
			if(rt.edits) {
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
						   {});
		    // self.network.on("doubleClick", function(params){
		    // 	if(params.nodes.length > 0){
		    // 	    self.open_snippet(params.nodes[0]);
		    // 	    self.set_mode('doc');
		    // 	}
		    // });
		});
	    },
	    query: function(){
		this.mango.m_send("query",{"nodes":"","routes":""});
	    },
	    details: function(name, event){
		this.detail_node = name;
		this.mode = 'detail';
	    },
	    addroute: function(){
		this.mango.m_send("addroute",{"group":"frontend", "spec":this.new_route});
	    },
	    delroute: function(src,dst,group){
		this.mango.m_send("delroute",{"src":src,"dst":dst,"group":group});
	    },
	    info_cb: function(header, args, event){
		this.nodes = {};
		var ns = args['nodes'];
		for(var i in ns){
		    var name = ns[i]['group'] + '/' + ns[i]['name']
		    this.nodes[name] = ns[i];
		    this.nodes[name].routes = [];
		}

		var rs = args['routes']
		for(var i in rs){
		    var src_name = rs[i]['src']['group']+'/'+rs[i]['src']['name'];
		    var dst_name = rs[i]['dst']['group']+'/'+rs[i]['dst']['name'];
		    this.nodes[src_name].routes.push({"src":src_name,"dst":dst_name,"group":rs[i].group,"edits":rs[i].edits});
		}
		for(var n in this.nodes){
		    console.log("N",n,this.nodes[n]);
		}
	    },
	    success: function(header,args,event){
		if(!args.success) this.error = args.message;
		else this.error = "success";
	    }
	}
    });
}
