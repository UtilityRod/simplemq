#! /usr/bin/env python3

import struct
import socket

HOST = "127.0.0.1"
PORT = 44567

def main():
    proto_name = b"SMQ"
    proto_version = 1
    client_name = b"qwerty123"
    username = b"user"
    password = b"pass"

    buffer = struct.pack(
        f"!H{len(proto_name)}sBH{len(client_name)}sH{len(username)}sH{len(password)}s",
        len(proto_name), proto_name,
        proto_version,
        len(client_name), client_name,
        len(username), username,
        len(password), password
        )

    fixed = struct.pack("!BI", 1, len(buffer))

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))
        s.sendall(fixed + buffer)

if __name__ == "__main__":
    main()