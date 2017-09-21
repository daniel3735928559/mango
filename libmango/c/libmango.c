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
  char *route;
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
  n->route = getenv("MANGO_ROUTE");
  n->serialiser = m_serialiser_new(n->version);
  n->interface = m_interface_new();
  n->ports = NULL;
  n->num_ports = 0;
  n->server_addr = getenv("MC_ADDR");
  char *node_if_path = malloc(strlen(getenv("LD_LIBRARY_PATH"))+strlen("/../node.yaml")+1);
  sprintf(node_if_path, "%s/../node.yaml",getenv("LD_LIBRARY_PATH"));
  n->zmq_context = zmq_ctx_new();
  m_interface_load(n->interface, node_if_path);
  free(node_if_path);
  n->local_gateway = m_transport_new(n->server_addr, n->zmq_context, n->route);
  n->dataflow = m_dataflow_new(n, n->local_gateway, n->serialiser, n->interface, m_node_dispatch, m_node_handle_error);
  m_node_handle(n, "heartbeat", m_node_heartbeat);
  m_debug_print(n,"HELLO",NULL);
  return n;
}

void m_node_add_interface(m_node_t *node, char *filename){
  m_interface_load(node->interface, filename);
}

int m_node_handle(m_node_t *node, char *fn_name, cJSON *(*handler)(m_node_t *, cJSON *, cJSON *, m_result_t *result)){
  return m_interface_handle(node->interface, fn_name, handler);
}

void m_node_dispatch(m_node_t *node, cJSON *header, cJSON *args){
  m_result_t *result = (m_result_t *)malloc(sizeof(m_result_t));
  m_debug_print(node, "RX", cJSON_GetObjectItem(header,"name")->valuestring);
  result->name = NULL;
  result->data = cJSON_CreateObject();
  m_interface_handler(node->interface, cJSON_GetObjectItem(header,"name")->valuestring)(node, header, args, result);
  if(cJSON_HasObjectItem(result->data,"error")){
    m_node_handle_error(node,
			cJSON_GetObjectItem(header,"src_node")->valuestring,
			cJSON_GetObjectItem(result->data,"error")->valuestring);
    return;
  }
  if(result->name != NULL){
    m_node_send(node, result->name, result->data, NULL, NULL);
  }
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  cJSON *args = cJSON_CreateObject();
  cJSON_AddStringToObject(args, "source", src);
  cJSON_AddStringToObject(args, "message", err);
  m_node_send(node, "error", args, NULL, "system");
  free(args);
}

cJSON *m_node_heartbeat(m_node_t *node, cJSON *header, cJSON *args, m_result_t *result){
  m_debug_print(node, "HB", NULL);
  m_node_send(node,"alive",cJSON_CreateObject(),cJSON_GetObjectItem(header,"mid")->valuestring,"system");
  return NULL;
}

cJSON *m_node_make_header(m_node_t *node, char *name, char *mid, char *type){
  cJSON *header = cJSON_CreateObject();
  cJSON_AddStringToObject(header,"name", name);
  if(mid) cJSON_AddStringToObject(header,"mid", mid);
  if(type) cJSON_AddStringToObject(header,"type", type);
  return header;
}

void m_node_send(m_node_t *node, char *name, cJSON *msg, char *mid, char *type){
  m_debug_print(node, "SENDING", name);
  cJSON *header = m_node_make_header(node,name,mid,type);
  m_dataflow_send(node->dataflow,header,msg);
}

void m_debug_print(m_node_t *node, char *tag, char *msg){
  if(node->debug) printf("[%s DEBUG] %s %s\n", node->node_id, tag, msg);
}

void m_node_run(m_node_t *node){
  while(1){
    zmq_pollitem_t items [] = {{node->local_gateway->socket, 0, ZMQ_POLLIN, 0}};
    zmq_poll(items, 1, 10);
    if(items[0].revents & ZMQ_POLLIN){
      m_dataflow_recv(node->dataflow);
    }
  }
}
