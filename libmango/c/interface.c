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
  int implemented;
  int size;
};

m_interface_t *m_interface_new(){
  m_interface_t *i = malloc(sizeof(m_interface_t));
  i->interface = cJSON_CreateObject();
  i->handlers = m_dict_new(LIBMANGO_DEFAULT_INTERFACE_SIZE);
  i->implemented = 0;
  i->size = 0;
  return i;
}

cJSON *m_interface_spec(m_interface_t *i){
  return i->interface;
}

cJSON *m_interface_process_yaml(yaml_parser_t *parser){
  cJSON *current = cJSON_CreateObject();
  char *current_key = NULL;
  yaml_event_t event;
  int storage = VAR;
  while(1) {
    yaml_parser_parse(parser, &event);

    if (event.type == YAML_SCALAR_EVENT) {
      if(storage){
	//printf("STORE EV=%s key=%s, CUR=%s\n", event.data.scalar.value, current_key, cJSON_Print(current));
	cJSON *o = cJSON_CreateString(event.data.scalar.value);
	cJSON_AddItemToObject(current, current_key, o);
      }
      else{
	//printf("MAP EV=%s\n", event.data.scalar.value);
	current_key = strdup(event.data.scalar.value);
      }
      storage ^= VAL;
    }
    
    else if (event.type == YAML_SEQUENCE_START_EVENT){
      cJSON *arr = cJSON_CreateArray();
      cJSON_AddItemToObject(current, event.data.scalar.value, arr);
      storage = SEQ;
    }

    else if (event.type == YAML_SEQUENCE_END_EVENT){
      storage = VAR;
    }

    else if (event.type == YAML_MAPPING_START_EVENT) {
      //printf("START MAP\n");
      if(current_key == NULL || strcmp(current_key,"") == 0)
	return m_interface_process_yaml(parser);
      else
	cJSON_AddItemToObject(current, current_key, m_interface_process_yaml(parser));
      storage ^= VAL;
    }
    
    else if(event.type == YAML_MAPPING_END_EVENT || event.type == YAML_STREAM_END_EVENT){
      //printf("END MAP\n");
      return current;
    }
    
    yaml_event_delete(&event);
  }
}

void m_interface_load(m_interface_t *i, char *filename){
  yaml_parser_t parser;
  FILE *source = fopen(filename, "rb");
  yaml_parser_initialize(&parser);
  yaml_parser_set_input_file(&parser, source);
  cJSON *obj = m_interface_process_yaml(&parser);
  yaml_parser_delete(&parser);
  fclose(source);
  char *name = cJSON_GetObjectItem(obj, "name")->valuestring;
  cJSON *new_if = cJSON_CreateObject();
  cJSON *new_inputs = cJSON_CreateObject();
  cJSON *new_outputs = cJSON_CreateObject();
  
  if(cJSON_HasObjectItem(obj, "inputs")){
    cJSON *o = cJSON_GetObjectItem(obj, "inputs")->child;
    while(o){
      cJSON_AddItemToObject(new_inputs, o->string, cJSON_Duplicate(o,1));
      o = o->next;
    }
  }
  if(cJSON_HasObjectItem(obj, "outputs")){
    cJSON *o = cJSON_GetObjectItem(obj, "outputs")->child;
    while(o){
      cJSON_AddItemToObject(new_outputs, o->string, cJSON_Duplicate(o,1));
      o = o->next;
    }
  }
  cJSON_AddItemToObject(new_if, "inputs", new_inputs);
  cJSON_AddItemToObject(new_if, "outputs", new_outputs);
  cJSON_AddItemToObject(i->interface, name, new_if);
}

int m_interface_handle(m_interface_t *i, char *fn_name, cJSON *handler(m_node_t *node, cJSON *header, cJSON *args)){
  //int present = cJSON_HasObjectItem(i->interface, fn_name);
  void *fn = m_dict_get(i->handlers, fn_name);
  //if(!present) return -1; // Function not in interface
  if(fn) return -2; // Already implemented
  m_dict_set(i->handlers, fn_name, handler);
  i->implemented++;
  return 0;
}

int m_interface_validate(m_interface_t *i, char *fn_name){
  return m_dict_get(i->handlers, fn_name) == NULL ? 0 : 1;
}

cJSON *(*m_interface_handler(m_interface_t *i, char *fn_name))(struct m_node *, cJSON *, cJSON *){
  return m_dict_get(i->handlers, fn_name);
}

int m_interface_ready(m_interface_t *i){
  return i->implemented == i->size; // Check for functions not yet implemented
}

char *m_interface_string(m_interface_t *i){
  return cJSON_PrintUnformatted(i->interface);
}
