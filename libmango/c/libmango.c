#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "cJSON/cJSON.h"
#include "error.h"

void m_node_new(char debug){
  m_node_t *n = malloc(sizeof(m_node_t));
  n->version = LIBMANGO_VERSION;
  n->debug = debug;
  n->node_id = getenv('MANGO_ID');
  n->mid = 0;
  n->serialiser = m_serialiser_new(n->version);
  n->interface = m_interface_new();
  n->ports = NULL;
  n->server_addr = getenv('MC_ADDR');
  
  n->zmq_context = zmq_ctx_new();
  m_interface_load(n->interface, '/home/zoom/suit/mango/libmango/node_if.yaml');
  m_interface_handle(n->interface, 'reg', m_node_reg);
  m_interface_handle(n->interface, 'reply', m_node_reply);
  m_interface_handle(n->interface, 'heartbeat', m_node_heartbeat);
  n->local_gateway = m_transport_new(n->server_addr);
  socket s = n->local_gateway->socket;  
  n->dataflow = m_dataflow_new(n->interface, n->local_gateway, n->serialiser, m_node_dispatch, m_node_handle_error);
  printf("SOCK %d",zmq_fileno(s));
  return n;
}

void m_node_dispatch(m_node_t *node, cJSON *header, cJSON *args){
  cJSON *result = *(m_interface_handler(node->interface, cJSON_getObjectItem(header,"command")->valuestring))(node, header, args);
  if(result->error){
    m_node_handle_error(cJSON_getObjectItem(header,"src_node")->valuestring,
			result->error);
    return;
  }
  if(result != NULL && cJSON_getObjectItem(header,"callback") != NULL){
    m_node_send(cJSON_getObjectItem(header,"callback")->valuestring,
		result,
		NULL,
		cJSON_getObjectItem(header,"mid")->valueint,
		cJSON_getObjectItem(header,"port")->valuestring);
  }
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  cJSON *args = cJSON_CreateObject();
  cJSON_AddStringToObject(args, 'source', src);
  cJSON_AddStringToObject(args, 'message', err);
  m_node_send('error',args,NULL,NULL,"mc");
  free(args);
}

cJSON *m_node_reg(m_node_t *node, cJSON *header, cJSON *args){
  node->node_id = strdup(args.get("id"));
}

cJSON *m_node_reply(m_node_t *node, cJSON *header, cJSON *args){
  printf("REPLY\n");
}

cJSON *m_node_heartbeat(m_node_t *node, cJSON *header, cJSON *args){
  m_node_send(node,"alive",NULL,NULL,NULL,"mc");
}

cJSON *m_node_make_header(m_node_t *node, char *command, char *callback, int mid, char *src_port){
  if(!callback) callback = LIBMANGO_REPLY;
  if(!mid) mid = m_node_get_mid(node);
  if(!src_port) src_port = LIBMANGO_STDIO;
  cJSON *header = cJSON_CreateObject();
  cJSON_AddStringToObject(header,"src_node", node->node_id);
  cJSON_AddStringToObject(header,"src_port", src_port);
  cJSON_AddStringToObject(header,"mid", mid);
  cJSON_AddStringToObject(header,"command", command);
  cJSON_AddStringToObject(header,"callback", callback);
  return header;
}
    
int m_node_get_mid(m_node_t *node){
  return node->mid++;
}

void m_node_ready(m_node_t *node, cJSON *header, cJSON *args){
  char *iface = m_interface_string(node->interface);
  cJSON *args = cJSON_CreateObject();
  cJSON *ports = cJSON_CreateObject();
  cJSON_AddStringToObject(args, 'id', node->node_id);
  cJSON_AddStringToObject(args, 'if', iface);
  cJSON_AddItemToObject(args, 'ports', node->ports);
  self.m_send(node,"hello",args,"reg",NULL,"mc");
  free(args);
  free(iface);
}

int m_node_send(m_node_t *node, char *command, cJSON *msg, char *callback, int mid, char *port){
  cJSON *header = m_make_header(node,command,callback,mid,port);
  m_dataflow_send(node->dataflow,header,msg);
  return cJSON_getObjectItem(header,"mid")->valueint;
}

void m_node_run(m_node_t *node){
  while(1){
    zmq_pollitem_t items [] = {
      {node->local_gateway->socket, 0, ZMQ_POLLIN, 0},
    };
    zmq_poll (items, 2, -1);
    if(items[0].revents & ZMQ_POLLIN){
      m_dataflow_recv(node->dataflow);
    }
  }
}
