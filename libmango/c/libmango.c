#include "zmq.h"
#include "libmango.h"
#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "cJSON/cJSON.h"
#include "error.h"

struct m_node {
  char *version;
  char *node_id;
  int mid;
  const char **ports;
  int num_ports;
  char debug;
  char *server_addr;
  m_interface_t *interface;
  m_serialiser_t *serialiser;
  m_transport_t *local_gateway;
  m_dataflow_t *dataflow;
  void *zmq_context;
};

m_node_t *m_node_new(char debug){
  m_node_t *n = malloc(sizeof(m_node_t));
  n->version = LIBMANGO_VERSION;
  n->debug = debug;
  n->node_id = getenv("MANGO_ID");
  n->mid = 0;
  n->serialiser = m_serialiser_new(n->version);
  n->interface = m_interface_new();
  n->ports = NULL;
  n->num_ports = 0;
  n->server_addr = getenv("MC_ADDR");
  
  n->zmq_context = zmq_ctx_new();
  m_interface_load(n->interface, "/home/zoom/suit/mango/libmango/node_if.yaml");
  m_interface_handle(n->interface, "reg", m_node_reg);
  m_interface_handle(n->interface, "reply", m_node_reply);
  m_interface_handle(n->interface, "heartbeat", m_node_heartbeat);
  n->local_gateway = m_transport_new(n->server_addr, n->zmq_context);
  void *s = n->local_gateway->socket;  
  n->dataflow = m_dataflow_new(n->interface, n->local_gateway, n->serialiser, m_node_dispatch, m_node_handle_error);
  // printf("SOCK %d",zmq_fileno(s));
  return n;
}

void m_node_dispatch(m_node_t *node, cJSON *header, cJSON *args){
  cJSON *result = m_interface_handler(node->interface, cJSON_GetObjectItem(header,"command")->valuestring)(node, header, args);
  if(cJSON_HasObjectItem(result,"error")){
    m_node_handle_error(node,
			cJSON_GetObjectItem(header,"src_node")->valuestring,
			cJSON_GetObjectItem(result,"error")->valuestring);
    return;
  }
  if(result != NULL && cJSON_GetObjectItem(header,"callback") != NULL){
    m_node_send(node,
		cJSON_GetObjectItem(header,"callback")->valuestring,
		result,
		NULL,
		cJSON_GetObjectItem(header,"mid")->valueint,
		cJSON_GetObjectItem(header,"port")->valuestring);
  }
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  cJSON *args = cJSON_CreateObject();
  cJSON_AddStringToObject(args, "source", src);
  cJSON_AddStringToObject(args, "message", err);
  m_node_send(node, "error", args, NULL, 0, "mc");
  free(args);
}

cJSON *m_node_reg(m_node_t *node, cJSON *header, cJSON *args){
  node->node_id = strdup(cJSON_GetObjectItem(args,"id")->valuestring);
}

cJSON *m_node_reply(m_node_t *node, cJSON *header, cJSON *args){
  printf("REPLY\n");
}

cJSON *m_node_heartbeat(m_node_t *node, cJSON *header, cJSON *args){
  m_node_send(node,"alive",NULL,NULL,0,"mc");
}

cJSON *m_node_make_header(m_node_t *node, char *command, char *callback, int mid, char *src_port){
  if(!callback) callback = LIBMANGO_REPLY;
  if(!mid) mid = m_node_get_mid(node);
  if(!src_port) src_port = LIBMANGO_STDIO;
  cJSON *header = cJSON_CreateObject();
  cJSON_AddStringToObject(header,"src_node", node->node_id);
  cJSON_AddStringToObject(header,"src_port", src_port);
  cJSON_AddNumberToObject(header,"mid", mid);
  cJSON_AddStringToObject(header,"command", command);
  cJSON_AddStringToObject(header,"callback", callback);
  return header;
}
    
int m_node_get_mid(m_node_t *node){
  return node->mid++;
}

void m_node_ready(m_node_t *node, cJSON *header, cJSON *args){
  char *iface = m_interface_string(node->interface);
  cJSON *hello_args = cJSON_CreateObject();
  cJSON *ports = cJSON_CreateStringArray(node->ports, node->num_ports);
  cJSON_AddStringToObject(hello_args, "id", node->node_id);
  cJSON_AddStringToObject(hello_args, "if", iface);
  cJSON_AddItemToObject(hello_args, "ports", ports);
  m_node_send(node,"hello",args,"reg",0,"mc");
  free(args);
  free(iface);
}

int m_node_send(m_node_t *node, char *command, cJSON *msg, char *callback, int mid, char *port){
  cJSON *header = m_node_make_header(node,command,callback,mid,port);
  m_dataflow_send(node->dataflow,header,msg);
  return cJSON_GetObjectItem(header,"mid")->valueint;
}

void m_node_run(m_node_t *node){
  while(1){
    zmq_pollitem_t items [] = {
      {node->local_gateway->socket, 0, ZMQ_POLLIN, 0},
    };
    zmq_poll(items, 2, -1);
    if(items[0].revents & ZMQ_POLLIN){
      m_dataflow_recv(node->dataflow);
    }
  }
}
