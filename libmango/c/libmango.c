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
  n->serialiser = m_serialiser_new(n->version);
  n->interface = m_interface_new();
  n->ports = NULL;
  n->num_ports = 0;
  n->server_addr = getenv("MC_ADDR");
  
  n->zmq_context = zmq_ctx_new();
  m_interface_load(n->interface, "/home/zoom/suit/mango/libmango/node_if.yaml");
  int x;
  x = m_interface_handle(n->interface, "reg", m_node_reg);
  x = m_interface_handle(n->interface, "reply", m_node_reply);
  x = m_interface_handle(n->interface, "heartbeat", m_node_heartbeat);
  n->local_gateway = m_transport_new(n->server_addr, n->zmq_context);
  n->dataflow = m_dataflow_new(n, n->local_gateway, n->serialiser, n->interface, m_node_dispatch, m_node_handle_error);
  return n;
}

void m_node_add_interface(m_node_t *node, char *filename){
  m_interface_load(node->interface, filename);
}

int m_node_handle(m_node_t *node, char *fn_name, cJSON *(*handler)(m_node_t *, cJSON *, cJSON *)){
  return m_interface_handle(node->interface, fn_name, handler);
}

void m_node_dispatch(m_node_t *node, cJSON *header, cJSON *args){
  cJSON *result = m_interface_handler(node->interface, cJSON_GetObjectItem(header,"command")->valuestring)(node, header, args);
  if(cJSON_HasObjectItem(result,"error")){
    m_node_handle_error(node,
			cJSON_GetObjectItem(header,"src_node")->valuestring,
			cJSON_GetObjectItem(result,"error")->valuestring);
    return;
  }
  if(result != NULL){
    m_node_send(node,
		"reply",
		result,
		cJSON_GetObjectItem(header,"port")->valuestring);
  }
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  cJSON *args = cJSON_CreateObject();
  cJSON_AddStringToObject(args, "source", src);
  cJSON_AddStringToObject(args, "message", err);
  m_node_send(node, "error", args, "mc");
  free(args);
}

cJSON *m_node_reg(m_node_t *node, cJSON *header, cJSON *args){
  if(strcmp(cJSON_GetObjectItem(header,"src_node")->valuestring,"mc")){
    return NULL;
  }
  node->node_id = strdup(cJSON_GetObjectItem(args,"id")->valuestring);
  return NULL;
}

cJSON *m_node_reply(m_node_t *node, cJSON *header, cJSON *args){
  return NULL;
}

cJSON *m_node_heartbeat(m_node_t *node, cJSON *header, cJSON *args){
  m_node_send(node,"alive",cJSON_CreateObject(),"mc");
  return NULL;
}

cJSON *m_node_make_header(m_node_t *node, char *command, char *src_port){
  if(!src_port) src_port = LIBMANGO_STDIO;
  cJSON *header = cJSON_CreateObject();
  cJSON_AddStringToObject(header,"src_port", src_port);
  cJSON_AddStringToObject(header,"command", command);
  return header;
}
    
void m_node_ready(m_node_t *node){
  char *iface = m_interface_string(node->interface);
  cJSON *hello_args = cJSON_CreateObject();
  cJSON *ports = cJSON_CreateStringArray(node->ports, node->num_ports);
  cJSON_AddStringToObject(hello_args, "id", node->node_id);
  cJSON_AddItemToObject(hello_args, "if", cJSON_Duplicate(m_interface_spec(node->interface),1));
  cJSON_AddItemToObject(hello_args, "ports", ports);
  m_node_send(node,"hello",hello_args,"mc");
  free(iface);
}

void m_node_send(m_node_t *node, char *command, cJSON *msg, char *port){
  cJSON *header = m_node_make_header(node,command,port);
  m_dataflow_send(node->dataflow,header,msg);
}

void m_node_run(m_node_t *node){
  m_node_ready(node);
  while(1){
    zmq_pollitem_t items [] = {{node->local_gateway->socket, 0, ZMQ_POLLIN, 0}};
    zmq_poll(items, 1, 10);
    if(items[0].revents & ZMQ_POLLIN){
      m_dataflow_recv(node->dataflow);
    }
  }
}
