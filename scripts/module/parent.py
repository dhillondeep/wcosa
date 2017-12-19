"""@package module
Parent provides similar functionality that each component of the tool will use.
This includes paths and information about operating system 
"""

import argparse
import os
import abc
from sys import platform
import module.helper as helper


class Parent(object):
    """Parent class is an abstract class for all the build, create and serial moduels"""

    __metaclass__ = abc.ABCMeta

    def __init__(self, path, board, ide):
        """Initialize the paths, operating system and the parser"""

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

        if path is None:
            self.curr_path = helper.linux_path(os.getcwd(), self.operating_system)
        else:
            self.curr_path = path

        self.dir_name = os.path.basename(self.curr_path)

        self.wcosa_path = helper.linux_path(os.path.abspath(os.path.dirname(
            os.path.abspath(__file__)) + "/../.."), self.operating_system)
        self.cmake_templates_path = self.wcosa_path + "/build/cmake-files"
        self.config_files_path = self.wcosa_path + "/build/config-files"

        if ide is None:
            self.ide = "None"
        else:
            self.ide = ide

        self.board = board
