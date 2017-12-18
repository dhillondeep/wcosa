"""
Handle regular in project usage of Wcosa
"""

from module.parent import Parent


def build():
    pass


def clean():
    pass


def upload(port):
    pass


def boards():
    pass


def config(board):
    pass


class Use(Parent):
    """Use is used to build, clean and upload WCosa projects"""

    def handle_args(self, args):
        """Allocates tasks for building, cleaning and uploading based on the args received"""

        if args[0] == "build":
            build()
        elif args[0] == "clean":
            clean()
        elif args[0] == "upload":
            upload(args[1])
        elif args[0] == "boards":
            boards()
        elif args[0] == "config":
            config(args[1])


if __name__ == '__main__':
    use = Use()
    use.start()
