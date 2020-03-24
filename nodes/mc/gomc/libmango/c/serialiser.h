#ifndef M_SERIALISER_H_
#define M_SERIALISER_H_

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "cJSON/cJSON.h"

#define LIBMANGO_PREAMBLE "MANGO"

typedef struct m_serialiser m_serialiser_t;

m_serialiser_t *m_serialiser_new(char *version);
char *m_serialiser_serialise(m_serialiser_t *s, cJSON *header, cJSON *args);
cJSON **m_serialiser_deserialise(m_serialiser_t *s, char *data);

#endif
