from hashlib import *
import hmac, re, json
from error import *

HMAC_LENGTH = 32

class m_serialiser:
    def __init__(self,version,method="json"):
        self.version = version
        self.method = method
        self.serialisers = {"json":m_json_serialiser()}

    def serialise(self,header,msg):
        return "{}\n{}".format(json.dumps(header), json.dumps(msg))
        # return self.make_preamble()+self.serialisers[self.method].pack(header,msg)

    def deserialise(self,message):
        
        header_str, args_str = message.decode().split("\n", 1)
        return json.loads(header_str), json.loads(args_str)
        # ver,method,msg = self.parse_preamble(message)
        # if ver != self.version:
        #     raise m_error(m_error.VERSION_MISMATCH, ver + " given, " + self.version + " expected")
        # h,a = self.serialisers[method].unpack(msg)
        # return h,a

class m_json_serialiser:
    def pack(self,header,msg):
        d = {"header":header,"args":msg}
        return bytes(json.dumps(d),"ASCII")

    def unpack(self,message):
        d = json.loads(message.decode())
        return d['header'],d['args']
