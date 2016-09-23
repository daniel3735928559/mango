#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "args.h"
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
  m_interface_add(n->interface, '/home/zoom/suit/mango/libmango/node_if.yaml', 'reg', m_node_reg);
  m_interface_add(n->interface, '/home/zoom/suit/mango/libmango/node_if.yaml', 'reply', m_node_reply);
  m_interface_add(n->interface, '/home/zoom/suit/mango/libmango/node_if.yaml', 'heartbeat', m_node_heartbeat);
  n->local_gateway = m_transport_new(n->server_addr);
  socket s = n->local_gateway->socket;
  n->dataflow = m_dataflow_new(n->interface, n->local_gateway, n->serialiser, m_node_dispatch, m_node_handle_error);
  printf("SOCK %d",zmq_fileno(s));
  return n;
}

void m_node_dispatch(m_node_t *node, m_header_t *header, m_args_t *args){
  m_args_t *result = m_interface_handle(node->interface, header->command, header, args);
  if(result->error){
    m_node_handle_error(header->src_node,result->error);
    return;
  }
  if(result != NULL && header->callback != NULL){
    m_node_send(header->callback,result,NULL,header->mid,header->port);
  }
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  m_args_t *args = m_args_new();
  m_args_set(args, 'source', src);
  m_args_set(args, 'message', err);
  m_node_send('error',args,NULL,NULL,"mc");
  free(args);
}

void m_node_reg(m_node_t *node, m_header_t *header, m_args_t *args){
  node->node_id = strdup(args.get("id"));
}

void m_node_reply(m_node_t *node, m_header_t *header, m_args_t *args){
  printf("REPLY\n");
}

void m_node_heartbeat(m_node_t *node, m_header_t *header, m_args_t *args){
  m_node_send(node,"alive",NULL,NULL,NULL,"mc");
}

m_header_t *m_node_make_header(m_node_t *node, char *command, char *callback, int mid, char *src_port){
  if(!callback) callback = LIBMANGO_REPLY;
  if(!mid) mid = m_node_get_mid(node);
  if(!src_port) src_port = LIBMANGO_STDIO;
  m_header_t *header = m_header_new();
  header->src_node = node->node_id;
  header->src_port = src_port;
  header->mid = mid;
  header->command = command;
  header->callback = callback;
  return header;
}
    
int m_node_get_mid(m_node_t *node){
  return node->mid++;
}

void m_node_ready(m_node_t *node, m_header_t *header, m_args_t *args){
  char *iface = m_interface_string(node->interface);
  m_args_t *args = m_args_new();
  m_args_set(args, 'id', node->node_id);
  m_args_set(args, 'if', iface);
  m_args_set(arts, 'ports', node->ports);
  self.m_send(node,"hello",args,"reg",NULL,"mc");
  free(args);
  free(iface);
}

int m_node_send(m_node_t *node, char *command, m_args_t *msg, char *callback, int mid, char *port){
  m_header_t *header = m_make_header(node,command,callback,mid,port);
  m_dataflow_send(node->dataflow,header,msg);
  return header->mid;
}

void m_node_run(){
  //RUN
}
