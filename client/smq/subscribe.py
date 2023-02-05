#! /usr/bin/env python3

import struct
from .fixed_header import Fixed

class Subscribe:
    def __init__(self, topic):
        self._topic = bytes(topic, 'utf-8')

    @property
    def byte_string(self):
        subscribe =  struct.pack(
            f"!H{len(self._topic)}s",
            len(self._topic), self._topic
        )

        fixed = Fixed(3, len(subscribe)).byte_string
        return fixed + subscribe