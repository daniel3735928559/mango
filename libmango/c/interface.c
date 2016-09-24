#include <stdlib.h>
#include "dict.h"

m_interface_t *m_interface_new(){
  m_interface_t *i = malloc(sizeof(m_interface_t));
  i->interface = m_dict_new(LIBMANGO_DEFAULT_INTERFACE_SIZE);
  i->implemented = 0;
  i->size = 0;
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
  i->implemented++;
  return 0;
}

int m_interface_validate(m_interface_t *i, char *fn_name){
  return m_dict_get(i->interface, fn_name) == NULL ? 0 : 1;
}

int m_interface_handler(m_interface_t *i, char *fn_name){
  return m_dict_get(i->interface, fn_name)->handler;
}

int m_interface_ready(m_interface_t *i){
  return i->implemented == i->size; // Check for functions not yet implemented
}

char *m_interface_str(m_interface_i *i){
  int len = 0;
  int idx = 0;
  int k = 0;
  char **keys = malloc(sizeof(char *)*i->size);
  char **vals = malloc(sizeof(char *)*i->size);
  while(idx = m_dict_next(i->interface, idx) != -1){
    m_dict_entry_t *e = interface->data[idx];
    keys[k] = strdup(e->key);
    vals[k] = strdup((char *)(e->val));
    len += strlen(keys[k]) + strlen(vals[k]);
    idx++;
    k++;
  }
  char *ans = malloc(len+4*i->size/*size of quotes and curly brackets??*/);
}
