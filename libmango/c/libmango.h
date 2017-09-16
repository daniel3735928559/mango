#ifndef LIBMANGO_H_
#define LIBMANGO_H_
#include "cJSON/cJSON.h"

typedef struct m_result {
  char *name;
  cJSON *data;
} m_result_t;

#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "error.h"
#include "zmq.h"

#define LIBMANGO_VERSION "0.1"
#define LIBMANGO_REPLY "reply"
#define LIBMANGO_HELLO "hello"
#define LIBMANGO_STDIO "stdio"

typedef struct m_node m_node_t;

m_node_t *m_node_new(char debug);
void m_node_dispatch(m_node_t *node, cJSON *header, cJSON *args);
void m_node_handle_error(m_node_t *node, char *src, char *err);
cJSON *m_node_heartbeat(m_node_t *node, cJSON *header, cJSON *args, m_result_t *result);
cJSON *m_node_make_header(m_node_t *node, char *command, char *mid, char *type);
void m_node_send(m_node_t *node, char *command, cJSON *msg, char *mid, char *type);
void m_node_add_interface(m_node_t *node, char *filename);
int m_node_handle(m_node_t *node, char *fn_name, cJSON *(*handler)(m_node_t *, cJSON *, cJSON *, m_result_t *result));
void m_node_run(m_node_t *node);

#endif
