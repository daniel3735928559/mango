var outstanding = {};
var commands = {};
var mid = 0;

var cache_killer = 100
var xhrpoller;

function xhr_recv(){
    if(xhrpoller.readyState == 4){
	console.log('GOt one!');
	console.log(xhrpoller.responseText);
	if(xhrpoller.responseText.length > 0){
	    m_recv(xhrpoller.responseText);
	}
	poll();
    }
}

function poll(){
    console.log('polling...');
    xhrpoller = new XMLHttpRequest();
    xhrpoller.onreadystatechange = function(){ console.log("booyah",xhrpoller.readyState,xhrpoller.responseText); xhr_recv() };
    xhrpoller.open("POST", "/poller", true);
    cache_killer++;
    xhrpoller.timeout = 30000;
    xhrpoller.ontimeout = poll;
    xhrpoller.setRequestHeader("Content-type", "text/plain");
    xhrpoller.send('');
}

window.onload = poll;

function m_send(dict,cb){
    console.log(JSON.stringify(dict));
    dict['mid'] = mid+'';
    outstanding[mid] = cb;
    mid++;
    var msg = "MANGO0.1 json\n"+JSON.stringify({"header":{"source/stdio":"mu","mid":mid+""},"args":dict});
    console.log(dict);
    var xhr = new XMLHttpRequest();
    var tries = 0;
    xhr.onreadystatechange = function(){
	if(xhr.readyState == 4){
	    if(xhr.status == 200){
		console.log(xhr.responseText);
	    }
	    else{
		console.log(xhr.status);
		tries++;
		if(tries < 3){
		    console.log("An error occurred--resending... (retry "+tries+")");
		    xhr.open("POST", "/", true);
		    xhr.setRequestHeader("Content-type", "text/plain");
		    xhr.send(msg);
		}
	    }
	}
    }
    xhr.open("POST", "/", true);
    xhr.setRequestHeader("Content-type", "text/plain");
    xhr.send(msg);
}

function m_recv(dict){
    try{
	console.log("AD")
	var nl1 = dict.indexOf("\n")
	console.log("HD",dict.substring(nl1))
	dict = JSON.parse(dict.substring(nl1));
    }
    catch(ex){
	console.log('Got non-JSON: ',dict);
	return
    }
    console.log("recv...",dict);
    var args = dict["args"]
    if(args["command"] == "reply"){
	console.log("OUT",outstanding);
	if(args['reply'] in outstanding)
	    outstanding[args['reply']](args);
	else
	    console.log('Fake reply',JSON.stringify(dict));
    }
    else if('command' in args){
	if(args['command'] in commands)
	    commands[args['command']](args);
	else
	    console.log('Bad command',JSON.stringify(dict));
    }
}
