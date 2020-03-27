#!/usr/bin/python2.7

# ==============================================
# otto -- a dbus-based auto-responding plugin. 
# Helper utility: /suit/mini/goaway
# ==============================================

#import park
import dbus, gobject, re, os, sys, random, time
import dbus.service
from dbus.mainloop.glib import DBusGMainLoop
from threading import Timer

PURPLE_CONV_TYPE_IM = 1
PURPLE_CONV_TYPE_CHAT = 2

class otto(dbus.service.Object):
    def __init__(self):
        self.logfile = "/home/zoom/suit/otto/otto.log"
        self.help_text = self.read_whole_file("/home/zoom/suit/otto/help")
        self.convs = {}
        self.away_msg = ""
        self.timeout = 60*15
        self.idle_checked = 0
        self.idle_update_freq = 10
        self.idle_time = 0
        self.idle_cutoff = 15*60
        self.update_away()
        self.update_idle()
        #self.park = park.park()
        self.jokes = []
        self.blacklist = ["ak41s"]
        self.listeners = {}
        self.one_shot_listeners = {}
        self.timer_time = -1
        self.timer_start_time = -1
        f = open("/home/zoom/suit/otto/jokes","r")
        for l in f:
            self.jokes += [l[:-1].split("...")]

        abus = dbus.SessionBus()
        abus.add_signal_receiver(self.report, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="ReceivedImMsg")
        self.obj = abus.get_object("im.pidgin.purple.PurpleService", "/im/pidgin/purple/PurpleObject")
        self.purple = dbus.Interface(self.obj, "im.pidgin.purple.PurpleInterface")

        bus_name = dbus.service.BusName('Ziggy.otto', bus=dbus.SessionBus())
        dbus.service.Object.__init__(self, bus_name, '/Ziggy/otto')

    @dbus.service.method('Ziggy.otto')
    def status(self):
        return "Listeners: " + str(self.listeners) + "\nGag list: " + str(self.blacklist) + "\nAway: " + self.get_away_msg()

    @dbus.service.method('Ziggy.otto')
    def notify(self, notification):
        for s in self.listeners.values() + self.one_shot_listeners.values():
            self.purple.PurpleConvImSend(s,"[otto.notify]: " + notification)
        self.one_shot_listeners = {}
        return "Notified: " + ",".join(self.listeners.keys() + self.one_shot_listeners.keys())

    @dbus.service.method('Ziggy.otto')
    def goaway(self, s):
        self.convs = {}
        self.away_msg = s
        if(s == ""): 
            return "Not away"
        return "Away: " + self.away_msg

    @dbus.service.method('Ziggy.otto')
    def goawait(self, s, t):
        self.convs = {}
        self.away_msg = s
        if(s == ""): 
            return "Not away"
        else:
            self.timer_start_time = time.time()
            self.timer_time = t
            return "Away: " + self.away_msg + " for " + str(t) + " seconds" 

    def time_remaining(self):
        print("TIME")
        if(self.timer_time == -1):
            return 0
        else:
            now = time.time()
            if(now - self.timer_start_time >= self.timer_time):
                self.timer_start_time = -1
                self.timer_time = -1
                self.away_msg = ""
                return 0
            else:
                return self.timer_start_time + self.timer_time - now

    def get_away_msg(self):
        secs_left = self.time_remaining()
        msg = self.away_msg
        if(secs_left > 0):
            if(secs_left < 60):
                msg += " (%d seconds remaining)"%(secs_left)
            elif(secs_left < 3600):
                msg += " (%d minutes remaining)"%(secs_left//60)
            else:
                msg += " (%dh:%dm remaining)"%(secs_left//3600, (secs_left%3600)//60)
        return msg
    
    def read_whole_file(self,fn):
        with open(fn) as f: return f.read()

    def write_to_logfile(self,s):
        with open(self.logfile,'a') as f: f.write(time.strftime("%Y%m%d:%H%M%S %Z",time.localtime()) + ": " + s + "\n")

    def update_away(self):
        self.away_msg = self.read_whole_file("/home/zoom/AWAY").strip()

    def update_idle(self):
        os.system("echo $(($(xssstate -i)/1000)) > /home/zoom/IDLE")
        self.idle_time = int(self.read_whole_file("/home/zoom/IDLE"))

    def report(self,account, sender, message, conv, flags):
        msg = re.sub("<[^>]*>","",str(message)).strip()
        target = self.purple.PurpleConvIm(conv)
        if(int(conv) == 0):
            target = self.purple.PurpleConversationNew(PURPLE_CONV_TYPE_IM,account,str(sender))
            print("AAAAA: " + str(target));
            target = self.purple.PurpleConvIm(target)
            print("BBBBB: " + str(target));
        print("Hello: "+str(sender)+" -- "+str(message))
        current_away_msg = self.get_away_msg()
        if(current_away_msg != ""):
            now = time.time()
            s = str(sender)
            print([b for b in self.blacklist if b in s])
            if((not (s in self.convs.keys()) or (now - self.convs[s] >= self.timeout)) and len([b for b in self.blacklist if b in s]) == 0):
                self.convs[s] = now
                self.write_to_logfile("Away: " + str(sender))
                self.purple.PurpleConvImSend(target, "[otto.away]: " + current_away_msg)
        if(msg == "!help"):
            self.purple.PurpleConvImSend(self.purple.PurpleConvIm(conv),"[otto.help]: " + self.help_text)
#        elif(msg[:3] == "!p "):
#            self.write_to_logfile("Park: " + str(sender))
#            resp = self.park.handle(str(sender),msg[3:])
#            self.purple.PurpleConvImSend(target,"[otto.park]: " + resp)
        elif(msg == "!gag"):
            self.write_to_logfile("Gag: " + str(sender))
            if not str(sender) in self.blacklist:
                self.blacklist += [str(sender)]
            print(self.blacklist)
            self.purple.PurpleConvImSend(target,"[otto.gag]: Added to gag list")
        elif(msg == "!ungag"):
            self.write_to_logfile("Ungag: " + str(sender))
            if str(sender) in self.blacklist:
                self.blacklist.remove(str(sender))
            self.purple.PurpleConvImSend(target,"[otto.ungag]: Not in gag list")
        elif(msg == "!listen"):
            self.listeners[str(sender)] = target
            self.write_to_logfile("Listening: " + str(sender))
            self.purple.PurpleConvImSend(target,"[otto.listen]: Added to listeners")
        elif(msg == "!listen1"):
            self.one_shot_listeners[str(sender)] = target
            self.write_to_logfile("Listening1: " + str(sender))
            self.purple.PurpleConvImSend(target,"[otto.listen]: Added to one-shot listeners")
        elif(msg == "!unlisten"):
            self.write_to_logfile("Unlisten: " + str(sender))
            if str(sender) in self.listeners.keys():
                del self.listeners[str(sender)]
            self.purple.PurpleConvImSend(target,"[otto.listen]: Not listening")
        elif(msg == "!away"):
            self.write_to_logfile("Away ping: " + str(sender))
            self.purple.PurpleConvImSend(target,"[otto.away]: " + ("[None]" if (self.away_msg == "") else self.get_away_msg()))
        elif(msg == "!joke"):
            self.write_to_logfile("Joke ping: " + str(sender))
            j = random.choice(self.jokes)
            print(j)
            i = 0
            l = len(j)
            for line in j:
                self.purple.PurpleConvImSend(target,"[otto.joke]: " + line)
                i += 1
                print(i,l)
                if i < l and l > 1:
                    print('asleep')
                    time.sleep(2)
        elif(msg == "!idle"):
            self.write_to_logfile("Idle ping: " + str(sender))
            self.update_idle()
            #now = time()
            # if((now - self.idle_checked) >= self.idle_update_freq):
            #     self.update_idle()
            #     self.idle_checked = now
            if(self.idle_time >= self.idle_cutoff):
                self.purple.PurpleConvImSend(target,"[otto.idle]: Idling for over 15 minutes")
            else:
                self.purple.PurpleConvImSend(target,"[otto.idle]: Idling for %d:%02d"%(int(self.idle_time/60),self.idle_time%60))
                # sanitised = 15*int(self.idle_time/15)
                # if(sanitised != 0): 
                #     self.purple.PurpleConvImSend(self.purple.PurpleConvIm(conv),"[otto.idle]: Idling for at least " + str(sanitised) + " minutes")
                # else:
                #     self.purple.PurpleConvImSend(self.purple.PurpleConvIm(conv),"[otto.idle]: Idling for under 15 minutes")
	

#def cc(conv):
#    away_msg = open("/home/zoom/AWAY").read()
#    if(away_msg != ""):
#        purple.PurpleConvImSend(purple.PurpleConvIm(conv),"[otto.away]: " + away_msg)

    #if(message == "!away" or message == "yt?" or message == "allos" or "you there?" in message or "if you're there" in message or re.match(".*you [^.A-Z]*there\?", message) != None):
    #    if(away_msg != ""):
    #        purple.PurpleConvImSend(purple.PurpleConvIm(conv),"[otto.away]: " + away_msg)
    #elif(message == "!idle"):
    #    purple.PurpleIdleGetUIOps().getTimeIdle()
    #    purple.PurpleConvImSend(purple.PurpleConvIm(conv),"[otto.idle]: Error -- Idle time not successfully retrieved")
    #elif(message[0] == '!help'):
    #    purple.PurpleConvImSend(purple.PurpleConvIm(conv),"[otto.help]: !away -- reports away message, if any; !idle -- reports (fuzzy) idle time")
            
#away_msg = read_away()

#bus.add_signal_receiver(cc, dbus_interface="im.pidgin.purple.PurpleInterface", signal_name="ConversationCreated")


dbus.mainloop.glib.DBusGMainLoop(set_as_default=True)
o = otto()
myloop = gobject.MainLoop()
myloop.run()
