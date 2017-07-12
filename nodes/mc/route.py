import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from transform import *
from mc_dataflows import *
from dataflow import m_dataflow
from transport import *
from libmango import m_node

class Route:
    def __init__(self, start, transforms, end, source_code):
        self.source_code = source_code
        self.start = nodelist[0]
        self.end = nodelist[-1]
        self.transforms = transforms

    def to_string(self):
        print(self.source_code)
        
    def apply(self,message,header,args):
        data = args
        env = {"raw":message,"name":header["name"]}
        
        for t in self.transforms:
            if t[0] == 'filter':
                if t[1].evaluate(env, data): continue
                else: return None,None,None
            else:
                  new_header = t[1].evaluate(env, data)
                  env = t[1].env
                  if env['raw'] != message:
                      return env['raw'],None,None
        new_message = message
        new_header = {"name":env.get('name','')}
        new_args = data
        return None, new_header, new_args

    def send(self,message,header,args):
        m,h,a = self.apply(message,header,args)
        if h is None and m is None: return
        if m is None:
            print("ROUTE send",m,h,a,self.end.node_id)

            # Special case so that mc can reply directly
            
            if str(self.end.node_id) == "mc":
                self.end.send(h,a,self.start.route)
            else:
                self.end.send(h,a)

        else:
            print("R SEND RAW",self.end,self.end.dataflow)
            self.end.dataflow.send_raw(m,bytearray(self.endpoint.owner.node_id,'utf-8'))

    def __repr__(self):
        return str(self.source) + " > " + "".join([str(t)+" > " for t in self.transmogrifiers]) + str(self.endpoint)

    # def reply(self,port,header,reply,raw,route=None):
    #     print("MC REPLY",port,header,reply,raw,route)
    #     self.source.owner.conn.send(a,port=self.source.name,dest=self.endpoint.get_id(),route=bytearray(self.source.owner.node_id,'utf-8'),source_node=header['source'],mid=header['mid'])
