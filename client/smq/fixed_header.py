#!/usr/bin/env python3

import struct


class Fixed:
    def __init__(self, type, size):
        self._type = type
        self._size = size

    @property
    def byte_string(self):
        return struct.pack("!BI", self._type, self._size)

