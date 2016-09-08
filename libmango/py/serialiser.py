from hashlib import *
import hmac, re, json
from error import *

HMAC_LENGTH = 32

class m_serialiser:
    def __init__(self,version,method="json"):
        self.version = version
        self.method = method
        self.serialisers = {"json":m_json_serialiser()}

    def make_preamble(self):
        return bytes("MANGO{0} {1}\n".format(self.version,self.method),"ASCII")

    def parse_preamble(self,msg):
        nl1 = msg.find(b'\n')
        m = re.match("^MANGO([0-9.]*) ([^ ]*)$",msg[:nl1].decode("ASCII"))
        if(m is None):
            raise m_error(m_error.SERIALISATION_ERROR,"Preamble failed to parse")
        return m.group(1),m.group(2),msg[nl1+1:]

    def serialise(self,header,msg):
        return self.make_preamble()+self.serialisers[self.method].pack(header,msg)

    def deserialise(self,message):
        ver,method,msg = self.parse_preamble(message)
        if ver != self.version:
            raise m_error(m_error.VERSION_MISMATCH, ver + " given, " + self.version + " expected")
        h,a = self.serialisers[method].unpack(msg)
        return h,a

class m_json_serialiser:
    def pack(self,header,msg):
        d = {"header":header,"args":msg}
        return bytes(json.dumps(d),"ASCII")

    def unpack(self,message):
        d = json.loads(message.decode())
        return d['header'],d['args']
