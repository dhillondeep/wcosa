"""Helper functions to be used through the app"""

import json
import os
import shutil


def linux_path(path, operating_system):
    """Converts Windows style path to linux style path"""

    path = path.replace("\\", "/")

    return path


def fill_template(file_path, data):
    """Fills the template cmake files"""

    file = open(file_path)
    file_str = file.read()
    file.close()

    for key in data:
        # ignore the comment
        if key == "__comment__":
            continue

        file_str = file_str.replace("%" + key, data[key])

    file = open(file_path, "w")
    file.write(file_str)
    file.close()


def create_folder(path, override=False):
    """Creates a folder at the given path"""

    if override is True:
        if os.path.exists(path):
            shutil.rmtree(path)
        os.mkdir(path)
    elif os.path.exists(path) is False:
        os.mkdir(path)


def write_conf(path, board, framework):
    config_file = open(path)
    config_data = json.load(config_file)
    config_data["board"] = board
    config_data["framework"] = framework
    config_file.close()
    config_file = open(path, "w")
    json.dump(config_data, config_file)
    config_file.close()
