#include <stdlib.h>
#include "dict.h"

#define LIBMANGO_DEFAULT_INTERFACE_SIZE 64

typedef struct m_function{
  m_dict_t *args;
  m_dict_t *handler(m_node_t *node, m_dict_t *header, m_dict_t *args);
} m_function_t;

typedef struct m_interface{
  m_dict_t *interface;
  int implemented;
  int size;
} m_interface_t;

m_interface_t *m_interface_new();
int m_interface_load(m_interface_t *i, char *filename);
int m_interface_handle(m_interface_t *i, char *fn_name, m_dict_t* handler(m_node_t *node, m_dict_t *header, m_dict_t *args));
int m_interface_validate(m_interface_t *i, char *fn_name);
int m_interface_handler(m_interface_t *i, char *fn_name);
int m_interface_ready(m_interface_t *i);
