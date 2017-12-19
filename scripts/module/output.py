from __future__ import print_function
from colorama import init
from colorama import Style

init()

def write(string="", color=Style.RESET_ALL):
    """write a string with color specified. No new line"""

    print(color, end="")
    print(string, end="")
    print(Style.RESET_ALL, end="")

def writeln(string="", color=Style.RESET_ALL):
    """write a string with color specified. New line"""

    print(color, end="")
    print(string, end="")
    print(Style.RESET_ALL)
