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
    colour, colour_bg, colour_fg = '', '', ''
    tokens = []
    lines = msg.split('\n')
    for i in range(len(lines)):
        line = lines[i]
        text = ''
        # print(f'{line=}')
        for ch in line:
            if ch == '\033':
                is_colour = True
                if text:
                    tokens.append((colour_fg+colour_bg, text))
                    colour = ch
                    text = ''
                else:
                    colour += ch
            elif is_colour:
                colour += ch
                if ch == 'm':
                    is_colour = False
                    if '[38' in colour or '[39' in colour:
                        colour_fg = colour
                    elif '[48' in colour or '[49' in colour:
                        colour_bg = colour
            else:
                text += ch
        if colour:
            tokens.append((colour_fg+colour_bg, text))
            colour = ''
        if i < len(lines) - 1:
            if lines[i+1].strip().startswith('\x1b[49m'):
                tokens.append(('\x1b[49m', ''))  # reset background colour
            elif lines[i+1].strip().startswith('\x1b[0m'):
                tokens.append(('\x1b[0m', ''))
            elif lines[i+1].strip().startswith('\x1b[39m'):
                tokens.append(('\x1b[39m', ''))
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
        rev_line = [('', ' ' * (max_len - length))]
        for colour, text in reversed(line):
            rev = (colour, text[::-1])
            # print(f'{rev=}')
            rev_line += (colour, text[::-1])
        # yield ' '*4 + rev_line + str(length)
        if set(rev_line[-1][1]) == {' '}:
            rev_line[-1] = (rev_line[-1][0] + '\x1b[49m', rev_line[-1][1])

        yield ''.join(chain.from_iterable(rev_line))

def print_reversed_ansi(msg: str) -> None:
    # for line in tokenise_ansi(msg):
    #     print(f'{line=}')
    # print('original:')
    # print(msg)

    # print('scanned:')
    max_len = max_line_len(msg)
    lengths = []
    scanned = []
    for line in tokenise_ansi(msg):
        length = line_len(line)
        lengths.append(length)
        # print(''.join(chain.from_iterable(line)), end=' '*(max_len-length) + str(length) + '\n')
        # print(''.join(chain.from_iterable(line)), end='\x1b[0m\n')
        scanned.append(''.join(chain.from_iterable(line)) + ' ' * (max_len - length))
    print('\x1b[0m')

    # print('reversed:')
    rvrsed = []
    for rev_line in reverse_ansi(msg):
        # print(f'{rev_line=}')
        # print(rev_line, end='\x1b[0m\n')
        rvrsed.append(rev_line)
    print('\033[0m')

    print(msg)

    for i, (o, r, s) in enumerate(zip(msg.split('\n'), rvrsed, scanned)):
        print(f'{o+" "*(max_len-lengths[i])}\x1b[0m  |{s}\x1b[0m |  {r}')
            # print(f'\033[31m{s} | {r}\033[0m')


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
        print(sys.argv[1])
        print_reversed_ansi(open(sys.argv[1]).read())
