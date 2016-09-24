#include "dataflow.h"
#include "transport.h"
#include "serialiser.h"
#include "error.h"

m_dataflow_t *m_dataflow_new(m_interface_t *interface,
			     m_transport_t *transport,
			     m_serialiser_t *serialiser,
			     m_args_t *dispatch(m_header_t *, m_args_t *),
			     void error(m_error_t *)){
  m_dataflow_t *d = malloc(sizeof(m_dataflow_t));
  d->interface = interface;
  d->transport = transport;
  d->serialiser = serialiser;
  d->dispath = dispatch;
  d->error = error;
}

void m_dataflow_send(m_dataflow_t *d, m_header_t *header, m_args_t *args){
  char *data = m_serialiser_serialise(d->serialiser, header, args);
  m_transport_tx(d->transport, data);
  free(data);
}
    
void m_dataflow_recv(m_dataflow_t *d){
  char *data = m_transport_rx(d->transport);
  m_dict_t *m = m_serialiser_deserialise(d->serialiser,data);
  if(!m_interface_validate(d->interface, m_dict_get(m,"header"))){
    d->error(UNKNOWN_COMMAND);
    return;
  }
  d->dispatch_cb(m_dict_get(m,"header"),m_dict_get(m,"args"));
  m_dict_free(m);
}
