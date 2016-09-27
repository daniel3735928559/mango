#include "transport.h"
#include "string.h"

m_transport_t *m_transport_new(char *addr, void *context){
  m_transport_t *t = malloc(sizeof(m_transport_t));
  t->target = strdup(addr);
  t->socket = zmq_socket(context, ZMQ_DEALER);
  zmq_connect(t->socket, t->target);
}

void m_transport_tx(m_transport_t *t, char *data){
  zmq_send(t->socket, data);
}

char *m_transport_rx(m_transport_t *t){
  int cur_max = 256;
  int size = 0;
  char *msg = malloc(cur_max);
  int size = zmq_recv(t->socket, msg, cur_max-1, 0);
  int msg_size = size;
  if(size == -1){
    return NULL;
  }
  while(1){
    char *msg2 = malloc(cur_max);
    int size = zmq_recv(t->socket, msg, cur_max-1, 0);
    if(size == -1){
      free(msg2);
      return msg;
    }
    else{
      char *new_msg = malloc(2*cur_max);
      memcopy(new_msg, msg, msg_size);
      memcopy(new_msg+msg_size, msg2, size);
      msg_size += size;
      cur_max *= 2;
      free(msg);
      free(msg2);
      msg = new_msg;
    }
  }
  return msg;
}

void m_transport_free(m_transport_t *t){
  free(t->addr);
  zmq_close(t->socket);
  free(t);
}
