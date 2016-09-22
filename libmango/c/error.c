m_error_t *m_error_new(char *message, m_errno error){
  m_error_t *e = malloc(sizeof(m_error_t));
  e->message = message;
  e->error = error;
  return e;
}
