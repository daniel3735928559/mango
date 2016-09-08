# Many thanks to https://gist.github.com/rich20bb/4190781

import socket, hashlib, base64, struct

MAGICGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
TEXT = 0x01
BINARY = 0x02

class mu_server_ws:
    def __init__(self, bind, port):
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self.socket.bind((bind, port))
        self.bind = bind
        self.port = port
        self.connections = {}
        self.listeners = [self.socket]
        self.socket.listen(5)
        self.socket.setblocking(0)
        
    def die(self):
        self.socket.close()

    def tx(self,data):
        pass

    def rx(self):
        c,a = self.socket.accept()
        return c,a

class mu_client_ws:
    def __init__(self, sock):
        self.handshake = (
            "HTTP/1.1 101 Web Socket Protocol Handshake\r\n"
            "Upgrade: WebSocket\r\n"
            "Connection: Upgrade\r\n"
            "Sec-WebSocket-Accept: %(acceptstring)s\r\n"
            "Server: mu\r\n"
            "Access-Control-Allow-Origin: http://localhost\r\n"
            "Access-Control-Allow-Credentials: true\r\n"
            "\r\n"
        )
        self.sock = sock
        self.handshaken = False
        self.alive = True
        self.header = ""

    def rx(self):
        data = self.sock.recv(4096)
        if not self.handshaken:
            print("[DEBUG] No handshake yet")
            self.header += data.decode('ASCII')
            if self.header.find('\r\n\r\n') != -1:
                parts = self.header.split('\r\n\r\n', 1)
                self.header = parts[0]
                if self.do_handshake_(self.header, parts[1]):
                    print("[INFO] Handshake successful")
                    self.handshaken = True
            return None
        else:
            print("DATA",data)
            if(len(data) > 0):
                return self.decode_(data)
            return None
        
    def tx(self, s):
        message = bytes([])
        b1 = 0x80
        if type(s) == str:
            b1 |= TEXT
            payload = s.encode('utf-8')
        else:
            payload = s
            
        # Append 'FIN' flag to the message
        message += bytes([b1])

        # never mask frames from the server to the client
        b2 = 0
        
        # How long is our payload?
        length = len(payload)
        if length < 126:
            b2 |= length
            message += bytes([b2])
        
        elif length < (2 ** 16) - 1:
            b2 |= 126
            message += bytes([b2])
            l = struct.pack(">H", length)
            message += l
        
        else:
            l = struct.pack(">Q", length)
            b2 |= 127
            message += bytes([b2])
            message += l

        # Append payload to message
        message += payload

        # Send to the client
        print("ASA",message)
        self.sock.send(message)


    def decode_(self, byte_array):
        # Turn string values into opererable numeric byte values
        #byte_array = [character for character in str_stream]
        data_length = byte_array[1] & 127
        index_first_mask = 2

        if data_length == 126:
            index_first_mask = 4
        elif data_length == 127:
            index_first_mask = 10

        # Extract masks
        masks = [m for m in byte_array[index_first_mask : index_first_mask+4]]
        index_first_data_byte = index_first_mask + 4
        
        # List of decoded characters
        decoded = []
        i = index_first_data_byte
        j = 0
        
        # Loop through each byte that was received
        while i < len(byte_array):
        
            # Unmask this byte and add to the decoded buffer
            decoded.append(byte_array[i] ^ masks[j % 4])
            i += 1
            j += 1

        # Return the decoded string
        return bytes(decoded)

    def do_handshake_(self, header, key=None):    
        print("[DEBUG] Begin handshake: %s" % header)
        handshake = self.handshake
        
        # Step through each header
        for line in header.split('\r\n')[1:]:
            name, value = line.split(': ', 1)
            
            # If this is the key
            if name.lower() == "sec-websocket-key":
            
                # Append the standard GUID and get digest
                combined = value + MAGICGUID
                response = base64.b64encode(hashlib.sha1(combined.encode('utf-8')).digest())
                
                # Replace the placeholder in the handshake response
                handshake = handshake % { 'acceptstring' : response.decode() }

        print("[DEBUG] Sending handshake %s" % handshake)
        self.sock.send(bytes(handshake,'utf-8'))
        return True

    def close(self):
        self.alive = False
        self.sock.close()
