#include <stdlib.h>
#include <string.h>
#include "zmq.h"
#include "transport.h"
#include "string.h"

m_transport_t *m_transport_new(char *addr, void *context){
  m_transport_t *t = malloc(sizeof(m_transport_t));
  t->target = strdup(addr);
  t->socket = zmq_socket(context, ZMQ_DEALER);
  zmq_connect(t->socket, t->target);
  return t;
}

void m_transport_tx(m_transport_t *t, char *data){
  zmq_send(t->socket, data, strlen(data), 0);
}

char *m_transport_rx(m_transport_t *t){
  zmq_msg_t msg;
  zmq_msg_init(&msg);
  int rc = zmq_msg_recv(&msg, t->socket, 0);
  if(rc == -1){
    perror("Problem");
    return NULL;
  }
  size_t sz = zmq_msg_size(&msg);
  char *data = malloc(sz+1);
  data[sz] = 0;
  memcpy(data, zmq_msg_data(&msg), sz);
  zmq_msg_close(&msg);
  return data;
}

void m_transport_free(m_transport_t *t){
  free(t->target);
  zmq_close(t->socket);
  free(t);
}
