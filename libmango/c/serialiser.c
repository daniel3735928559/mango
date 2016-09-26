#include <strlib.h>
#include <stdio.h>
#include <string.h>

char *LIBMANGO_PREAMBLE = "MANGO";

typedef struct m_serialiser {
  char *version;
  char *method;
} m_serialiser_t;

m_serialiser_t *m_serialiser_new(char *version){
  m_serialiser_t *s = malloc(sizeof(m_serialiser_t));
  s->version = version;
  s->method = "json";
}

char *m_serialiser_make_preamble(m_serialiser_t *s, char *buf){
  sprintf(buf, "MANGO%s %s\n", s->version, s->method);
}

int m_serialiser_len_preamble(m_serialiser_t *s){
  return strlen(LIBMANGO_PREAMBLE) + strlen(s->version) + strlen(s->method) + 2;
}

char *m_serialiser_parse_preamble(m_serialiser_t *s, char *data){
  int l = m_serialiser_len_preamble(s);
  char *sample_preamble = malloc(m_serialiser_len_preamble(s)+1);
  m_serialiser_make_preamble(s, sample_preamble);
  char *ans = NULL;
  if(strncmp(sample_preamble, data, l) == 0){
    ans data+l+1;
  }
  free(sample_preamble);
  return ans;
}

char *m_serialiser_serialise(m_serialiser_t *s, m_dict_t *header, m_dict_t *args){
  m_dict_t *data_dict = m_dict_new();
  m_dict_set(data_dict, "header", header);
  m_dict_set(data_dict, "args", args);

  int l = m_serialiser_len_preamble(s);
  char *content = m_dict_str(data_dict);
  char *data = malloc(strlen(content)+l+1);
  m_serialiser_make_preamble(s, data);
  strcopy(data+l, content);
  free(content);
  m_dict_free(data_dict);
  return data;
}

m_dict_t *m_serialiser_deserialise(m_serialiser_t *s, char *data){
  char *content = m_serialiser_validate_preamble(s, data);
  if(!content) return NULL;
  
}

    this.deserialise = function(data){
	try{
	    var d = JSON.parse(message);
	    return [d['header'],d['args']]
	} catch (e) {
	    console.log(e);
	    throw new MError("Failed to parse message");
	}
    }
}
