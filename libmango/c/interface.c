#include <stdlib.h>
#include "dict.h"

m_interface_t *m_interface_new(){
  m_interface_t *i = malloc(sizeof(m_interface_t));
  i->interface = m_dict_new(LIBMANGO_DEFAULT_INTERFACE_SIZE);
  i->unimplemented = 0;
  return i;
}

int m_interface_load(m_interface_t *i, char *filename){
  // load YAML and set handlers to NULL
}

int m_interface_handle(m_interface_t *i, char *fn_name, m_dict_t* handler(m_node_t *node, m_dict_t *header, m_dict_t *args)){
  m_function_t *fn = m_dict_get(i->interface, fn_name);
  if(!fn) return -1; // Function not in interface
  if(fn->handler) return -2; // Already implemented
  fn->handler = handler;
  i->unimplemented--;
  return 0;
}

int m_interface_validate(m_interface_t *i, char *fn_name){
  return m_dict_get(i->interface, fn_name) == NULL ? 0 : 1;
}

int m_interface_handler(m_interface_t *i, char *fn_name){
  return m_dict_get(i->interface, fn_name)->handler;
}

int m_interface_ready(m_interface_t *i){
  return i->unimplemented == 0; // Check for functions not yet implemented
}
