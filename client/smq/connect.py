#! /usr/bin/env python3

import random
import string
import struct
from .fixed_header import Fixed

class Connect:
    def __init__(self, username, password):
        self._proto_name = b"SMQ"
        self._proto_version = 1
        self._client_name = bytes("".join(
            random.choices(
                string.ascii_uppercase +
                string.ascii_lowercase +
                string.digits,
                k=10
            )
        ), 'utf-8')
        self._username = bytes(username, 'utf-8')
        self._password = bytes(password, 'utf-8')

    @property
    def byte_string(self):
        connect = struct.pack(
            f"!H{len(self._proto_name)}sBH{len(self._client_name)}sH{len(self._username)}sH{len(self._password)}s",
            len(self._proto_name), self._proto_name,
            self._proto_version,
            len(self._client_name), self._client_name,
            len(self._username), self._username,
            len(self._password), self._password
            )

        fixed = Fixed(1, len(connect)).byte_string
        return fixed + connect