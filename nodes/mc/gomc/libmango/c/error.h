#ifndef M_ERROR_H_
#define M_ERROR_H_

enum m_errno {
  UNKNOWN_COMMAND,
  VERSION_MISMATCH,
  VALIDATION_ERROR,
  SERIALISATION_ERROR,
  INVALID_INTERFACE,
  BAD_ARGUMENT
};

typedef enum m_errno m_errno_t;
typedef struct m_error m_error_t;

m_error_t *m_error_new(char *message, m_errno_t error);

#endif
