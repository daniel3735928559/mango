ws = new ReconnectingWebSocket("ws://"  + location.host + '/zeromq')
ws2 = new ReconnectingWebSocket("ws://"  + location.host + '/test')

ws.onmessage = function(message) {
    payload = JSON.parse(message.data);
    document.getElementById('latest_data').innerHTML = '<h2> Data: ' + message.data + '</h2>';
};

function stuff(){
  ws2.send("asdasda");
}
