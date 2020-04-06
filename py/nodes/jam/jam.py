from IPython import embed

import dbus, re, sys, time, os, threading
import zmq
from libmango import *
from fuzzywuzzy import process
from gi.repository import GLib
from dbus.mainloop.glib import DBusGMainLoop


class jam(m_node):
    def __init__(self):
        super().__init__(debug=True)
        self.interface.add_interface({'finduser':self.finduser,'sendto':self.sendto,'sendconv':self.sendconv,'getconvs':self.getconvs,'getbuddies':self.getbuddies})
        # Setup socket for receiving messages
        self.purple_sock = self.context.socket(zmq.ROUTER)
        self.purple_sock.bind("inproc://purplerecv")
        self.add_socket(self.purple_sock, self.purple_recv, self.purple_err)
        self.purple_thread = threading.Thread(target=run_recv_thread, args=(self.context,))
        self.purple_thread.daemon = True
        self.purple_thread.start()
        #process = subprocess.Popen(["python", "purple.py"], cwd=os.path.dirname(os.path.realpath(__file__)), env=os.environ)
        
        # Get DBus instance for sending messages
        bus = dbus.SessionBus()
        obj = bus.get_object("im.pidgin.purple.PurpleService", "/im/pidgin/purple/PurpleObject")
        self.purple = dbus.Interface(obj, "im.pidgin.purple.PurpleInterface")

        self.accounts = self.purple.PurpleAccountsGetAll()
        self.buddies = []
        for acc in self.accounts:
            bids = self.purple.PurpleFindBuddies(acc, "")
            
            for bid in bids:
                buddy = {
                    "id": bid,
                    "name": self.purple.PurpleBuddyGetName(bid),
                    "alias": self.purple.PurpleBuddyGetAlias(bid),
                    "account_id": acc
                }
                self.buddies.append(buddy)
        self.match_quality_threshold = 95
        self.debug_print("jamming")
        self.run()

    def buddy_by_name(self, name):
        for b in self.buddies:
            if b["name"] == name:
                return "{} ({}) @{}".format(b["name"], b["alias"], b["account_id"])
        return "[buddy not found]"
        
    def purple_recv(self):
        rt,msg = self.purple_sock.recv_multipart()
        msg= json.loads(msg)
        self.debug_print("[JAM.PY] MSG",msg)
        cmd = msg["command"]
        del msg["command"]
        self.m_send(cmd, msg)
        
    def purple_err(self):
        self.debug_print("PRPL DIED")
                        
    def finduser(self, args):
        candidates = process.extract({"name":args['name']},self.buddies,processor=lambda b: b['name'])
        candidates += process.extract({"alias":args['name']},self.buddies,processor=lambda b: b['alias'])
        candidates = sorted([c for c in candidates if c[1] > self.match_quality_threshold],key=lambda c: c[1])
        self.debug_print("Filtered candidates:",candidates)
        return "users",{"users":[{"id":str(c[0]["id"]),"name":c[0]["name"],"alias":c[0]["alias"]} for c in candidates]}
                        
    def sendto(self, args):
        for b in self.buddies:
            if str(b["id"]) == str(args["id"]):
                conv = self.purple.PurpleConversationNew(1, int(b["account_id"]), str(b['name']))
                im = self.purple.PurpleConvIm(conv)
                self.debug_print("[JAM.PY] SENDING",args,b,conv,im)
                self.purple.PurpleConvImSend(im, args['msg'])

    def sendconv(self, args):
        try:
            c = int(args['conv'])
            chat = self.purple.PurpleConversationGetChatData(c)
            if chat == 0:
                im = self.purple.PurpleConvIm(c)
                self.debug_print("SENDING",args,im)
                self.purple.PurpleConvImSend(im, args['msg'])
            else:
                self.purple.PurpleConvChatSend(chat, args['msg'])
        except:
            pass
        
    def getbuddies(self, args):
        bl = []
        for b in self.buddies:
            bl.append({"id":str(b["id"]),"name":b["name"],"alias":b["alias"],"account":str(b["account_id"])})
        return "buddies",{"buddylist":bl}
            
    def getconvs(self, args):
        cs = self.purple.PurpleGetConversations()
        convs = []
        for c in cs:
            chat = self.purple.PurpleConversationGetChatData(c)
            if chat == 0:
                user = self.purple.PurpleConversationGetName(c)
                convs.append({"id":str(c),"participants":[self.buddy_by_name(user)]})
            else:
                users = self.purple.PurpleConvChatGetUsers(chat)
                usernames = []
                for u in users:
                    name = self.purple.PurpleConvChatCbGetName(u)
                    usernames.append(self.buddy_by_name(name))
                convs.append({"id":str(c),"participants":usernames})
        return "convs",{"convs":convs}
        self.debug_print("SENDING",args,im)
        self.purple.PurpleConvImSend(im, args['msg'])



def run_recv_thread(ctx):
    print("STARTING THREAD")
    DBusGMainLoop(set_as_default=True)
    bus = dbus.SessionBus()
    tx = ctx.socket(zmq.DEALER)
    tx.connect("inproc://purplerecv")

    obj = bus.get_object("im.pidgin.purple.PurpleService", "/im/pidgin/purple/PurpleObject")
    purple = dbus.Interface(obj, "im.pidgin.purple.PurpleInterface")
    #purple = bus.get("im.pidgin.purple.PurpleService", "/im/pidgin/purple/PurpleObject")
    
    def _recv(account, sender, message, conv, flags):
        print("maybe got one?",account,sender,message,conv,flags,tx)        
        print(account)
        buddy = purple.PurpleFindBuddy(account,sender)
        print("BUDDY",buddy)
        name = purple.PurpleBuddyGetName(buddy)
        alias = purple.PurpleBuddyGetAlias(buddy)
        print("BUDDY NAME",name,"ALIAS",alias)
        data = {"command":"recv", "msg":message, "from":"{} ({})".format(name,alias), "conv":str(conv), "account":account, "ischat":False}
        tx.send_string(json.dumps(data))
        print("sender: {} message: {}, account: {}, conversation: {}, flags: {}, conv: {}".format(sender,message,account,conv,flags,conv))

    def _recvchat(account, sender, message, conv, flags):
        print("maybe got one chat?",account,sender,message,conv,flags,tx)        
        print(account)
        buddy = purple.PurpleFindBuddy(account,sender)
        print("BUDDY",buddy)
        name = purple.PurpleBuddyGetName(buddy)
        alias = purple.PurpleBuddyGetAlias(buddy)
        print("BUDDY NAME",name,"ALIAS",alias)
        data = {"command":"recv", "msg":message, "from":"{} ({})".format(name,alias), "conv":str(conv), "account":account, "ischat":True}
        tx.send_string(json.dumps(data))
        print("sender: {} message: {}, account: {}, conversation: {}, flags: {}, conv: {}".format(sender,message,account,conv,flags,conv))

    def _sent(account, recipient, message):
        print(account)
        buddy = purple.PurpleFindBuddy(account,recipient)
        print("BUDDY",buddy)
        name = purple.PurpleBuddyGetName(buddy)
        alias = purple.PurpleBuddyGetAlias(buddy)
        print("BUDDY NAME",name,"ALIAS",alias)
        data = {"command":"send", "msg":message, "to":"{} ({})".format(name,alias), "account":account}
        tx.send_string(json.dumps(data))

    def _sentchat(account, message, chat_id):
        print(account)
        data = {"command":"send", "msg":message, "conv":str(chat_id), "account":account}
        tx.send_string(json.dumps(data))

    def _sign(buddy, online):
        print("BUDDY",buddy)
        name,alias = purple.PurpleBuddyGetName(buddy),purple.PurpleBuddyGetAlias(buddy)
        print("BUDDY NAME",name,"ALIAS",alias)
        data = {"command":"signon" if online else "signoff", "who":"{} ({})".format(name,alias)}
        tx.send_string(json.dumps(data))

    def _idle(buddy, oldidle, newidle):
        if oldidle == newidle:
            return
        print("BUDDY",buddy)
        name,alias = purple.PurpleBuddyGetName(buddy),purple.PurpleBuddyGetAlias(buddy)
        print("BUDDY NAME",name,"ALIAS",alias)
        data = {"command":"idle" if newidle else "unidle", "who":"{} ({})".format(name,alias)}
        tx.send_string(json.dumps(data))


    # OLD
    bus.add_signal_receiver(_recv, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="ReceivedImMsg")
    bus.add_signal_receiver(_recvchat, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="ReceivedChatMsg")
    
    bus.add_signal_receiver(_sent, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="SentImMsg")
    bus.add_signal_receiver(_sentchat, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="SentChatMsg")

    #bus.add_signal_receiver(_away, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="BuddyAway")
    bus.add_signal_receiver(_idle, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="BuddyIdle")
    bus.add_signal_receiver(lambda b: _sign(b, True), dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="BuddySignedOn")
    bus.add_signal_receiver(lambda b: _sign(b, False), dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="BuddySignedOff")

    # void (*buddy_away)(PurpleBuddy *buddy, PurpleStatus *old_status, PurpleStatus *status);
    # void (*buddy_idle)(PurpleBuddy *buddy, gboolean old_idle, gboolean idle);
    # void (*buddy_signed_off)(PurpleBuddy *buddy);
    # void (*buddy_signed_on)(PurpleBuddy *buddy);



    
    # NEW
    #purple.ReceivedImMsg.connect(_recv)
    
    
    loop = GLib.MainLoop()
    loop.run()


jam()
