#include "dataflow.h"
#include "transport.h"
#include "serialiser.h"
#include "error.h"
#include "cJSON/cJSON.h"

struct m_dataflow {
  m_interface_t *interface;
  m_transport_t *transport;
  m_serialiser_t *serialiser;
  void (*dispatch)(cJSON *, cJSON *);
  void (*error)(m_error_t *);
};
  
m_dataflow_t *m_dataflow_new(m_interface_t *interface,
			     m_transport_t *transport,
			     m_serialiser_t *serialiser,
			     void (*dispatch)(cJSON *, cJSON *),
			     void (*error)(char *, char *)){
  m_dataflow_t *d = malloc(sizeof(m_dataflow_t));
  d->interface = interface;
  d->transport = transport;
  d->serialiser = serialiser;
  d->dispatch = dispatch;
  d->error = error;
}

void m_dataflow_send(m_dataflow_t *d, cJSON *header, cJSON *args){
  char *data = m_serialiser_serialise(d->serialiser, header, args);
  m_transport_tx(d->transport, data);
  free(data);
}
    
void m_dataflow_recv(m_dataflow_t *d){
  char *data = m_transport_rx(d->transport);
  cJSON *m = m_serialiser_deserialise(d->serialiser,data);
  if(!m_interface_validate(d->interface, cJSON_GetObjectItem(cJSON_GetObjectItem(m,"header"),"command")->valuestring)){
    d->error(UNKNOWN_COMMAND);
    return;
  }
  d->dispatch(cJSON_GetObjectItem(m,"header"), cJSON_GetObjectItem(m,"args"));
  free(data);
  cJSON_Delete(m);
}
