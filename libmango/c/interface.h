#ifndef M_INTERFACE_H_
#define M_INTERFACE_H_

#include <stdlib.h>
#include "libmango.h"
#include "dict.h"
#include "cJSON/cJSON.h"

#define LIBMANGO_DEFAULT_INTERFACE_SIZE 64

enum storage_flags{VAR, VAL, SEQ};
typedef struct m_function m_function_t;
typedef struct m_interface m_interface_t;

struct m_node;

m_interface_t *m_interface_new();
int m_interface_load(m_interface_t *i, char *filename);
int m_interface_handle(m_interface_t *i, char *fn_name, cJSON* handler(struct m_node *node, cJSON *header, cJSON *args));
int m_interface_validate(m_interface_t *i, char *fn_name);
cJSON *(*m_interface_handler(m_interface_t *i, char *fn_name))(struct m_node *, cJSON *, cJSON *);
int m_interface_ready(m_interface_t *i);
char *m_interface_string(m_interface_t *i);

#endif
