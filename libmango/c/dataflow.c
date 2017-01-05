#include "libmango.h"
#include "dataflow.h"
#include "transport.h"
#include "serialiser.h"
#include "error.h"
#include "cJSON/cJSON.h"

struct m_dataflow {
  m_node_t *node;
  m_interface_t *interface;
  m_transport_t *transport;
  m_serialiser_t *serialiser;
  void (*dispatch)(m_node_t *, cJSON *, cJSON *);
  void (*error)(m_node_t *, char *, char *);
};
  
m_dataflow_t *m_dataflow_new(m_node_t *node,
			     m_transport_t *transport,
			     m_serialiser_t *serialiser,
			     m_interface_t *interface,
			     void (*dispatch)(m_node_t *, cJSON *, cJSON *),
			     void (*error)(m_node_t *, char *, char *)){
  m_dataflow_t *d = malloc(sizeof(m_dataflow_t));
  d->node = node;
  d->interface = interface;
  d->transport = transport;
  d->serialiser = serialiser;
  d->dispatch = dispatch;
  d->error = error;
  return d;
}

void m_dataflow_send(m_dataflow_t *d, cJSON *header, cJSON *args){
  char *data = m_serialiser_serialise(d->serialiser, header, args);
  m_transport_tx(d->transport, data);
  free(data);
}
    
void m_dataflow_recv(m_dataflow_t *d){
  char *data = m_transport_rx(d->transport);
  cJSON *m = m_serialiser_deserialise(d->serialiser,data);
  cJSON *header = cJSON_GetObjectItem(m,"header");
  cJSON *args = cJSON_GetObjectItem(m,"args");
  
  if(!m_interface_validate(d->interface, cJSON_GetObjectItem(header,"name")->valuestring)){
    d->error(d->node, cJSON_GetObjectItem(header,"src_port")->valuestring, "Unkown message");
    return;
  }
  d->dispatch(d->node, header, args);
  free(data);
  cJSON_Delete(m);
}
