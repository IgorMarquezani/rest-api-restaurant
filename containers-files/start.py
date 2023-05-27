# 1/bin/python3

import subprocess
import sys

start_container = ['docker-compose', 'up']

stop_container = ['docker-compose', 'down']

list_images = ['docker', 'image', 'ls']

force_image_remove = ['docker', 'image', 'rmi', '-f']

go_run = ['go', 'run']

def start_api():
    message = subprocess.run('pwd', stdout=subprocess.PIPE).stdout.decode()
    message = message.split('/')[1:message.split('/').__len__()-1]

    path = ['/']
    for msg in message:
        path += msg.__str__() + '/'

    subprocess.Popen(start_container)

    status = subprocess.call(go_run + [''.join(path) + 'main.go'])
    while status == 1:
        status = subprocess.call(go_run + [''.join(path) + 'main.go'])


def reset_containers():
    subprocess.call(stop_container)
    containers = subprocess.check_output(list_images).decode().splitlines()[1:]

    for container in containers:
        if "postgres" in container or "dpage/pgadmin4" in container:
            subprocess.call(force_image_remove + [container.split()[3]])
            break

    subprocess.call(['docker', 'volume', 'rm', '$(docker',
                    'volume', 'ls', '-q)', '2>', '/dev/null'])


def show_help():
    print("help:  show this message")
    print("start: start the containers")
    print("reset: completely reset containers, images and volumes")
        

if len(sys.argv) > 1:
    if "shutdown" in sys.argv[1:]:
        show_help()

    if "reset" in sys.argv[1:]:
        reset_containers()

    if "start" in sys.argv[1:]:
        start_api()

    if "help" in sys.argv[1:]:
        show_help()

if len(sys.argv) <= 1:
    print("Please provide a argument (start, shutdown, reset or help).")
