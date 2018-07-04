var intersect_lists = function(l1,l2){
    var result = [];
    for(var i = 0; i < l1.length; i++){
	for(var j = 0; j < l2.length; j++){
	    if(l1[i] == l2[j]) result.push(l1[i]);
	}
    }
    return result;
}

var union_lists = function(l1,l2){
    var result = [].concat(l1);
    for(var i = 0; i < l2.length; i++){
	var found = false;
	for(var j = 0; j < l1.length; j++){
	    if(l2[i] == l1[j]){
		found = true;
		break;
	    }
	}
	if(!found) result.push(l2[i]);
    }
    return result;
}

var search = function(q, nodes){
    var result = [];
    if(!q || q.length == 0){
	for(var n in nodes){
	    result.push(n);
	}
    }
    else if(q[0] == "and"){
	result = search(q[1][0],nodes);
	for(var i = 1; i < q[1].length; i++){
	    var res = search(q[1][i],nodes);
	    result = intersect_lists(result,res);
	}
    }
    else if(q[0] == "or"){
	for(var i = 0; i < q[1].length; i++){
	    var res = search(q[1][i],nodes);
	    result = union_lists(result,res);
	}
    }
    else if(q[0] == "edge"){
	var edge = q[1];
	var res = search(edge.query,nodes);
	for(var i = 0; i < res.length; i++){
	    var node_id = res[i];
	    for(var n in nodes){
		if(edge.name == "*"){
		    for(var e in nodes[n].edges[edge.dir]){
			if(nodes[n].edges[edge.dir][e].indexOf(node_id) >= 0){
			    result = union_lists(result,[n]);
			    break;
			}
		    }
		}
		else if(edge.name in nodes[n].edges[edge.dir] && nodes[n].edges[edge.dir][edge.name].indexOf(node_id) >= 0){
		    result = union_lists(result,[n]);
		}
	    }
	}
    }
    else if(q[0] == "name"){
	var name = q[1].toLowerCase();
	for(var n in nodes){
	    if(name == "*" || nodes[n].name.toLowerCase().indexOf(name) >= 0){
		result.push(n);
	    }
	}
    }
    return result;
}
