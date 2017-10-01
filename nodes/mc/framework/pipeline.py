import io, re, time, traceback
from . transform import *
from . route import *

class Pipeline(Route):
    def __init__(self, route_id, start, end, group, source_code, steps):
        super().__init__(route_id, start, [], end, group, source_code)
        i = 0
        self.routes = []
        src = steps[0][1]
        transforms = []
        dst = None
        for x in steps[1:]:
            if x[0] == 'node':
                self.routes.append({"src":src,"dst":x[1],"transforms":transforms})
                transforms = []
                src = x[1]
            else:
                transforms.append(x[1])

        print("NP",self.routes)
        self.length = len(self.routes)
        self.outstanding = {}

    def has_src(self, n):
        for x in self.routes:
            print("SRC",x['src'],n)
            if x['src'] == n:
                return True
        return False
        
    def send(self,src,message,header,args):
        mid = header['mid']
        if not mid in self.outstanding and src == self.routes[0]['src']:
            # We haven't seen this message before--start a new pipeline:
            self.outstanding[mid] = 0
        elif mid in self.outstanding and src == self.routes[self.outstanding[mid]]['src']:
            # This message is in the middle of a pipeline
            pass
        else:
            # This message isn't actually part of the current
            # pipeline, but may be part of other routes, so just
            # return
            return

        route = self.routes[self.outstanding[mid]]
        env = {"name":header["name"]} # No raw in pipelines
        data = args        
        for t in route['transforms']:
            if t.kind == 'filter':
                if t.evaluate(env, data):
                    continue
                else:
                    # Message got filtered; drop it from pipeline
                    del self.outstanding[mid]
            else:
                data = t.evaluate(env, data)
                env = t.env
                
        header = {"name":env.get('name',''),'mid':header['mid']}
        args = data
        
        print("PIPELINE send",route['dst'].node_id)
        route['dst'].handle(header,args)

        self.outstanding[mid] += 1
        if self.outstanding[mid] >= len(self.routes):
            del self.outstanding[mid]

    def transform_spec(self):
        return " > ".join([str(t) for t in self.transforms])
            
    def __repr__(self):
        name = "{}/{}".format(self.group, self.route_id)
        if len(self.transforms) > 0:
            spec = "{} > {} > {}".format(str(self.src), self.transform_spec(), str(self.dst))
        else:
            spec = "{} > {}".format(str(self.src), str(self.dst))
        return "{}: {}".format(name, spec)
