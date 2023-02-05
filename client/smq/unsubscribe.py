#! /usr/bin/env python3

import struct
from .fixed_header import Fixed


class Unsubscribe:
    def __init__(self, topic):
        self._topic = bytes(topic, 'utf-8')

    @property
    def byte_string(self):
        unsubscribe = struct.pack(
            f"!H{len(self._topic)}s",
            len(self._topic), self._topic
        )

        fixed = Fixed(4, len(unsubscribe)).byte_string
        return fixed + unsubscribe