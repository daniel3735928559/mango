#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "args.h"
#include "error.h"

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
} m_node_t;

void m_node_new(char debug);
void m_node_ready(m_node_t *node, m_header_t *header, m_args_t *args);
int m_node_send(m_node_t *node, char *command, m_args_t *msg, char *callback, int mid, char *port);
void m_node_run();
