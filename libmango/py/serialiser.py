from hashlib import *
import hmac, re, json

HMAC_LENGTH = 32

class m_serialiser:
    def __init__(self,version):
        self.version = version
        self.serialisers = {"json":m_json_serialiser(self.version),
                            "std":m_std_serialiser(self.version)}

    def make_preamble(self,method,dest=None):
        if(dest is None):
            return bytes("MANGO{0} {1}\n".format(self.version,method),"ASCII")
        else:
            return bytes("MANGO{0} {1} {2}\n".format(self.version,method,dest),"ASCII")

    def parse_preamble(self,msg):
        nl1 = msg.find(b'\n')
        preamble = msg[:nl1]
        msg = msg[nl1+1:]
        m = re.match("^MANGO([0-9.]*) ([^ ]*)$",preamble.decode("ASCII"))
        if(m is None):
            raise m_error(m_error.SERIALISATION_ERROR,"Preamble failed to parse")
        return m.group(1),m.group(2),msg

    def serialise(self,method,msg_dict,source,mid,dport=None):
        return self.make_preamble(method,dport)+self.serialisers[method].pack(msg_dict,source,mid,dport)

    def deserialise(self,message):
        ver,method,msg = self.parse_preamble(message)
        if ver != self.version:
            raise m_error(m_error.VERSION_MISMATCH, ver + " given, " + self.version + " expected")
        h,a = self.serialisers[method].unpack(msg)
        return ver,h,a


    
class m_json_serialiser:
    def pack(self,msg_dict,source,mid,dport=None):
        d = {"header":{"source":source,"mid":mid},"args":msg_dict}
        return bytes(json.dumps(d),"ASCII")

    def unpack(self,message):
        d = json.loads(message.decode())
        return d['header'],d['args']



class m_std_serialiser:
    def __init__(self,version):
        self.version = version

    def dict_pack(self,msg_dict):
        msg = bytearray()
        xref = []
        pos = 0
        for a in msg_dict:
            entry = str(pos) + ":"
            pos += len(a)
            entry += str(pos) + ":"
            pos += len(msg_dict[a])
            entry += str(pos)
            xref += [entry]
            msg += bytearray(a+":","ASCII")
            msg += bytearray(msg_dict[a])
            msg += bytearray("\n","ASCII")
            
        return bytearray(" ".join(xref)+"\n","ASCII") + msg

    def dict_unpack(self,msg,binary=True):
        nl1 = msg.find(ord('\n'))
        xref = msg[:nl1].decode().split(" ")
        msg = msg[nl1+1:]
        result = {}
        for x in xref:
            s,t,e = [int(i) for i in x.split(":")]
            result[msg[s:t].decode("ASCII")] = msg[t+1:e] if binary else msg[t+1:e].decode("ASCII")
        return result
            
    def pack(self,msg_dict,source,mid,dport=None):
        h = {"source":source,"mid":mid}
        header = self.dict_pack(h)
        body = self.dict_pack(msg_dict)
        line0 = "{0} {1}\n".format(len(header), len(body))
        result = bytearray(line0,"ASCII")+header+body
        return result

    def unpack(self,msg):
        nl1 = msg.find(ord('\n'))
        if(nl1 == -1): return;
        m = re.search("^([0-9]*) ([0-9]*)$",msg[:nl1].decode("ASCII"))
        if(len(m.groups())!=2): 
            print("invalid syntax")
            return(NULL)
        header_length = int(m.group(1))
        message_length = int(m.group(2))
        
        header = msg[nl1+1:nl1+1+header_length]
        message = msg[nl1+1+header_length:nl1+1+header_length+message_length]
        if(len(header) != header_length or len(message) != message_length):
            print(len(message))
            print(len(header))
            print("invalid syntax")
            return
        return dict_unpack(header),dict_unpack(message)
