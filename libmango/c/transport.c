#include "transport.h"

void m_transport_tx(m_transport_t *t, char *data){
  zmq_send(t->socket, data);
}

char *m_transport_rx(){
  return zmq_recv(t->socket);
}
