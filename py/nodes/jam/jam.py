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
        self.interface.add_interface({'send':self.send,'sendto':self.send_to})
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

    def purple_recv(self):
        rt,msg = self.purple_sock.recv_multipart()
        msg= json.loads(msg)
        self.debug_print("[JAM.PY] MSG",msg)
        cmd = msg["command"]
        del msg["command"]
        self.m_send(cmd, msg)
        
    def purple_err(self):
        self.debug_print("PRPL DIED")
                        
    def send(self, args):
        candidates = process.extract({"name":args['to']},self.buddies,processor=lambda b: b['name'],limit=3)
        candidates += process.extract({"alias":args['to']},self.buddies,processor=lambda b: b['alias'],limit=3)
        candidates = sorted([c for c in candidates if c[1] > self.match_quality_threshold],key=lambda c: c[1])
        self.debug_print("Filtered candidates:",candidates)
        if len(candidates) == 1 or (len(candidates) == 2 and candidates[0][0] == candidates[1][0]):
            buddy = candidates[0][0]
        elif len(candidates) == 0:
            self.debug_print("Target not found: {}".format(args['to']))
            return
        else:
            self.debug_print("Ambiguous target: {}\nSome possible matches:\n{}".format(args['to'],"\n".join(["- Id: {}, Alias: {}, Quality: {}".format(c[0]['name'], c[0]['alias'], c[1]) for c in candidates])))
            return
        
        conv = self.purple.PurpleConversationNew(1, int(buddy['account_id']), str(buddy['name']))
        im = self.purple.PurpleConvIm(conv)
        self.debug_print("SENDING",args,buddy,conv,im)
        self.purple.PurpleConvImSend(im, args['msg'])

    def send_to(self, args):
        im = self.purple.PurpleConvIm(int(args['conv']))
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
        data = {"command":"recv", "msg":message, "from":"{} ({})".format(name,alias), "conv":str(conv), "account":account}
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
    bus.add_signal_receiver(_recv, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="ReceivedChatMsg")
    
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
