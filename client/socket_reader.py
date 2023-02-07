#!/usr/bin/env python3

import threading
import struct
import signal
import socket
from selectors import DefaultSelector, EVENT_READ

class SocketReader(threading.Thread):
    def __init__(self, sock, handlers):
        threading.Thread.__init__(self)
        self._sock = sock
        self._handlers = handlers
        self._private_read, self._private_write = socket.socketpair()
        signal.signal(signal.SIGUSR1, self._die)

    def run(self):
        sel = DefaultSelector()
        sel.register(self._private_read, EVENT_READ)
        sel.register(self._sock, EVENT_READ)
        quit_flag = False

        while not quit_flag:
            for key, _ in sel.select():

                if key.fileobj == self._sock:

                    try:
                        buffer = self._sock.recv(5)
                        packet_type, size = struct.unpack("!BI", buffer)
                        buffer = self._sock.recv(size)
                    except OSError:
                        continue 
                    self._handlers[packet_type](buffer)
                else:
                    quit_flag = True
                    break

    def _die(self, signum, frame):
        self._private_write.sendall(b'\0')
        
