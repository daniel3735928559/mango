#include <time.h>
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

char *gen_id(int len) {
  char *ans = (char *)malloc(len+1);
  char *s = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  for(int i = 0; i < len; i++) {
    ans[i] = s[rand()%66];
  }
  ans[len] = '\0';
  return ans;
}

m_node_t *m_node_new(char debug){
  srand(time(NULL));
  m_node_t *n = malloc(sizeof(m_node_t));
  n->version = LIBMANGO_VERSION;
  n->debug = debug;
  n->node_id = getenv("MANGO_COOKIE");
  n->serialiser = m_serialiser_new(n->version);
  n->interface = m_interface_new();
  n->ports = NULL;
  n->num_ports = 0;
  n->server_addr = getenv("MANGO_SERVER");
  n->zmq_context = zmq_ctx_new();
  n->local_gateway = m_transport_new(n->server_addr, n->zmq_context, gen_id(16));
  n->dataflow = m_dataflow_new(n, n->local_gateway, n->serialiser, n->interface, m_node_dispatch, m_node_handle_error);
  m_node_handle(n, "heartbeat", m_node_heartbeat);
  m_debug_print(n,"HELLO",NULL);
  return n;
}

int m_node_handle(m_node_t *node, char *fn_name, cJSON *(*handler)(m_node_t *, cJSON *, m_result_t *result)){
  return m_interface_handle(node->interface, fn_name, handler);
}

void m_node_dispatch(m_node_t *node, cJSON *header, cJSON *args){
  m_result_t *result = (m_result_t *)malloc(sizeof(m_result_t));
  m_debug_print(node, "RX", cJSON_GetObjectItem(header,"command")->valuestring);
  result->command = NULL;
  result->data = cJSON_CreateObject();
  m_interface_handler(node->interface, cJSON_GetObjectItem(header,"command")->valuestring)(node, args, result);
  if(cJSON_HasObjectItem(result->data,"error")){
    m_node_handle_error(node,
			cJSON_GetObjectItem(header,"source")->valuestring,
			cJSON_GetObjectItem(result->data,"error")->valuestring);
    return;
  }
  if(result->command != NULL){
    m_node_send(node, result->command, result->data, cJSON_GetObjectItem(header,"mid")->valuestring);
  }
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  cJSON *args = cJSON_CreateObject();
  cJSON_AddStringToObject(args, "source", src);
  cJSON_AddStringToObject(args, "message", err);
  m_node_send(node, "error", args, NULL);
  free(args);
}

cJSON *m_node_heartbeat(m_node_t *node, cJSON *args, m_result_t *result){
  m_debug_print(node, "HB", NULL);
  m_node_send(node,"alive",cJSON_CreateObject(),NULL);
  return NULL;
}

cJSON *m_node_make_header(m_node_t *node, char *command, char *mid){
  cJSON *header = cJSON_CreateObject();
  cJSON_AddStringToObject(header,"format", "json");
  cJSON_AddStringToObject(header,"command", command);
  cJSON_AddStringToObject(header,"cookie", node->node_id);
  if(mid){
    cJSON_AddStringToObject(header,"mid", mid);
  } else {
    cJSON_AddStringToObject(header,"mid", gen_id(16));
  }
  return header;
}

void m_node_send(m_node_t *node, char *name, cJSON *msg, char *mid){
  m_debug_print(node, "SENDING", name);
  cJSON *header = m_node_make_header(node,name,mid);
  m_dataflow_send(node->dataflow,header,msg);
}

void m_debug_print(m_node_t *node, char *tag, char *msg){
  if(node->debug) printf("[%s DEBUG] %s %s\n", node->node_id, tag, msg);
}

void m_node_run(m_node_t *node){
  m_node_send(node,"alive",cJSON_CreateObject(),NULL);
  while(1){
    zmq_pollitem_t items [] = {{node->local_gateway->socket, 0, ZMQ_POLLIN, 0}};
    zmq_poll(items, 1, 10);
    if(items[0].revents & ZMQ_POLLIN){
      m_dataflow_recv(node->dataflow);
    }
  }
}
