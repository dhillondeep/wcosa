"""
Handle creating and update of WCosa projects
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
        self.curr_path = self.linux_path(os.getcwd())
        self.dir_name = os.path.basename(self.curr_path)
        self.wcosa_path = self.linux_path(os.path.abspath(os.path.dirname(
            os.path.abspath(__file__)) + "/.."))



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

        # create cmake folder
        if not os.path.exists(self.curr_path + "/required"):
            os.makedirs(self.curr_path + "/required")

        # copy config file to main dir
        copyfile(self.wcosa_path + "/scripts/required-files/config.json",
                 self.curr_path + "/config.json")

        writeln("done")


        """

        

        print(self.Colors.OKGREEN + "Creating project configuration" + self.Colors.ENDC)
        print("You will now have to provide the folder path where avr is stored")
        print(self.Colors.BOLD + "Make sure you provide the path to the first instance of avr" + self.Colors.ENDC)
        print(
            self.Colors.BOLD + "Example: " + self.Colors.ENDC + "C:\\Program Files (x86)\\Arduino\\hardware\\tools\\avr")
        avr_path = input("Please type avr path: ")
        avr_path = os.path.abspath(avr_path)

        # load json and write changes
        config_file = open(curr_path + "/config.json")
        config_data = json.load(config_file)
        config_data["board"] = board
        config_data["avr-path"] = avr_path
        config_file.close()
        config_file = open(curr_path + "/config.json", "w")
        json.dump(config_data, config_file)
        config_file.close()

        print(
            self.Colors.OKBLUE + "Avr path set to " + avr_path + " and the board is " + board + self.Colors.ENDC + "\n")

        print(self.Colors.OKGREEN + "Setting up build system with cmake and make ... " + self.Colors.ENDC, end="")

        # set up CMakeLists and build system
        # copy CMakeLists file to main dir if it does not exist
        if not os.path.exists(curr_path + "/lib/CMakeLists.txt"):
            copyfile(curr_path + "/scripts/required-files/CMakeLists-lib.txt", curr_path + "/lib/CMakeLists.txt")
        if not os.path.exists(curr_path + "CMakeLists.txt"):
            copyfile(curr_path + "/scripts/required-files/CMakeLists-main.txt", curr_path + "/CMakeLists.txt")

        print("done!")
        print(self.Colors.HEADER + "---Creation of WCosa finished!!---" + self.Colors.ENDC)
        """

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
    create = Create()
    create.start()
