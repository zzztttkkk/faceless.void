import typing
import tempfile
import os

import fvsconsts
import fvscommon

if typing.TYPE_CHECKING:
    from argparse import Namespace
    from pathlib import Path


def do_cover_pkg(args: "Namespace"):
    match args.p:
        case "vld":
            return _do_vld_cover()


def _do_vld_cover():
    print("run code cover from vld")

    def exec(threshold: int) -> str:
        file_content = f"""package vld
const (
    PerferPtrVldSizeThreshold = uintptr({threshold})
)
"""
        with open(
            fvsconsts.PROJECT_ROOT.joinpath("vld/perferptr.go"),
            mode="w+",
            encoding="utf8",
        ) as f:
            f.truncate(0)
            f.write(file_content)

        out = f"./vld/cover.{threshold}.out"
        fvscommon.must(f"go test -coverprofile {out} {fvsconsts.ROOT_PKG_NAME}/vld")
        return out

    outs = [
        exec(0),
        exec(128),
    ]
    _merge_cover_outs(*outs)

    for fp in outs:
        os.remove(fp)


class _CoverFile(typing.TypedDict):
    modeline: str
    lines: bytes


def _read_cover_out_file(fp: typing.Union[str, "Path"]) -> _CoverFile:
    with open(fp, mode="rb") as f:
        modeline = f.readline()
        lines = f.read()
        return {"modeline": (modeline.decode("utf8")).strip(), "lines": lines}


def _merge_cover_outs(*fps: typing.Union[str, "Path"]):
    modeline: str | None = None
    lines: list[bytes] = []
    for fp in fps:
        cf = _read_cover_out_file(fp)
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
