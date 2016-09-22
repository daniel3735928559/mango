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
	/* console.log("DISPATCH",header,args,self.iface.iface,self.iface.iface[header['command']]); */
	/* try{ */
        /*     result = self.iface.iface[header['command']]['handler'](header,args); */
        /*     if(result && 'callback' in header){ */
	/* 	self.m_send(header['callback'],result,null,header['mid'],header['port']) */
	/*     } */
	/* } catch(e) { */
	/*     console.log(e); */
        /*     self.handle_error(header['src_node'],e+"") */
	/* } */
}

void m_node_handle_error(m_node_t *node, char *src, char *err){
  /* console.log('OOPS',src,err); */
  /* self.m_send('error',{'source':src,'message':err},null,null,"mc"); */
}

void m_node_reg(m_node_t *node, m_header_t *header, m_args_t *args){
  //self.node_id = args["id"];
}

void m_node_reply(m_node_t *node, m_header_t *header, m_args_t *args){
  //console.log("REPLY",header,args);
}

void m_node_heartbeat(m_node_t *node, m_header_t *header, m_args_t *args){
  //self.m_send("alive",{},null,null,"mc");
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
  
	/* var ifce = {}; */
	/* for(var i in this.iface.iface){ */
	/*     ifce[i] = JSON.parse(JSON.stringify(this.iface.iface[i])); */
	/*     delete ifce[i]['handler']; */
	/* } */
	/* console.log("IF",ifce) */
	/* self.m_send('hello',{'id':self.node_id,'if':ifce,'ports':self.ports},"reg",null,"mc") */
}

int m_node_send(m_node_t *node, char *command, m_args_t *msg, char *callback, int mid, char *port){
  /* console.log('sending',command,msg,mid,port) */
  /*   header = self.make_header(command,callback,mid,port) */
  /*   self.dataflow.send(header,msg) */
  /*   return header['mid'] */
}

void m_node_run(){
  //RUN
}
