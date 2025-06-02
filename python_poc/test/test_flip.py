from .. import flip

import sys
from itertools import chain

class TestTokenise:

    def test_tokenise_line_continuation(self):
        result = []
        for line_tokens in flip.tokenise_ansi(open('test/data/1_2.cow').read()):
            result.append(''.join(chain.from_iterable(line_tokens)))

        expected = [
            '''    \x1b[49m   \x1b[38;5;16m▄▄\x1b[48;5;16m\x1b[38;5;142m▄▄▄\x1b[49m\x1b[38;5;16m▄▄''',
            '''     ▄\x1b[48;5;16m\x1b[38;5;58m▄\x1b[48;5;58m\x1b[38;5;70m▄\x1b[48;5;70m \x1b[48;5;227m    \x1b[48;5;237m\x1b[38;5;227m▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[49m\x1b[38;5;16m▄''',
        ]


        sys.stdout.write('\n'.join(['expected:', *expected]) + '\n')
        sys.stdout.write('\n'.join(['result:', *result]) + '\n')

        assert result == expected

    def test_tokenise_umbreon(self):
        result = []
        for line_tokens in flip.tokenise_ansi(open('test/data/umbreon.cow').read()):
            result.append(''.join(chain.from_iterable(line_tokens)))

        expected = [
            '''    \x1b[49m   \x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;237m▄▄ \x1b[49m''',
            '''      \x1b[48;5;16m \x1b[48;5;214m▄ \x1b[48;5;237m\x1b[38;5;214m▄\x1b[48;5;16m \x1b[49m \x1b[38;5;16m▄▄\x1b[48;5;16m\x1b[38;5;237m▄▄▄ \x1b[49m   \x1b[38;5;16m▄▄▄''',
            '''     \x1b[48;5;16m \x1b[48;5;237m \x1b[38;5;232m▄▄\x1b[48;5;16m \x1b[49m\x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;235m▄\x1b[48;5;94m\x1b[38;5;237m▄\x1b[48;5;214m \x1b[48;5;237m\x1b[38;5;214m▄\x1b[38;5;235m▄\x1b[48;5;232m\x1b[38;5;16m▄\x1b[49m▀ ▄\x1b[48;5;16m\x1b[38;5;214m▄\x1b[48;5;235m\x1b[38;5;94m▄ \x1b[48;5;16m \x1b[49m''',
            '''     \x1b[48;5;16m \x1b[48;5;232m\x1b[38;5;214m▄\x1b[48;5;214m\x1b[38;5;237m▄\x1b[48;5;237m\x1b[38;5;214m▄  \x1b[48;5;232m\x1b[38;5;237m▄\x1b[48;5;237m\x1b[38;5;235m▄\x1b[48;5;235m \x1b[48;5;94m\x1b[38;5;16m▄\x1b[49m▀  \x1b[48;5;16m \x1b[48;5;94m\x1b[38;5;235m▄\x1b[48;5;214m▄ \x1b[48;5;94m\x1b[38;5;16m▄\x1b[49m▀''',
            '''    \x1b[48;5;16m \x1b[48;5;231m\x1b[38;5;52m▄\x1b[48;5;214m\x1b[38;5;237m▄\x1b[48;5;237m\x1b[38;5;214m▄\x1b[48;5;214m\x1b[38;5;237m▄\x1b[48;5;237m\x1b[38;5;235m▄\x1b[38;5;231m▄\x1b[38;5;16m▄\x1b[48;5;235m\x1b[38;5;237m▄\x1b[48;5;16m \x1b[49m\x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;235m▄▄\x1b[38;5;232m▄\x1b[48;5;235m \x1b[38;5;16m▄▄\x1b[49m▀''',
            '''    \x1b[48;5;16m \x1b[48;5;232m\x1b[38;5;237m▄\x1b[48;5;237m  \x1b[38;5;235m▄\x1b[48;5;16m\x1b[38;5;52m▄\x1b[48;5;196m\x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;237m▄\x1b[48;5;237m\x1b[38;5;232m▄\x1b[48;5;232m \x1b[38;5;235m▄\x1b[48;5;235m  \x1b[38;5;214m▄\x1b[48;5;16m\x1b[38;5;235m▄\x1b[49m\x1b[38;5;16m▄''',
            '''     ▀\x1b[48;5;239m▄\x1b[48;5;237m▄\x1b[38;5;232m▄▄\x1b[48;5;235m \x1b[48;5;232m\x1b[38;5;235m▄\x1b[48;5;235m    \x1b[48;5;214m \x1b[48;5;235m \x1b[48;5;214m \x1b[48;5;16m\x1b[38;5;232m▄\x1b[49m\x1b[38;5;16m▄''',
            '''       \x1b[48;5;16m \x1b[48;5;235m  \x1b[38;5;214m▄\x1b[48;5;214m\x1b[38;5;237m▄\x1b[48;5;235m\x1b[38;5;94m▄\x1b[38;5;16m▄\x1b[48;5;232m▄\x1b[49m▀▀\x1b[48;5;94m▄\x1b[48;5;232m\x1b[38;5;235m▄▄\x1b[48;5;16m \x1b[49m''',
            '''      \x1b[38;5;16m▄\x1b[48;5;16m\x1b[38;5;235m▄\x1b[48;5;235m \x1b[48;5;16m \x1b[48;5;214m\x1b[38;5;237m▄\x1b[48;5;237m\x1b[38;5;94m▄\x1b[48;5;16m \x1b[49m     \x1b[38;5;16m▀▀''',
            '''      ▀\x1b[48;5;235m▄\x1b[48;5;16m \x1b[48;5;237m  \x1b[48;5;16m \x1b[49m''',
            '''         ▀▀\x1b[39m\x1b[39m''',
            '',
        ]

        sys.stdout.write('\n'.join(['expected:', *expected]) + '\n')
        sys.stdout.write('\n'.join(['result:', *result]) + '\n')

        assert result == expected
