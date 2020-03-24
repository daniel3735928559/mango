#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "cJSON/cJSON.h"
#include "serialiser.h"

struct m_serialiser {
  char *version;
  char *method;
};

m_serialiser_t *m_serialiser_new(char *version){
  m_serialiser_t *s = malloc(sizeof(m_serialiser_t));
  s->version = version;
  s->method = "json";
  return s;
}

char *m_serialiser_serialise(m_serialiser_t *s, cJSON *header, cJSON *args){
  char *header_str = cJSON_PrintUnformatted(header);
  char *args_str = cJSON_PrintUnformatted(args);
  int header_len = strlen(header_str);
  int args_len = strlen(args_str);
  char *data = malloc(header_len + 1 + args_len + 1);
  strcpy(data, header_str);
  data[header_len] = '\n';
  strcpy(data+header_len+1, args_str);
  data[header_len + 1 + args_len] = '\0';
  
  free(header_str);
  free(args_str);
  return data;
}

cJSON **m_serialiser_deserialise(m_serialiser_t *s, char *data){
  int l = strlen(data);
  int brk = -1;
  for(int i = 0; i < l; i++) {
    if (data[i] == '\n') {
      brk = i;
      break;
    }
  }
  data[brk] = '\0';
  
  cJSON *header = cJSON_Parse(data);
  cJSON *args = cJSON_Parse(data+brk+1);
  cJSON **ans = malloc(sizeof(cJSON *)*2);
  ans[0] = header;
  ans[1] = args;
  return ans;
}
