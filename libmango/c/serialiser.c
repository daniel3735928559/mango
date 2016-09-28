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

void m_serialiser_make_preamble(m_serialiser_t *s, char *buf){
  sprintf(buf, "MANGO%s %s\n", s->version, s->method);
}

int m_serialiser_len_preamble(m_serialiser_t *s){
  return strlen(LIBMANGO_PREAMBLE) + strlen(s->version) + strlen(s->method) + 2;
}

char *m_serialiser_parse_preamble(m_serialiser_t *s, char *data){
  int l = m_serialiser_len_preamble(s);
  char *sample_preamble = malloc(m_serialiser_len_preamble(s)+1);
  m_serialiser_make_preamble(s, sample_preamble);
  char *ans = NULL;
  if(strncmp(sample_preamble, data, l) == 0){
    ans = data+l+1;
  }
  free(sample_preamble);
  return ans;
}

char *m_serialiser_serialise(m_serialiser_t *s, cJSON *header, cJSON *args){
  cJSON *data_dict = cJSON_CreateObject();
  cJSON_AddItemToObject(data_dict, "header", header);
  cJSON_AddItemToObject(data_dict, "argsr", args);

  int l = m_serialiser_len_preamble(s);
  char *content = cJSON_Print(data_dict);
  char *data = malloc(strlen(content)+l+1);
  m_serialiser_make_preamble(s, data);
  strcpy(data+l, content);
  free(content);
  cJSON_Delete(data_dict);
  return data;
}

cJSON *m_serialiser_deserialise(m_serialiser_t *s, char *data){
  char *content = m_serialiser_parse_preamble(s, data);
  if(!content) return NULL;
  return cJSON_Parse(content);
}
