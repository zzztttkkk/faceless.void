import typing
import tempfile
import os

import fvsconsts
import fvscommon

if typing.TYPE_CHECKING:
    from argparse import Namespace
    from pathlib import Path


def do_cover(args: "Namespace"):
    match args.p:
        case "vld":
            return do_vld_cover()


def do_vld_cover():
    print("run code cover from vld")

    _merge_cover_outs(
        fvsconsts.PROJECT_ROOT.joinpath("vld/cover.0.out"),
        fvsconsts.PROJECT_ROOT.joinpath("vld/cover.128.out"),
    )


class CoverFile(typing.TypedDict):
    modeline: str
    lines: bytes


def read_cover_out_file(fp: typing.Union[str, "Path"]) -> CoverFile:
    with open(fp, mode="rb") as f:
        modeline = f.readline()
        lines = f.read()
        return {"modeline": (modeline.decode("utf8")).strip(), "lines": lines}


def _merge_cover_outs(*fps: typing.Union[str, "Path"]):
    modeline: str | None = None
    lines: list[bytes] = []
    for fp in fps:
        cf = read_cover_out_file(fp)
        if modeline is None:
            modeline = cf["modeline"]
        elif cf["modeline"] != modeline:
            raise ValueError()

        lines.append(cf["lines"])

    with tempfile.NamedTemporaryFile("wb+", delete_on_close=False) as tmpf:
        tmpf.write(typing.cast(str, modeline).encode("utf8"))
        tmpf.write(b"\r\n")
        for bs in lines:
            tmpf.write(bs)

        tmpf.close()
        fvscommon.must(f"go tool cover -html {tmpf.name}")
        os.remove(tmpf.name)
