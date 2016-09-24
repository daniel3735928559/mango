#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "dict.h"
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

void m_node_dispatch(m_node_t *node, m_dict_t *header, m_dict_t *args){
  m_dict_t *result = *(m_interface_handler(node->interface, m_dict_get(header,"command")))(node, header, args);
  if(result->error){
    m_node_handle_error(m_dict_get(header,"src_node"),result->error);
    return;
  }
  if(result != NULL && m_dict_get(header,"callback") != NULL){
    m_node_send(m_dict_get(header,"callback"),result,NULL,m_dict_get(header,"mid"),m_dict_get(header,"port"));
  }
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  m_dict_t *args = m_dict_new();
  m_dict_set(args, 'source', src);
  m_dict_set(args, 'message', err);
  m_node_send('error',args,NULL,NULL,"mc");
  free(args);
}

m_dict_t *m_node_reg(m_node_t *node, m_dict_t *header, m_dict_t *args){
  node->node_id = strdup(args.get("id"));
}

m_dict_t *m_node_reply(m_node_t *node, m_dict_t *header, m_dict_t *args){
  printf("REPLY\n");
}

m_dict_t *m_node_heartbeat(m_node_t *node, m_dict_t *header, m_dict_t *args){
  m_node_send(node,"alive",NULL,NULL,NULL,"mc");
}

m_dict_t *m_node_make_header(m_node_t *node, char *command, char *callback, int mid, char *src_port){
  if(!callback) callback = LIBMANGO_REPLY;
  if(!mid) mid = m_node_get_mid(node);
  if(!src_port) src_port = LIBMANGO_STDIO;
  m_dict_t *header = m_dict_new(0);
  m_dict_set(header,"src_node", node->node_id);
  m_dict_set(header,"src_port", src_port);
  m_dict_set(header,"mid", mid);
  m_dict_set(header,"command", command);
  m_dict_set(header,"callback", callback);
  return header;
}
    
int m_node_get_mid(m_node_t *node){
  return node->mid++;
}

void m_node_ready(m_node_t *node, m_dict_t *header, m_dict_t *args){
  char *iface = m_interface_string(node->interface);
  m_dict_t *args = m_dict_new();
  m_dict_set(args, 'id', node->node_id);
  m_dict_set(args, 'if', iface);
  m_dict_set(args, 'ports', node->ports);
  self.m_send(node,"hello",args,"reg",NULL,"mc");
  free(args);
  free(iface);
}

int m_node_send(m_node_t *node, char *command, m_dict_t *msg, char *callback, int mid, char *port){
  m_dict_t *header = m_make_header(node,command,callback,mid,port);
  m_dataflow_send(node->dataflow,header,msg);
  return m_dict_get(header,"mid");
}

void m_node_run(){
  //RUN
}
