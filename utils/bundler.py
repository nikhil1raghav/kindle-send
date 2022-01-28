import os
import subprocess
from config import store_path
from ebooklib import epub
def bundler(urls, title):
    cnt = 0
    #downloaded all epubs
    command = ["percollate","epub","--wait=1"]
    for url in urls:
        command.append(url);
    command.append("--title={}".format(title))
    command.append("-o")
    command.append(store_path+title)
    
    subprocess.call(command)


    



