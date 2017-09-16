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

cJSON *excite(m_node_t *node, cJSON *header, cJSON *args, m_result_t *result){
  char *str = cJSON_GetObjectItem(args,"str")->valuestring;
  unsigned long len = strlen(str);
  char *excited = malloc(len+2);
  snprintf(excited, len+2, "%s!",str);
  cJSON_AddStringToObject(result->data,"str",excited);
  result->name = "excited";
  result->data = ans;
}
