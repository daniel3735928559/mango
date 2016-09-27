#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "error.h"
#include "cJSON/cJSON.h"
#include "zmq.h"

char *LIBMANGO_VERSION = "0.1";
char *LIBMANGO_REPLY = "reply";
char *LIBMANGO_HELLO = "hello";
char *LIBMANGO_STDIO = "stdio";

typedef struct m_node {
  char *version;
  char *node_id;
  uint32 mid;
  char **ports;
  char debug;
  char *server_addr;
  m_interface_t *interface;
  m_serialiser_t *serialiser;
  m_transport_t *local_gateway;
  m_dataflow_t *dataflow;
  void *zmq_context;
} m_node_t;

void m_node_new(char debug);
int m_node_send(m_node_t *node, char *command, cJSON *msg, char *callback, int mid, char *port);
void m_node_ready(m_node_t *node, cJSON *header, cJSON *args);
void m_node_run(m_node_t *node);
