#include <stdio.h>
#include "excite.h"
#include "libmango.h"
#include "cJSON/cJSON.h"

int main(int argc, char **argv){
  m_node_t *node = m_node_new(0);
  m_node_add_interface(node, "./excite.yaml");
  m_node_handle(node, "excite", excite);
  m_node_run(node);
}

cJSON *excite(m_node_t *node, cJSON *header, cJSON *args){
  cJSON *ans = cJSON_CreateObject();
  char *s = cJSON_GetObjectItem(args,"str")->valuestring;
  unsigned long l = strlen(s);
  char *excited = malloc(l+2);
  sprintf(excited, "%s!",cJSON_GetObjectItem(args,"str")->valuestring);
  cJSON_AddStringToObject(ans,"excited",excited);
  return ans;
}
