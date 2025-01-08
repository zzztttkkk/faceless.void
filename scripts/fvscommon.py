import os


def must(cmd: str):
    code = os.system(cmd)
    if code == 0:
        return
    raise ValueError(f"exec `{cmd}` failed, code: {code}")
