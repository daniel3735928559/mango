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
  cJSON **m = m_serialiser_deserialise(d->serialiser,data);
  cJSON *header = m[0];
  cJSON *args = m[1];

  if(!cJSON_HasObjectItem(header,"command")){
    m_debug_print(d->node, "ERROR", "No command");
    d->error(d->node, "no command", "Invalid header");
  } else {
    char *cmd = cJSON_GetObjectItem(header,"command")->valuestring;
    if(!m_interface_validate(d->interface, cmd)){
      m_debug_print(d->node, "ERROR", cmd);
      d->error(d->node, cmd, "Unknown message");
    } else {
      m_debug_print(d->node, "DISPATCH", cmd);
      d->dispatch(d->node, header, args);
    }
  }
  free(data);
  cJSON_Delete(header);
  cJSON_Delete(args);
}
