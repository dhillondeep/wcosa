"""
Handle Serial related tasks
"""

from module.parent import Parent


def open_serial(param):
    pass


def open_plotter(param):
    pass


class Serial(Parent):
    """Serial is used to show serial and plotter for WCosa projects"""

    def handle_args(self, args):
        """Allocates tasks for opening serial and plotter based on the args received"""

        if args[0] == "create":
            open_serial(args[1])
        elif args[0] == "update":
            open_plotter(args[1])


if __name__ == '__main__':
    serial = Serial()
    serial.start()
