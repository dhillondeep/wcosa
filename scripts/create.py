"""
Handle creating and updating WCosa projects
"""
import os
import json
from shutil import copyfile
from colorama import Fore
from module.parent import Parent
from module.output import write, writeln


class Create(Parent):
    """Create is used to create and update WCosa projects"""

    def __init__(self):
        """Initialize paths to use for creation"""

        super(Create, self).__init__()
        self.curr_path = self.linux_path(os.getcwd())
        self.dir_name = os.path.basename(self.curr_path)
        self.wcosa_path = self.linux_path(os.path.abspath(os.path.dirname(
            os.path.abspath(__file__)) + "/.."))


    def create_internal_config(self):
        """Creates and fills the internal config file so that it can be used"""

        copyfile(self.wcosa_path + "/scripts/required-files/internal-config.json",
                 self.curr_path + "/wcosa/internal-config.json")

        config_file = open(self.curr_path + "/wcosa/internal-config.json")
        config_data = json.load(config_file)
        config_data["project-name"][1] = self.dir_name
        config_data["build-type"][1] = "RELEASE"
        config_data["cosa-path"][1] = self.wcosa_path + "/platform/cosa/cores/cosa"
        config_data["cosa-lib-path"][1] = self.wcosa_path + "/platform/cosa/libraries"
        config_data["generic-gcc-avr-cmake-path"][1] = self.curr_path + "/wcosa/generic-gcc-avr.cmake"
        config_data["avr-mcu"][1] = "mcu"
        config_data["avr-h-fuse"][1] = "h-fuse"
        config_data["avr-l-fuse"][1] = "l-fuse"
        config_data["mcu-speed"][1] = "140000000L"
        config_data["cmake-c-flags-release"][1] = "None"
        config_data["cmake-cxx-flags-release"][1] = "None"
        config_data["cmake-c-flags-debug"][1] = "None"
        config_data["cmake-cxx-flags-debug"][1] = "None"
        config_file.close()
        config_file = open(self.curr_path + "/wcosa/internal-config.json", "w")
        json.dump(config_data, config_file)
        config_file.close()

        return config_data


    def fill_cmake_template(self, config_data):
        """Fills the root level CMakeLists file templates"""

        cmake_file = open(self.curr_path + "/CMakeLists.txt")
        cmake_str = cmake_file.read()
        cmake_file.close()

        for key in config_data:
            # ignore the comment
            if key == "__comment__":
                continue

            whether_quotation = config_data[key][0]
            value = config_data[key][1]

            if whether_quotation is True:
                adder = "\""
            else:
                adder = ""

            cmake_str = cmake_str.replace("%" + key, adder + value + adder)

        cmake_file = open(self.curr_path + "/CMakeLists.txt", "w")
        cmake_file.write(cmake_str)
        cmake_file.close()


    def create_cosa(self, board):
        """Creates WCosa project"""

        write("Creating work environment - ", color=Fore.CYAN)

        # create src folder
        if not os.path.exists(self.curr_path + "/src"):
            os.makedirs(self.curr_path + "/src")

        # create lib folder
        if not os.path.exists(self.curr_path + "/lib"):
            os.makedirs(self.curr_path + "/lib")

        # create bin folder
        if not os.path.exists(self.curr_path + "/bin"):
            os.makedirs(self.curr_path + "/bin")
        
        # create wcosa folder
        if not os.path.exists(self.curr_path + "/wcosa"):
            os.makedirs(self.curr_path + "/wcosa")

        # copy config file to main dir
        copyfile(self.wcosa_path + "/scripts/required-files/config.json",
                 self.curr_path + "/config.json")

        # copy cmake files
        copyfile(self.wcosa_path + "/scripts/required-files/CMakeLists-lib.txt",
                 self.curr_path + "/lib/CMakeLists.txt")

        copyfile(self.wcosa_path + "/scripts/required-files/CMakeLists-src.txt",
                 self.curr_path + "/src/CMakeLists.txt")

        copyfile(self.wcosa_path + "/scripts/required-files/CMakeLists-main.txt",
                 self.curr_path + "/CMakeLists.txt")

        copyfile(self.wcosa_path + "/scripts/required-files/generic-gcc-avr.cmake",
                 self.curr_path + "/wcosa/generic-gcc-avr.cmake")

        writeln("done")

        write("Creating project configuration - ", Fore.CYAN)

        # load config file and write changes
        config_file = open(self.curr_path + "/config.json")
        config_data = json.load(config_file)
        config_data["board"] = board
        config_file.close()
        config_file = open(self.curr_path + "/config.json", "w")
        json.dump(config_data, config_file)
        config_file.close()

        # create internal config file and fill it
        config_data = self.create_internal_config()

        # fill in CMakeLists templates
        self.fill_cmake_template(config_data)

        writeln("done")

        writeln("Finished Creation: ", Fore.CYAN)
        writeln("src        -> Source files", Fore.CYAN)
        writeln("lib        -> Library files (each library in seperate folder)", Fore.CYAN)
        writeln("bin        -> Binary files", Fore.CYAN)
        writeln("wcosa      -> Internal files used for build process", Fore.CYAN)
        writeln("Do not touch bin and wcosa folder", Fore.YELLOW)


    def update_cosa(self, board):
        """Updates WCosa project"""
        pass

    def handle_args(self, args):
        """Allocates tasks for creating and updating based on the args received"""

        if "-create" in args:
            self.create_cosa(args["-create"])
        elif "-update" in args:
            self.update_cosa(args["-update"])


if __name__ == '__main__':
    Create().start()
