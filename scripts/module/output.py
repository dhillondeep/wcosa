from __future__ import print_function
from colorama import init
from colorama import Fore, Back, Style

init()

def write(string="", color=Style.RESET_ALL):
    print(color, end="")
    print(string, end="")
    print(Style.RESET_ALL, end="")

def writeln(string="", color=Style.RESET_ALL):
    print(color, end="")
    print(string, end="")
    print(Style.RESET_ALL)
