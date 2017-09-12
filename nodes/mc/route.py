import io, re, socketserver, socket, threading, time, signal, os, sys, random, zmq, subprocess, shlex, json, traceback
from transform import *
from mc_dataflows import *
from dataflow import m_dataflow
from transport import *
from libmango import m_node

class Route:
    def __init__(self, start, transforms, end, group, source_code):
        self.source_code = source_code
        self.src = start
        self.dst = end
        self.src_name = str(start)
        self.dst_name = str(end)
        self.group = group
        self.transforms = transforms

    def apply(self,message,header,args):
        data = args
        env = {"raw":message,"name":header["name"]}
        
        for t in self.transforms:
            if t.kind == 'filter':
                if t.evaluate(env, data): continue
                else: return None,None,None
            else:
                data = t.evaluate(env, data)
                env = t.env
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
            print("ROUTE send",m,h,a,self.dst.node_id)

            # Special case so that mc can reply directly
            
            if str(self.dst.node_id) == "mc":
                self.dst.handle(h,a,self.src.route)
            else:
                self.dst.handle(h,a)

        else:
            print("R SEND RAW",self.dst,self.dst.dataflow)
            self.dst.dataflow.send_raw(m,bytearray(self.dst.node_id,'utf-8'))

    def __repr__(self):
        return "{} > {} > {}".format(self.src.node_id, str(self.transforms), self.dst.node_id)
