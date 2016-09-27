#ifndef M_DATAFLOW_H_
#define M_DATAFLOW_H_

#include "libmango.h"
#include "transport.h"
#include "serialiser.h"
#include "interface.h"
#include "error.h"
#include "libmango.h"
#include "cJSON/cJSON.h"

typedef struct m_dataflow m_dataflow_t;
struct m_node;

m_dataflow_t *m_dataflow_new(struct m_node *node,
			     m_transport_t *transport,
			     m_serialiser_t *serialiser,
			     m_interface_t *interface,
			     void (*dispatch)(struct m_node *, cJSON *, cJSON *),
			     void (*error)(struct m_node *, char *, char *));
void m_dataflow_send(m_dataflow_t *d, cJSON *header, cJSON *args);
void m_dataflow_recv(m_dataflow_t *d);

#endif
