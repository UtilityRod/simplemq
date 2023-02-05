#! /usr/bin/env python3

import struct
import socket
import time
from smq.connect import Connect
from smq.subscribe import Subscribe
from smq.publish import Publish
from smq.unsubscribe import Unsubscribe

HOST = "127.0.0.1"
PORT = 44567

def main():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        connect = Connect("admin", "password")
        s.sendall(connect.byte_string)
        subscribe = Subscribe("test")
        s.sendall(subscribe.byte_string)
        publish = Publish("test", "value1")
        s.sendall(publish.byte_string)
        publish = Publish("test", "value2")
        s.sendall(publish.byte_string)
        unsub = Unsubscribe("test")
        s.sendall(unsub.byte_string)
        publish = Publish("test", "value3")
        s.sendall(publish.byte_string)

        for i in range(3):
            packet_type, size = read_fixed(s)
            if packet_type == 2:
                topic, value = read_publish(s, size)
                print(f"{topic=} {value=}")
            else:
                print("Unknown packet type")
        

def read_fixed(sock):
    buffer = sock.recv(5)
    return struct.unpack("!BI", buffer)

def read_publish(sock, size):
    buffer = sock.recv(size)
    nread = 0
    size = struct.unpack("!H", buffer[:2])[0]
    nread += 2
    topic = struct.unpack(f"!{size}s", buffer[nread:size + 2])[0]
    nread += size
    size = struct.unpack("!H", buffer[nread:nread + 2])[0]
    nread += 2
    value = struct.unpack(f"!{size}s", buffer[nread:])[0]
    return topic, value

if __name__ == "__main__":
    main()