# -*- coding: utf-8 -*-

# This script builds backend and starts a Docker container

import os
import argparse


def runCommand(cmd: str, errMsg: str = ""):
    '''Calls os.system(cmd). If it gets an exception, it calls exit(1)'''
    try:
        os.system(cmd)
    except Exception as e:
        print("[ERR] ", e)
        exit(1)


# Args
parser = argparse.ArgumentParser(description="This script builds backend and starts a Docker container")

parser.add_argument("--image-name",
                    type=str,
                    default="dev-tags-drive",
                    help="name of a Docker image (default: 'dev-tags-drive')")
parser.add_argument("--image-tag",
                    type=str,
                    default="latest",
                    help="tag of a Docker image (default: 'lastest')")
parser.add_argument("--container-name",
                    type=str,
                    default="dev-tags-drive",
                    help="name of a Docker container (default: 'dev-tags-drive')")
parser.add_argument("--container-port",
                    type=int,
                    default=80,
                    help="port of a Docker container (default: 80)")
parser.add_argument("--mount-folder",
                    type=str,
                    default=os.path.realpath("."),
                    help="folder for mount to a Docker container (default: root folder)")
parser.add_argument("--build-only",
                    action="store_true",
                    help="don't run s Docker container (default: False)")

# Available args:
#    image_name
#    image_tag
#    container_name
#    container_port
#    mount_folder
#    build_only
args = parser.parse_args()


def buildDockerImage():
    # Build a Docker image (run in root folder)
    runCommand(f"docker build -t {args.image_name}:{args.image_tag} -f scripts/docker/Dockerfile .",
               "can't build a Docker image")


def runDockerContainer():
    runCommand("docker run -d --rm " +
               f"--name {args.container_name} " +
               f"-p {args.container_port}:80 " +
               f"-v {args.mount_folder}/var:/app/var " +
               f"-v {args.mount_folder}/var/data:/app/data " +
               f"{args.image_name}:{args.image_tag}")


if __name__ == "__main__":
    print("[INF] build a Docker image")

    buildDockerImage()

    if not args.build_only:
        print("[INF] run a Docker container")
        runDockerContainer()
    else:
        print("[INF] skip running a Docker container")
