# -*- coding: utf-8 -*-

import os


def isEmpty(s: str) -> bool:
    if len(s) == 0:
        return True
    if len(s) == 1 and s[0] == "\n":
        return True

    return False


def setupEnv(file: str):
    try:
        f = open(file, "r")
    except Exception as e:
        print("can't open a file:", e)
        exit(1)

    for line in f.readlines():
        if isEmpty(line):
            continue

        if line[0] == "#":
            # Skip commented lines
            continue

        if line[-1] == "\n":
            # Trim \n
            line = line[:len(line)-1]

        key, value = line.split("=", 1)
        os.environ[key] = value

    f.close()


setupEnv("./scripts/run/go/run.env")

try:
    os.system("go run -mod=vendor main.go")
except KeyboardInterrupt:
    pass
