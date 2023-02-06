#! /usr/bin/env python3

import struct
import socket
import argparse
import ipaddress
import os
import signal
import sys
from smq.connect import Connect
from smq.subscribe import Subscribe
from smq.publish import Publish
from smq.unsubscribe import Unsubscribe
from socket_reader import SocketReader

HOST = "127.0.0.1"
PORT = 44567
quit = False

def main():
    parser = argparse.ArgumentParser(
        prog ="SMQ Client",
        description="Python client that communicates with SMQ server"
    )
    parser.add_argument(
        '-a', '--addr',
        type=ipv4_addr, default=HOST,
        help="IPv4 address of the SMQ server"
    )
    parser.add_argument(
        '-p', '--port',
        type=port, default=PORT,
        help="Remote port for the SMQ Server"
    )

    args = parser.parse_args()
    print("SMQ client V1.0.0")
    username = input("Username:> ")
    password = input("Password:> ")

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.connect((args.addr, args.port))
        connect = Connect(username, password)
        sock.sendall(connect.byte_string)
        response_code = read_response(sock)

        if response_code == 0xFF:
            print("General error during authentication with server")
            return
        elif response_code == 0x02:
            print("Invalid username/password")
            return
        else:
            print("Authentication successful")
        read_handlers = {
            2: publish_handler,
            6: puback_handler
        }
        socket_reader = SocketReader(sock, read_handlers)
        socket_reader.start()
        input_handler(sock)

    socket_reader.join()

def input_handler(sock):
    handlers = {
        "pub": publish_request,
        "sub": subscribe_handler,
        "unsub": unsubscribe_handler,
        "quit": disconnect_handler, 
    }
    while True:
        choice = input("SMQ:> ").split(" ")

        if choice[0] not in handlers:
            print(f"invalid choice: {choice[0]}")
            continue
        elif choice[0] == 'quit':
            disconnect_handler(sock)
            break
        else:
            handlers[choice[0]](sock, choice[1:])


def publish_request(sock, input):
    publish = Publish(input[0], input[1])
    sock.sendall(publish.byte_string)

def subscribe_handler(sock, input):
    subscribe = Subscribe(input[0])
    sock.sendall(subscribe.byte_string)

def unsubscribe_handler(sock, input):
    unsubscribe = Unsubscribe(input[0])
    sock.sendall(unsubscribe.byte_string)

def disconnect_handler(sock):
    print("Disconnecting...")
    os.kill(os.getpid(), signal.SIGUSR1)
    
def publish_handler(buffer):
    nread = 0
    size = struct.unpack("!H", buffer[:2])[0]
    nread += 2
    topic = struct.unpack(f"!{size}s", buffer[nread:size + 2])[0]
    nread += size
    size = struct.unpack("!H", buffer[nread:nread + 2])[0]
    nread += 2
    value = struct.unpack(f"!{size}s", buffer[nread:])[0]
    
    print(f"Topic: {topic} Message: {value}\nSMQ:> ", end="")

def puback_handler(buffer):
    response = buffer[0]

    if response == 1:
        print("Message was publishing successfully")
    else:
        print("Message was not published successfully")

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

def read_response(sock):
    buffer = sock.recv(5)
    packet_type, size = struct.unpack("!BI", buffer)
    buffer = sock.recv(size)
    response_code = struct.unpack("!B", buffer)

    if packet_type == 6:
        rtn = response_code[0]
    else:
        print(f"Unknown packet type {packet_type}")
        rtn = -1

    return rtn

def ipv4_addr(ip_string:str) -> ipaddress.IPv4Address:
    error_flag = False
    try:
        ip_addr = ipaddress.ip_address(ip_string)
    except ValueError:
        error_flag = True

    if error_flag or type(ip_addr) != ipaddress.IPv4Address:
        raise argparse.ArgumentTypeError(f"invalid IPv4 address {ip_string}")
    
    return ip_string

def port(port_str:str) -> int:
    port = int(port_str)

    if port <= 1024 or port >= 65535:
        raise argparse.ArgumentTypeError("port must be within range 1024 < x < 65535")

    return port

if __name__ == "__main__":
    main()