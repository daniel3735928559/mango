#include <stdio.h>
#include "excite.h"
#include "libmango.h"
#include "cJSON/cJSON.h"

int main(int argc, char **argv){
  printf("hello from ex.c\n");
  m_node_t *node = m_node_new(1);
  m_node_handle(node, "excite", excite);
  m_node_run(node);
}

cJSON *excite(m_node_t *node, cJSON *args, m_result_t *result){
  char *str = cJSON_GetObjectItem(args,"message")->valuestring;
  unsigned long len = strlen(str);
  char *excited = malloc(len+2);
  snprintf(excited, len+2, "%s!",str);
  cJSON_AddStringToObject(result->data,"message",excited);
  result->command = "excited";
}
