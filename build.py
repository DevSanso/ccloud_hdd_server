import os
import sys
import subprocess
import platform


current_dir = os.path.dirname(os.path.realpath(__file__))
dst_dir = os.path.join(current_dir,"dst")

class Flag:
    def __init__(self,args = list()):
        self.m = dict()
        self.m["server"] = False
        self.m["cli"] = False
        for arg in args:
            self.m[arg] = True

    def is_server(self):
        return self.m["server"]
    def is_cli(self):
        return self.m["cli"]

def is_dst_dir():
    return os.path.exists(dst_dir)

def is_windows():
    return platform.system() == "Windows"

def source_real_path(source_dir,file_list=list()):
    res=list()
    for source in file_list:
        res.append(os.path.join(source_dir,source))
    return res

def build_server():
    source_dir=os.path.join(current_dir, "cmd", "server")
    source_list = source_real_path(source_dir,os.listdir(source_dir))
    output=os.path.join(dst_dir,"hdd_server")

    if is_windows():
        output = output + ".exe"
    args = ["go","build","-o",output] + source_list
    subprocess.run(args)

def build_cli():
    source_dir = os.path.join(current_dir, "cmd", "admin_cli")
    source_list = source_real_path(source_dir,os.listdir(source_dir))
    output = os.path.join(dst_dir,"hdd_server_cli")

    if is_windows():
        output = output + ".exe"

    args = ["go", "build", "-o", output] + source_list
    subprocess.run(args)

if __name__ == "__main__":
    flag = Flag(sys.argv)
    if flag.is_server():
       build_server()

    if flag.is_cli():
        build_cli()