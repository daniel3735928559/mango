#include <stdio.h>
#include <stdlib.h>
#include <yaml.h>
#include "dict.h"
#include "cJSON/cJSON.h"

struct m_function{
  cJSON *args;
  cJSON *(*handler)(m_node_t *node, cJSON *header, cJSON *args);
};

struct m_interface{
  cJSON *interface;
  m_dict_t *handlers;
  int implemented;
  int size;
};

m_interface_t *m_interface_new(){
  m_interface_t *i = malloc(sizeof(m_interface_t));
  i->interface = cJSON_CreateObject;
  i->handlers = m_dict_new(LIBMANGO_DEFAULT_INTERFACE_SIZE);
  i->implemented = 0;
  i->size = 0;
  return i;
}

int m_interface_load(m_interface_t *i, char *filename){
  cJSON *obj = cJSON_createObject();
  yaml_parser_t parser;
  FILE *source = fopen(filename, "rb");
  yaml_parser_initialize(&parser);
  yaml_parser_set_input_file(&parser, source);
  m_interface_process_yaml(&parser, obj);
  yaml_parser_delete(&parser);
  fclose(source);
}

void m_interface_process_yaml(yaml_parser_t *parser, cJSON *node){
  cJSON *last_leaf;
  yaml_event_t event;
  int storage = VAR;
  while(1) {
    yaml_parser_parse(parser, &event);
    
    // Parse value either as a new leaf in the mapping
    //  or as a leaf value (one of them, in case it's a sequence)
    if (event.type == YAML_SCALAR_EVENT) {
      if (storage) g_node_append_data(last_leaf, g_strdup((gchar*) event.data.scalar.value));
      else last_leaf = g_node_append(data, g_node_new(g_strdup((gchar*) event.data.scalar.value)));
      storage ^= VAL; // Flip VAR/VAL switch for the next event
    }
    
    // Sequence - all the following scalars will be appended to the last_leaf
    else if (event.type == YAML_SEQUENCE_START_EVENT) storage = SEQ;
    else if (event.type == YAML_SEQUENCE_END_EVENT) storage = VAR;

    // depth += 1
    else if (event.type == YAML_MAPPING_START_EVENT) {
      process_layer(parser, last_leaf);
      storage ^= VAL; // Flip VAR/VAL, w/o touching SEQ
    }
    
    // depth -= 1
    else if(event.type == YAML_MAPPING_END_EVENT
	    || event.type == YAML_STREAM_END_EVENT)
      break;
    
    yaml_event_delete(&event);
  }
}

int m_interface_handle(m_interface_t *i, char *fn_name, cJSON *handler(m_node_t *node, cJSON *header, cJSON *args)){
  m_function_t *fn = m_dict_get(i->handlers, fn_name);
  if(!fn) return -1; // Function not in interface
  if(fn->handler) return -2; // Already implemented
  fn->handler = handler;
  i->implemented++;
  return 0;
}

int m_interface_validate(m_interface_t *i, char *fn_name){
  return m_dict_get(i->handlers, fn_name) == NULL ? 0 : 1;
}

int m_interface_handler(m_interface_t *i, char *fn_name){
  return m_dict_get(i->handlers, fn_name)->handler;
}

int m_interface_ready(m_interface_t *i){
  return i->implemented == i->size; // Check for functions not yet implemented
}

char *m_interface_string(m_interface_t *i){
  return cJSON_Print(i->interface);
}
