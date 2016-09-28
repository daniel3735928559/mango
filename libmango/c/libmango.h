#ifndef LIBMANGO_H_
#define LIBMANGO_H_

#include "transport.h"
#include "serialiser.h"
#include "dataflow.h"
#include "interface.h"
#include "error.h"
#include "cJSON/cJSON.h"
#include "zmq.h"

#define LIBMANGO_VERSION "0.1"
#define LIBMANGO_REPLY "reply"
#define LIBMANGO_HELLO "hello"
#define LIBMANGO_STDIO "stdio"

typedef struct m_node m_node_t;

m_node_t *m_node_new(char debug);
void m_node_dispatch(m_node_t *node, cJSON *header, cJSON *args);
void m_node_handle_error(m_node_t *node, char *src, char *err);
cJSON *m_node_reg(m_node_t *node, cJSON *header, cJSON *args);
cJSON *m_node_reply(m_node_t *node, cJSON *header, cJSON *args);
cJSON *m_node_heartbeat(m_node_t *node, cJSON *header, cJSON *args);
cJSON *m_node_make_header(m_node_t *node, char *command, char *callback, int mid, char *src_port);
int m_node_get_mid(m_node_t *node);
void m_node_ready(m_node_t *node, cJSON *header, cJSON *args);
int m_node_send(m_node_t *node, char *command, cJSON *msg, char *callback, int mid, char *port);
void m_node_run(m_node_t *node);

#endif
