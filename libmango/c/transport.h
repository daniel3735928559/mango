#ifndef M_TRANSPORT_H_
#define M_TRANSPORT_H_

typedef struct m_transport {
  char *target;
  void *socket;
} m_transport_t;

m_transport_t *m_transport_new(char *addr, void *context);
void m_transport_tx(m_transport_t *t, char *data);
char *m_transport_rx(m_transport_t *t);
void m_transport_free(m_transport_t *t);

#endif
