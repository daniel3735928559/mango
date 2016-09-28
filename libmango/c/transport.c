#include <stdlib.h>
#include <string.h>
#include "zmq.h"
#include "transport.h"
#include "string.h"

m_transport_t *m_transport_new(char *addr, void *context){
  m_transport_t *t = malloc(sizeof(m_transport_t));
  t->target = strdup(addr);
  t->socket = zmq_socket(context, ZMQ_DEALER);
  printf("%s\n",t->target);
  zmq_connect(t->socket, t->target);
  return t;
}

void m_transport_tx(m_transport_t *t, char *data){
  printf("TX %d %s\n",strlen(data),data);
  zmq_send(t->socket, data, strlen(data), 0);
  printf("ZMQ SENT\n");
}

char *m_transport_rx(m_transport_t *t){
  printf("RX\n");
  int cur_max = 256;
  char *msg = malloc(cur_max);
  memset(msg,0,cur_max);
  int size = zmq_recv(t->socket, msg, cur_max-1, 0);
  printf("RXED0 %s\n",msg);
  int msg_size = size;
  if(size == -1){
    return NULL;
  }
  else if(size < cur_max-1){
    printf("RXED FIN %s\n",msg);
    return msg;
  }
  while(1){
    int size = zmq_recv(t->socket, msg, cur_max-1, 0);
    printf("RXED1 %s\n",msg);
    if(size == -1){
      return msg;
    }
    else{
      char *msg2 = malloc(cur_max);
      memset(msg2,0,cur_max);
      char *new_msg = malloc(2*cur_max);
      memset(new_msg,0,2*cur_max);
      memcpy(new_msg, msg, msg_size);
      memcpy(new_msg+msg_size, msg2, size);
      msg_size += size;
      cur_max *= 2;
      free(msg);
      free(msg2);
      msg = new_msg;
    }
  }
  printf("RXED %s\n",msg);
  return msg;
}

void m_transport_free(m_transport_t *t){
  free(t->target);
  zmq_close(t->socket);
  free(t);
}
