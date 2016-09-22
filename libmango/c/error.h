enum m_errno {
  VERSION_MISMATCH,
  VALIDATION_ERROR,
  SERIALISATION_ERROR,
  INVALID_INTERFACE,
  BAD_ARGUMENT
}

typedef struct m_error {
  char *message;
  m_errno error;
} m_error_t;

m_error_t *m_error_new(char *message, m_errno error);
