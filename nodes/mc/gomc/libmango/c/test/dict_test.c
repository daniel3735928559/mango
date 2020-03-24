#include "dict.h"
#include <stdio.h>

#define NONE "Nothing"

int main(int argc, char **argv){
  m_dict_t *d = m_dict_new(0);
  printf("%d\n",m_dict_set(d, "a", "b"));
  printf("%d\n",m_dict_set(d, "a", "blah"));
  printf("%d\n",m_dict_set(d, "aasd", "foo"));
  printf("%d\n",m_dict_set(d, "adw", "fo4o"));
  printf("%d\n",m_dict_set(d, "aqw", "foo3"));
  printf("%d\n",m_dict_set(d, "afew", "foo2"));
  char *s1 = m_dict_get(d, "a");
  char *s2 = m_dict_get(d, "afew");
  char *s3 = m_dict_get(d, "askdj");
  printf("%s\n", s1 ? s1 : NONE);
  printf("%s\n", s2 ? s2 : NONE);
  printf("%s\n", s3 ? s3 : NONE);
}
