#! /usr/bin/evn python3

import struct
from .fixed_header import Fixed

class Publish:
    def __init__(self, topic, msg):
        self._topic = bytes(topic, 'utf-8')
        self._msg = bytes(msg, 'utf-8')

    @property
    def byte_string(self):
        publish = struct.pack(
            f"!H{len(self._topic)}sH{len(self._msg)}s",
            len(self._topic), self._topic,
            len(self._msg), self._msg
        )

        fixed = Fixed(2, len(publish)).byte_string
        return fixed + publish

