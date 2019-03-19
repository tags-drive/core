# -*- coding: utf-8 -*-

# This script builds backend and starts a Docker container

import os
import shutil


def runCommand(cmd: str, errMsg: str = ""):
    '''Calls os.system(cmd). If it get an exception, it calls exit(1)'''
    try:
        os.system(cmd)
    except Exception as e:
        print(f"[FAT] {errMsg}:", e)
        exit(1)


def removeFile(file: str):
    '''Tries to remove file'''
    try:
        os.remove(file)
    except Exception as e:
        print(f"[WRN] can't delete file '{file}':", e)


# Vars
ROOT = os.path.realpath(".")
CONTAINER_NAME = "dev-tags-drive"
PORT = 80


def buildBinary():
    # Build a binary
    os.environ["GOOS"] = "linux"
    os.environ["GOARCH"] = "amd64"

    runCommand("go build -o tags-drive", "can't build binary")

    # Try to delete old version (we can't move a file, if it exists)
    removeFile("./scripts/docker/tags-drive")

    try:
        shutil.move("tags-drive", "./scripts/docker")
    except Exception as e:
        print(f"[FAT] can't move 'tags-drive':", e)
        exit(1)


def buildDockerImage():
    # Build a Docker image
    os.chdir("./scripts/docker")

    runCommand("docker build -t dev-tags-drive:latest .", "can't build a Docker image")

    # Clear
    removeFile("./tags-drive")


def runDockerContainer():
    # Run container
    runCommand(f"docker run -d --name {CONTAINER_NAME} " +
               f"-p {PORT}:80 " +
               f"-v {ROOT}/configs:/app/configs " +
               f"-v {ROOT}/data:/app/data " +
               f"--restart=unless-stopped " +
               f"dev-tags-drive:latest")


buildBinary()

buildDockerImage()

runDockerContainer()
