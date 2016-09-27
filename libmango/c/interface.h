#include <stdlib.h>
#include "dict.h"
#include "cJSON/cJSON.h"

#define LIBMANGO_DEFAULT_INTERFACE_SIZE 64

enum storage_flags{VAR, VAL, SEQ};

typedef struct m_function{
  cJSON *args;
  cJSON *handler(m_node_t *node, cJSON *header, cJSON *args);
} m_function_t;

typedef struct m_interface{
  cJSON *interface;
  m_dict_t *handlers;
  int implemented;
  int size;
} m_interface_t;

m_interface_t *m_interface_new();
int m_interface_load(m_interface_t *i, char *filename);
int m_interface_handle(m_interface_t *i, char *fn_name, cJSON* handler(m_node_t *node, mcJSON *header,cJSON *args));
int m_interface_validate(m_interface_t *i, char *fn_name);
int m_interface_handler(m_interface_t *i, char *fn_name);
int m_interface_ready(m_interface_t *i);
