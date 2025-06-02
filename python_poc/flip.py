#!/usr/bin/env python3

import sys
from typing import Iterator
from itertools import chain

msg = [
    # red ▄, green ▄
    "\033[38;5;160m▄x \033[38;5;46m▄",
    # green (via continuation)▄, yellow ▄
    "x▄ \033[38;5;190m▄",
]

def tokenise_ansi(msg: str) -> Iterator[list[tuple[str, str]]]:
    is_colour = False
    colour = ''
    tokens = []
    for line in msg.split('\n'):
        text = ''
        # print(f'{line=}')
        for ch in line:
            if ch == '\033':
                is_colour = True
                if text:
                    tokens.append((colour, text))
                    colour = ch
                    text = ''
                else:
                    colour += ch
            elif is_colour:
                colour += ch
                if ch == 'm':
                    is_colour = False
            else:
                text += ch
        if colour:
            tokens.append((colour, text))
            colour = ''
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
        length = line_len(line)
        rev_line = ' ' * (max_len - length)
        for colour, text in reversed(line[1:]):
            rev = [colour, text[::-1]]
            # print(f'{rev=}')
            rev_line += ''.join(rev)
        yield ' '*3 + rev_line

def print_reversed_ansi(msg: str) -> None:
    for line in tokenise_ansi(msg):
        print(f'{line=}')
    print('original:')
    print(msg)

    print('scanned:')
    for line in tokenise_ansi(msg):
        print(''.join(chain.from_iterable(line)), end='\n')

    print('reversed:')
    for rev_line in reverse_ansi(msg):
        # print(f'{rev_line=}')
        print(rev_line, end='\n')
    print('\033[0m')


# t = list(tokenise_ansi('\n'.join(msg)))
# print('\n'.join(msg))
# print('\033[0m')
#
# for line in t:
#     print(f'\033[0m{line=}')
# print()
#
# print_reversed_ansi('\n'.join(msg))

if __name__ == '__main__':
    if len(sys.argv) > 1:
        print_reversed_ansi(open(sys.argv[1]).read())
