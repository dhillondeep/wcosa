"""
Handle creating and updating WCosa projects
"""
import os
import json
from collections import OrderedDict
from glob import glob
from shutil import copyfile
from shutil import which
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
        self.cmake_templates_path = self.wcosa_path + "/build/cmake-files"
        self.config_files_path = self.wcosa_path + "/build/config-files"
        self.ide = "None"

    def get_avr_paths(self):
        """Get Paths of AVR libraries and tools"""

        avr = []
        if self.operating_system == "windows" or self.operating_system == "cygwin":
            arduino_sdk_path = os.environ.get("ARDUINO_SDK_PATH")

            print(arduino_sdk_path)

            if arduino_sdk_path is None:
                write("ERROR: AVR tools and libraries not found,", Fore.RED)
                writeln("please install Arduino SDK and add ARDUINO_SDK_PATH as system variable", Fore.RED)
                quit(2)
            else:
                avr.insert(0, helper.linux_path(arduino_sdk_path, self.operating_system) + "/hardware/tools/avr/avr")
                avr.insert(1, helper.linux_path(arduino_sdk_path, self.operating_system) + "/hardware/tools/avr/bin")
        else:
            if os.environ.get("AVR_FIND_ROOT_PATH") is not None:
                avr.insert(0, os.environ.get("AVR_FIND_ROOT_PATH"))
            elif os.path.exists("/opt/local/avr"):
                avr.insert(0, "/opt/local/avr")
            elif os.path.exists("/usr/avr"):
                avr.insert(0, "/user/avr")
            elif os.path.exists("/usr/lib/avr"):
                avr.insert(0, "/user/lib/avr")
            elif os.path.exists("/usr/local/CrossPack-AVR"):
                avr.insert(0, "/usr/local/CrossPack-AVR")
            else:
                write("ERROR: AVR libraries not found,", Fore.RED)
                writeln("please set AVR_FIND_ROOT_PATH in your environment", Fore.RED)
                quit(2)

            which_result = which("avr-gcc")
            if which_result is not None:
                avr.insert(1, which_result)
            else:
                writeln("ERROR: AVR tools not found, please install them")
                quit(2)

        return avr

    def update_internal_config(self, config_path, board_path):
        """Fill internal configuration file with all the info that is neede"""

        config_file = open(config_path)
        config_data = json.load(config_file, object_pairs_hook=OrderedDict)
        config_file.close()

        board_file = open(board_path)
        board_data = json.load(board_file,)
        board_file.close()

        user_config_file = open(self.curr_path + "/config.json")
        user_config_data = json.load(user_config_file)
        user_config_file.close()

        # get tool addresses based on operating system
        avr = self.get_avr_paths()

        # go through all the tags and update
        config_data["os"] = self.operating_system
        config_data["ide"] = self.ide
        config_data["wcosa-path"] = self.wcosa_path
        config_data["project-name"] = self.dir_name
        config_data["cmake-version"] = "2.8"
        config_data["avr-path"] = avr[0]
        config_data["avr-tool-path"] = avr[1]
        config_data["custom-definitions"] = board_data["board-flag"]

        if "build-flags" in user_config_data and not user_config_data["build-flags"] == "":
            config_data["custom-definitions"] += " " + user_config_data["build-flags"]

        config_data["board"] = user_config_data["board"]
        config_data["avr-mcu"] = board_data["mcu"]
        config_data["avr-h-fuse"] = board_data["hfuse"]
        config_data["avr-l-fuse"] = board_data["lfuse"]
        config_data["avr-f-cpu"] = board_data["f_cpu"]

        config_file = open(config_path, "w")
        json.dump(config_data, config_file, indent=4)
        config_file.close()

        writeln(which("avr-gcc.exe"))

    def parse_library_includes(self, lib_path, cmake_path):
        """Gathers library paths based on their folder structure and adds them to cmake file"""

        # rule to gather library paths
        include_lib_str = ""
        libraries = glob(lib_path + "/*/")

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

        cmake_file = open(cmake_path)
        cmake_str = cmake_file.read()
        cmake_file.close()

        cmake_str = cmake_str.replace("%" + "custom-directories-include", include_lib_str)

        cmake_file = open(cmake_path, "w")
        cmake_file.write(cmake_str)
        cmake_file.close()

    def update_clion_cmake(self, config_path):
        """Updates the template strings in clion CMakeLists file"""

        config_file = open(config_path)
        config_data = json.load(config_file, object_pairs_hook=OrderedDict)
        config_file.close()

        helper.fill_template(self.curr_path + "/CMakeLists.txt", config_data)
        helper.fill_template(self.curr_path + "/CMakeListsPrivate.txt", config_data)

        # parse libraries and the add them to cmake
        self.parse_library_includes(self.curr_path + "/lib", self.curr_path + "/CMakeListsPrivate.txt")

    def update_build_cmake(self, cmake_path, config_path):
        config_file = open(config_path)
        config_data = json.load(config_file, object_pairs_hook=OrderedDict)
        config_file.close()

        helper.fill_template(cmake_path, config_data)

        # parse libraries and the add them to cmake
        self.parse_library_includes(self.curr_path + "/lib", cmake_path)

    def create_cosa(self, board):
        """Creates WCosa project"""

        write("Creating work environment - ", color=Fore.CYAN)

        # create src, lib, bin and wcosa folders
        helper.create_folder(self.curr_path + "/src")
        helper.create_folder(self.curr_path + "/lib")
        helper.create_folder(self.curr_path + "/wcosa")
        helper.create_folder(self.curr_path + "/wcosa/bin")
        helper.create_folder(self.curr_path + "/wcosa/cmake")
        helper.create_folder(self.curr_path + "/wcosa/config")

        # copy cmake files and config files
        copyfile(self.cmake_templates_path + "/build/CMakeLists.txt",
                 self.curr_path + "/wcosa/cmake/CMakeLists.txt")
        copyfile(self.cmake_templates_path + "/build/generic-gcc-avr.cmake",
                 self.curr_path + "/wcosa/cmake/generic-gcc-avr.cmake")
        copyfile(self.config_files_path + "/internal-config.json",
                 self.curr_path + "/wcosa/config/internal-config.json")
        copyfile(self.config_files_path + "/config.json",
                 self.curr_path + "/config.json")
        copyfile(self.wcosa_path + "/build/.gitignore",
                 self.curr_path + "/.gitignore")

        writeln("done")
        write("Creating project configuration - ", Fore.CYAN)

        # update user config file
        user_config_file = open(self.curr_path + "/config.json")
        user_config_data = json.load(user_config_file, object_pairs_hook=OrderedDict)
        user_config_file.close()

        user_config_data["board"] = board
        user_config_data["framework"] = "cosa"

        user_config_file = open(self.curr_path + "/config.json", "w")
        json.dump(user_config_data, user_config_file, indent=4)
        user_config_file.close()

        # update internal config file based on information we have
        self.update_internal_config(self.curr_path + "/wcosa/config/internal-config.json",
                                    self.wcosa_path + "/build/boards/" + board + ".json")

        # copy ide specific CMakeFile
        if self.ide is not None:
            if self.ide == "clion":
                copyfile(self.cmake_templates_path + "/clion/CMakeLists.txt",
                         self.curr_path + "/CMakeLists.txt")
                copyfile(self.cmake_templates_path + "/clion/CMakeListsPrivate.txt",
                         self.curr_path + "/CMakeListsPrivate.txt")

                # update the CMakeLists files by filling in the templates
                self.update_clion_cmake(self.curr_path + "/wcosa/config/internal-config.json")
            else:
                writeln("ERROR: This ide is not supported", Fore.RED)
                quit(2)

        # update the build CMakeLists files by filling in the templates
        self.update_build_cmake(self.curr_path + "/wcosa/cmake/CMakeLists.txt", self.curr_path + "/wcosa/config/internal-config.json")

        writeln("done")

        writeln("Finished Creation: ", Fore.CYAN)
        writeln("src        -> Source files", Fore.CYAN)
        writeln("lib        -> Library files (each library in seperate folder)", Fore.CYAN)
        writeln("bin        -> Binary files", Fore.CYAN)
        writeln("wcosa      -> Internal files used for build process", Fore.CYAN)
        writeln("Do not touch bin and wcosa folder", Fore.YELLOW)

    def update_cosa(self, newBoard):
        """update the cosa project"""
        writeln("Updating " + self.dir_name + " project: ", Fore.CYAN)

         # create src, lib, bin and wcosa folders
        helper.create_folder(self.curr_path + "/src")
        helper.create_folder(self.curr_path + "/lib")
        helper.create_folder(self.curr_path + "/wcosa")
        helper.create_folder(self.curr_path + "/wcosa/bin")
        helper.create_folder(self.curr_path + "/wcosa/cmake")
        helper.create_folder(self.curr_path + "/wcosa/config")

        # create a .gitignore file and add wcosa and CMak

        # copy cmake files, config files and gitignore file
        copyfile(self.cmake_templates_path + "/build/CMakeLists.txt",
                 self.curr_path + "/wcosa/cmake/CMakeLists.txt")
        copyfile(self.cmake_templates_path + "/build/generic-gcc-avr.cmake",
                 self.curr_path + "/wcosa/cmake/generic-gcc-avr.cmake")

        # update user config file
        user_config_file = open(self.curr_path + "/config.json")
        user_config_data = json.load(user_config_file, object_pairs_hook=OrderedDict)
        user_config_file.close()

        user_config_data["board"] = newBoard

        user_config_file = open(self.curr_path + "/config.json", "w")
        json.dump(user_config_data, user_config_file, indent=4)
        user_config_file.close()

         # update internal config file based on information we have
        self.update_internal_config(self.curr_path + "/wcosa/config/internal-config.json",
                                    self.wcosa_path + "/build/boards/" + newBoard + ".json")
        
        # get configuration from config file
        user_config_file = open(self.curr_path + "/wcosa/config/internal-config.json")
        user_config_data = json.load(user_config_file, object_pairs_hook=OrderedDict)
        user_config_file.close()

        if user_config_data["ide"] == "clion":
            copyfile(self.cmake_templates_path + "/clion/CMakeListsPrivate.txt",
                         self.curr_path + "/CMakeListsPrivate.txt")
    
            self.update_clion_cmake(self.curr_path + "/wcosa/config/internal-config.json")

        # update the build CMakeLists files by filling in the templates
        self.update_build_cmake(self.curr_path + "/wcosa/cmake/CMakeLists.txt", self.curr_path + "/wcosa/config/internal-config.json")
    
    def handle_args(self, args):
        """Allocates tasks for creating and updating based on the args received"""

        if "-path" in args:
            self.curr_path = helper.linux_path(args["-path"], self.operating_system)

        if "-ide" in args:
            self.ide = args["-ide"]

        if "-create" in args:
            self.create_cosa(args["-create"])
        elif "-update" in args:
            self.update_cosa(args["-update"])


if __name__ == '__main__':
    Handler().start()
