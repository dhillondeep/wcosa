"""@package module
Parent class parses the command line arguments and makes every base class extend handle_args
method to do the work
"""

import abc
from sys import platform


class Parent(object):
    """Parent class parses and has an abstract argument handler"""

    __metaclass__ = abc.ABCMeta

    def __init__(self):
        """Initialize the operating system"""

         # check operating system
        if platform == "linux" or platform == "linux2":
            # linux
            self.operating_system = "linux"
        elif platform == "darwin":
            # OS X
            self.operating_system = "mac"
        elif platform == "win32":
            # Windows
            self.operating_system = "windows"
        elif platform == "cygwin":
            self.operating_system = "cygwin"
        else:
            # Other
            self.operating_system = platform

    def parse(self):
        """Parse command line arguments"""

        from sys import argv
        args = {}  # Empty dictionary to store key-value pairs.
        while argv:  # While there are arguments left to parse...
            if argv[0][0] == '-':  # Found a "-name value" pair.
                if len(argv) == 1:
                    args[argv[0]] = ""
                else:
                    args[argv[0]] = argv[1]  # Add key and value to the dictionary.
            argv = argv[1:]  # Reduce the argument list by copying it starting from index 1.
        return args

    @abc.abstractmethod
    def handle_args(self, args):
        """Handle command line arguments"""
        return

    def start(self):
        """Entry Point of script which starts everything"""

        self.handle_args(self.parse())
