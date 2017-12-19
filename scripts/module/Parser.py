"""@package module
Parser provides an abstarct parser which other classes extend to parse inputs
"""

import argparse
import abc

class Parser():
    """Abstarct parser used for parsing command line argumnets"""

    __metaclass__ = abc.ABCMeta

    def __init__(self, desc, name, max_pos):
        self.parser = argparse.ArgumentParser(description=desc,
        prog=name,
        formatter_class=lambda prog: argparse.HelpFormatter(prog,max_help_position=max_pos))

    @abc.abstractmethod
    def parse(self):
        """Abstract method to handle the cli arguments"""
