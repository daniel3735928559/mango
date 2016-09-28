#include <stdio.h>
#include "excite.h"
#include "libmango.h"
#include "cJSON/cJSON.h"

int main(int argc, char **argv){
  m_node_t *node = m_node_new(0);
  m_node_add_interface(node, "./excite.yaml");
  m_node_handle(node, "excite", excite);
  m_node_handle(node, "print", print);
  m_node_ready(node);
  m_node_run(node);
}

cJSON *excite(m_node_t *node, cJSON *header, cJSON *args){
  printf("EXCITING\n%s\n%s\n",cJSON_Print(header),cJSON_Print(args));
  cJSON *ans = cJSON_CreateObject();
  char *s = cJSON_GetObjectItem(args,"str")->valuestring;
  int l = strlen(s);
  char *excited = malloc(l+2);
  memcpy(excited,s,l);
  excited[l]='!';
  excited[l+1] = '\0';
  cJSON_AddStringToObject(ans,"excited",excited);
  return ans;
}

cJSON *print(m_node_t *node, cJSON *header, cJSON *args){
  printf("PRINT: \nHEADER = %s\nARGS = %s\n",cJSON_Print(header),cJSON_Print(args));
  return NULL;
}
