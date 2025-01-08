import argparse

from fvscover import do_cover_pkg

parser = argparse.ArgumentParser(prog="fv.scripts", description="fv dev scripts")

subparsers = parser.add_subparsers()

cover_parse = subparsers.add_parser("cover", description="run code cover test")
cover_parse.set_defaults(__action=do_cover_pkg)
cover_parse.add_argument(
    "-p", choices=(".", "vld"), type=str, help="which package", required=True
)

if __name__ == "__main__":
    args = parser.parse_args()
    if getattr(args, "__action"):
        getattr(args, "__action")(args)
