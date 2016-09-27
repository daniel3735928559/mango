#include "dict.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct m_dict_entry {
  char *key;
  void *val;
};

struct m_dict {
  m_dict_entry_t **data;
  int size;
};

m_dict_t *m_dict_new(int size){
  m_dict_t *d = malloc(sizeof(m_dict_t));
  if(!d) return NULL;
  d->data = malloc(sizeof(char *)*(size ? size : DICT_INIT_SIZE));
  if(!d->data) return NULL;
  d->size = DICT_INIT_SIZE;
  return d;
}

int m_dict_hash(m_dict_t *dict, char *key){
  long h = 0;
  long b = 2;
  while(*key) h = b*h + *key++;
  return h % dict->size;
}

void *m_dict_get(m_dict_t *dict, char *key){
  long h = m_dict_hash(dict, key);
  //printf("h=%d\n",h);
  int i = 0;
  while(i++ < dict->size && dict->data[h] && strcmp(dict->data[h]->key, key))
    h = (h+1)%(dict->size);
  if(dict->data[h] == NULL) return NULL;
  if(i >= dict->size) return NULL;
  return dict->data[h]->val;
}

int m_dict_set(m_dict_t *dict, char *key, void *val){
  long h = m_dict_hash(dict, key);
  //printf("h=%d\n",h);
  int i = 0;
  while(i++ < dict->size && dict->data[h] && strcmp(dict->data[h]->key, key)){
    h = (h+1)%(dict->size);
  }
  if(i >= dict->size);
  dict->data[h] = malloc(sizeof(m_dict_entry_t));
  dict->data[h]->key = strdup(key);
  dict->data[h]->val = val;
  return 0;
}

/* int m_dict_next(m_dict_t *dict, int idx){ */
/*   while(idx < dict->size && dict->data[idx] == NULL) idx++; */
/*   if(idx >= dict->size) return -1; */
/*   return idx; */
/* } */

void m_dict_expand(m_dict_t *dict){
  
}

void m_dict_free(m_dict_t *dict){
  int i;
  for(i = 0; i < dict->size; i++){
    if(dict->data[i] != NULL){
      free(dict->data[i]->key);
      free(dict->data[i]->val);
      free(dict->data[i]);
    }
  }
  free(dict);
}
