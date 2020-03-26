#include <stdio.h>
#include <stdlib.h>
#include <yaml.h>
#include "libmango.h"
#include "dict.h"
#include "cJSON/cJSON.h"

struct m_function{
  cJSON *(*handler)(m_node_t *node, cJSON *header, cJSON *args);
};

struct m_interface{
  cJSON *interface;
  m_dict_t *handlers;
};

m_interface_t *m_interface_new(){
  m_interface_t *i = malloc(sizeof(m_interface_t));
  i->handlers = m_dict_new(LIBMANGO_DEFAULT_INTERFACE_SIZE);
  return i;
}

int m_interface_handle(m_interface_t *i, char *fn_name, cJSON *handler(m_node_t *node, cJSON *args, m_result_t *result)){
  void *fn = m_dict_get(i->handlers, fn_name);
  if(fn) return -2; // Already implemented
  m_dict_set(i->handlers, fn_name, handler);
  return 0;
}

int m_interface_validate(m_interface_t *i, char *fn_name){
  return m_dict_get(i->handlers, fn_name) == NULL ? 0 : 1;
}

cJSON *(*m_interface_handler(m_interface_t *i, char *fn_name))(struct m_node *, cJSON *, m_result_t *result){
  return m_dict_get(i->handlers, fn_name);
}

char *m_interface_string(m_interface_t *i){
  return cJSON_PrintUnformatted(i->interface);
}
