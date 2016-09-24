#define DICT_INIT_SIZE 100

typedef struct m_dict_entry {
  char *key;
  void *val;
} m_dict_entry_t;

typedef struct m_dict {
  m_dict_entry_t **data;
  int size;
} m_dict_t;

m_dict_t *m_dict_new(int size);
void *m_dict_get(m_dict_t *dict, char *key);
int m_dict_set(m_dict_t *dict, char *key, void *val);
m_dict_entry_t *m_dict_next(m_dict_t *dict, int idx);
void m_dict_free(m_dict_t *dict);
