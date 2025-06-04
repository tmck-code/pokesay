#!/usr/bin/env python3

import os, sys
from typing import Iterator
from itertools import chain
from collections import deque

DEBUG = os.environ.get('DEBUG', '0') == '1'

def tokenise_ansi(msg: str) -> Iterator[list[tuple[str, str]]]:
    colour, colour_fg, colour_bg = '', deque([], 1), deque([], 1)
    is_colour, is_reeet = False, False
    for i, line in enumerate(msg.split('\n')):
        colour, text, tokens = '', '', []
        for ch in line:
            # start of colour sequence detected!
            if ch == '\033':
                is_colour = True
                # if there is text in the current token buffer,
                if text:
                    # save it,
                    if DEBUG:
                        print(f'> {text=}, {colour_fg=}, {colour_bg=}, {colour=}, {is_colour=}')
                    tokens.append((''.join([*colour_fg, *colour_bg]), text))
                    colour = ch # then start a new colour token (unknown if bg/fg)
                    text = ''   # & clear the text buffer

                    # Do NOT clear the colour_fg/bg variables!
                    # This allows the next line to continue with the same colour.

                # or, if the text buffer is empty,
                else:
                    colour = ch # keep appending to the current colour token
            # middle of colour sequence
            elif is_colour:
                # append to current colour token
                colour += ch
                # check if we're at the end of a colour sequence
                if ch == 'm':
                    # if so, set the colour detection to false
                    is_colour = False
                    # check if the current colour token is fg or bg
                    if '[38' in colour or '[39' in colour:
                        colour_fg.append(colour)
                    elif '[48' in colour or '[49' in colour:
                        colour_bg.append(colour)
                    elif '[0m' in colour:
                        # if we reach the end of a colour sequence, reset colours
                        is_reset = True
                        colour = ''
            # othwise, if we're not in a colour sequence
            else:
                text += ch # append to the current text buffer
        # if we reach the end of a line, we need to flush the text buffer
        if colour:
            if DEBUG:
                print(f'! {text=}, {colour_fg=}, {colour_bg=}, {colour=}, {is_colour=}\n')
            tokens.append((''.join([*colour_fg, *colour_bg]), text))
        yield tokens
        tokens = []

def max_line_len(msg: str) -> int:
    max_len = 0
    for line in tokenise_ansi(msg):
        line_len = sum(len(text) for _, text in line)
        if line_len > max_len:
            max_len = line_len
    return max_len

def line_len(line: list[tuple[str, str]]) -> int:
    return sum(len(text) for _, text in line)

def reverse_ansi(msg: str) -> Iterator[str]:
    max_len = max_line_len(msg)

    for line in tokenise_ansi(msg):
        yield (
            [('', ' ' * (max_len - line_len(line)))] +
            list((colour, text[::-1]) for colour, text in reversed(line))
        )

def print_reversed_ansi(msg: str) -> None:
    print('original:', msg, '\x1b[0m', sep='\n')

    print('scanned:')
    for line in tokenise_ansi(msg):
        print(''.join(chain.from_iterable(line)), end='\x1b[0m\n')

    print('reversed:')
    for rev_line in reverse_ansi(msg):

        line = ' '*4 + ''.join(chain.from_iterable(rev_line))
        print(line, end='\x1b[0m\n')

if __name__ == '__main__':
    if len(sys.argv) > 1:
        print(sys.argv[1])
        print_reversed_ansi(open(sys.argv[1]).read())
