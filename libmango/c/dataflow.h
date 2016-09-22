typedef struct m_dataflow {
  m_interface_t *interface;
  m_transport_t *transport;
  m_serialiser_t *serialiser;
  m_args_t *dispatch(m_header_t *, m_args_t *);
  void error(m_error_t *);
} m_dataflow_t;

m_dataflow_t *m_dataflow_new(m_interface_t *interface,
			     m_transport_t *transport,
			     m_serialiser_t *serialiser,
			     m_args_t *dispatch(m_header_t *, m_args_t *),
			     void error(m_error_t *));

void m_dataflow_send(m_dataflow_t *d, m_header_t *header, m_args_t *args);
    
void m_dataflow_recv(m_dataflow_t *d);
