# -*- coding: utf-8 -*-

import os


def setupEnv(file: str):
    try:
        f = open(file, "r")
    except Exception as e:
        print("can't open a file:", e)
        exit(1)

    for line in f.readlines():
        if line[-1] == "\n":
            # Trim \n
            line = line[:len(line)-1]

        key, value = line.split("=", 1)
        os.environ[key] = value

    f.close()


setupEnv("./scripts/run/run.env")

try:
    os.system("go run -mod=vendor main.go")
except KeyboardInterrupt:
    pass
