enum m_errno {
  VERSION_MISMATCH,
  VALIDATION_ERROR,
  SERIALISATION_ERROR,
  INVALID_INTERFACE,
  BAD_ARGUMENT
};

struct m_error {
  char *message;
  m_errno_t error;
};

m_error_t *m_error_new(char *message, m_errno error){
  m_error_t *e = malloc(sizeof(m_error_t));
  e->message = message;
  e->error = error;
  return e;
}
