# Script for delete useless files
# for example file what server created
# like database.db, players directory and etc

import os
import shutil

paths = ["database.db", "players"]

for path in paths:
    if os.path.exists(path):
        if os.path.isdir(path):
            try:
                shutil.rmtree(path)
                print(f"Directory '{path}' deleted succesfully")
            except OSError as e:
                print(f"An error occurred while deleting '{path}': {e}")
        elif os.path.isfile(path):
            try:
                os.remove(path)
                print(f"File '{path}' deleted succesfully")
            except PermissionError:
                    print(f"Permission denied: Unable to delete '{path}'")
            except Exception as e:
                    print(f"An error occurred while deleting '{path}': {e}")
