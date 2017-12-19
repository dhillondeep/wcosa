"""
Handle creating and updating WCosa projects
"""
import os
import json
from glob import glob
from shutil import copyfile
from colorama import Fore
from module.parent import Parent
from module.output import write, writeln
import module.helper as helper


class Handler(Parent):
    """Handles creating and update of WCosa projects"""

    def __init__(self):
        """Initialize paths to use for creation"""

        super(Handler, self).__init__()
        self.curr_path = helper.linux_path(os.getcwd(), self.operating_system)
        self.dir_name = os.path.basename(self.curr_path)
        self.wcosa_path = helper.linux_path(os.path.abspath(os.path.dirname(
            os.path.abspath(__file__)) + "/.."), self.operating_system)
        self.cmake_templates_path = self.wcosa_path + "/scripts/cmake-files"
        self.config_files_path = self.wcosa_path + "/scripts/config-files"

    def create_json(self, board_path):
        board_file = open(board_path)
        board_json = json.load(board_file)
        fill_json = {"os": self.operating_system, "cmake-version": "2.8",
                     "project-name": self.dir_name, "wcosa-path": self.wcosa_path}

        if self.operating_system == "windows" or self.operating_system == "cygwin":
            arduino_sdk_path = os.environ.get("ARDUINO_SDK_PATH")

            if arduino_sdk_path is None:
                writeln("Install Arduino SDK and add ARDUINO_SDK_PATH as system variable")
            fill_json["avr-path"] = helper.linux_path(
                arduino_sdk_path, self.operating_system) + "/hardware/tools/avr"
        else:
            fill_json["avr-path"] = "None"

        fill_json["avr-mcu"] = board_json["mcu"]
        fill_json["avr-h-fuse"] = board_json["hfuse"]
        fill_json["avr-l-fuse"] = board_json["lfuse"]

        # rule to gather library paths
        include_lib_str = ""
        libraries = glob(self.curr_path + "/lib/*/")

        for library in libraries:
            library = helper.linux_path(library, self.operating_system)
            files = glob(library + "/*/")

            source_exists = False
            for directory in files:
                directory = helper.linux_path(directory, self.operating_system)

                if directory == library + "src/" and source_exists is not True:
                    source_exists = True
                    include_lib_str += "include_directories(\"" + directory + "\")\n"
                elif source_exists is True:
                    break

            if source_exists is not True:
                include_lib_str += "include_directories(\"" + library + "\")\n"

        fill_json["include-directories"] = include_lib_str

        # get custom flags from the config file
        config_file = open(self.curr_path + "/config.json")
        config_data = json.load(config_file)
        flags = config_data["build-flags"]
        config_file.close()

        if len(flags) >= 1:
            fill_json["custom-compiler-flags"] = "add_definitions(\"" + flags + "\")"
        else:
            fill_json["custom-compiler-flags"] = ""

        return fill_json

    def create_cosa(self, board):
        """Creates WCosa project"""

        write("Creating work environment - ", color=Fore.CYAN)

        # create src, lib, bin and wcosa folders
        helper.create_folder(self.curr_path + "/src")
        helper.create_folder(self.curr_path + "/lib")
        helper.create_folder(self.curr_path + "/bin", True)
        helper.create_folder(self.curr_path + "/wcosa")

        # copy cmake files and config files
        copyfile(self.cmake_templates_path + "/CMakeLists.txt",
                 self.curr_path + "/CMakeLists.txt")
        copyfile(self.cmake_templates_path + "/generic-gcc-avr.cmake",
                 self.curr_path + "/wcosa/generic-gcc-avr.cmake")
        copyfile(self.config_files_path + "/config.json",
                 self.curr_path + "/config.json")

        writeln("done")
        write("Creating project configuration - ", Fore.CYAN)

        helper.write_conf(self.curr_path + "/config.json", board, "wcosa")

        # use internal config file based on the board to fill the template
        data = self.create_json(self.wcosa_path + "/scripts/boards/" + board + ".json")
        helper.fill_template(self.curr_path + "/CMakeLists.txt", data)

        writeln("done")

        writeln("Finished Creation: ", Fore.CYAN)
        writeln("src        -> Source files", Fore.CYAN)
        writeln("lib        -> Library files (each library in seperate folder)", Fore.CYAN)
        writeln("bin        -> Binary files", Fore.CYAN)
        writeln("wcosa      -> Internal files used for build process", Fore.CYAN)
        writeln("Do not touch bin and wcosa folder", Fore.YELLOW)

    def update_cosa(self, newBoard):
        writeln("Updating " + self.dir_name + " project: ", Fore.CYAN)

        # create src, lib, bin and wcosa folders
        helper.create_folder(self.curr_path + "/src")
        helper.create_folder(self.curr_path + "/lib")
        helper.create_folder(self.curr_path + "/bin", True)
        helper.create_folder(self.curr_path + "/wcosa")

        # copy cmake files and config files
        copyfile(self.cmake_templates_path + "/generic-gcc-avr.cmake",
                 self.curr_path + "/wcosa/generic-gcc-avr.cmake")

        cmake_file = open(self.curr_path + "/CMakeLists.txt")
        lines = cmake_file.readlines()
        cmake_file.close()
        cmake_file = open(self.curr_path + "/CMakeLists.txt", "w")

        # delete all the includes and add a template in their place so that new includes can be added
        added_first_time = False
        for line in lines:
            if "include_directories" in line:
                if not added_first_time:
                    cmake_file.write("%include-directories")
                    added_first_time = True
                continue
            elif "# building library and adding src (do not delete this line)" in line and not added_first_time:
                cmake_file.write("%include-directories")
                cmake_file.write("\n##########################################################################")
                cmake_file.write("\n\n\n")
                cmake_file.write("##########################################################################")
                cmake_file.write("\n" + line)
            else:
                cmake_file.write(line)
        cmake_file.close()

        # add compiler flag template
        cmake_file = open(self.curr_path + "/CMakeLists.txt")
        lines = cmake_file.readlines()
        cmake_file.close()
        cmake_file = open(self.curr_path + "/CMakeLists.txt", "w")

        # delete all the includes and add a template in their place so that new includes can be added
        last_add_definition = False
        done = False
        counter = 0
        for line in lines:
            if 0 < counter <= 2:
                counter += 1
                cmake_file.write("\n")
                continue

            if "add_definitions" in line:
                last_add_definition = True
            elif line == "#do not touch this line###################################################\n" and \
                    last_add_definition is True and done is False:
                cmake_file.write("#do not touch this line###################################################")
                cmake_file.write("\n%custom-compiler-flags")
                done = True
                counter += 1
                continue
            elif "\n" == line and last_add_definition is True and done is False:
                cmake_file.write("#do not touch this line###################################################")
                cmake_file.write("\n%custom-compiler-flags")
                done = True
            else:
                last_add_definition = False

            cmake_file.write(line)
        cmake_file.close()

        if newBoard == "" or newBoard == "this":
            config_file = open(self.curr_path + "/config.json")
            config_data = json.load(config_file)
            board = config_data["board"]
            config_file.close()
            writeln("Board is retrieved from project configuration: " + board, Fore.CYAN)
        else:
            helper.write_conf(self.curr_path + "/config.json", newBoard, "cosa")
            board = newBoard
            writeln("Board is changed to the provided board: " + board, Fore.CYAN)

        # use internal config file based on the board to fill the template
        data = self.create_json(self.wcosa_path + "/scripts/boards/" + board + ".json")
        helper.fill_template(self.curr_path + "/CMakeLists.txt", data)

        writeln("Update complete", Fore.YELLOW)

    def handle_args(self, args):
        """Allocates tasks for creating and updating based on the args received"""

        if "-path" in args:
            self.curr_path = helper.linux_path(args["-path"], self.operating_system)

        if "-create" in args:
            self.create_cosa(args["-create"])
        elif "-update" in args:
            self.update_cosa(args["-update"])


if __name__ == '__main__':
    Handler().start()
