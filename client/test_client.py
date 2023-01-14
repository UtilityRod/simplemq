#! /usr/bin/env python3

import struct
import socket
import time

HOST = "127.0.0.1"
PORT = 44567

def main():
    proto_name = b"SMQ"
    proto_version = 1
    client_name = b"qwerty123"
    username = b"user"
    password = b"pass"

    connect_buffer = struct.pack(
        f"!H{len(proto_name)}sBH{len(client_name)}sH{len(username)}sH{len(password)}s",
        len(proto_name), proto_name,
        proto_version,
        len(client_name), client_name,
        len(username), username,
        len(password), password
        )

    connect_fixed = struct.pack("!BI", 1, len(connect_buffer))

    topic = b"test"
    value = b"value"
    publish_buffer = struct.pack(
        f"!H{len(topic)}sH{len(value)}s",
        len(topic), topic,
        len(value), value
    )

    publish_fixed = struct.pack("!BI", 2, len(publish_buffer))

    subscribe_buffer = struct.pack(
        f"!H{len(topic)}s",
        len(topic), topic
    )

    subscribe_fixed = struct.pack("!BI", 3, len(subscribe_buffer))

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        s.sendall(connect_fixed + connect_buffer)
        s.sendall(subscribe_fixed + subscribe_buffer)
        s.sendall(publish_fixed + publish_buffer)
        buffer = s.recv(5)
        packet, remaining = struct.unpack("!BI", buffer)
        buffer = s.recv(remaining)

        nread = 0
        size = struct.unpack("!H", buffer[:2])[0]
        nread += 2
        topic = struct.unpack(f"!{size}s", buffer[nread:size + 2])[0]
        nread += size
        size = struct.unpack("!H", buffer[nread:nread + 2])[0]
        nread += 2
        value = struct.unpack(f"!{size}s", buffer[nread:])[0]
        print(packet, topic, value)

if __name__ == "__main__":
    main()