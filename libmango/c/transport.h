typedef struct m_transport {
  char *target;
  zmq_socket *socket;
} m_transport_t;

void m_transport_tx(m_transport_t *t, char *data);

char *m_transport_rx();
