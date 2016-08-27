class m_error(Exception):
    VERSION_MISMATCH = 0
    VALIDATION_ERROR = 1
    SERIALISATION_ERROR = 2
    INVALID_INTERFACE = 3
    BAD_ARGUMENT = 4
    reprs = {VERSION_MISMATCH:"Version mismatch",
             VALIDATION_ERROR:"Validation error",
             SERIALISATION_ERROR:"Serialisation error",
             INVALID_INTERFACE:"Invalid interface",
             BAD_ARGUMENT:"Bad argument"}
    
    def __init__(self,code,message):
        self.code = code
        self.message = message
        Exception.__init__(self,message)

    def __repr__(self):
        return m_error.reprs[self.code]+': '+self.message
