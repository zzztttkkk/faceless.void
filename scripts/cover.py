import typing

import consts

print(consts.PROJECT_ROOT)

if typing.TYPE_CHECKING:
    from argparse import Namespace


def do_cover(args: "Namespace"):
    match args.p:
        case "vld":
            return do_vld_cover()


def do_vld_cover():
    print("run code cover from vld")


def _merge_cover_outs(*fps: str):
    for x in fps:
        print(x)
