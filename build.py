#!python3

import subprocess
import sys

from pathlib import Path

FLAGS = ['-std=c++17']
SRC   = './src'
TESTS = './test'
BIN = './bin'
ROOT  = Path('./')

# make build directory
def init_bin_dir():
    Path(BIN).mkdir(exist_ok=True)
    for d in [SRC, TESTS]:
        dir = Path(BIN) / d
        dir.mkdir(exist_ok=True)
init_bin_dir()

def run(cmd, *args):
    proc = subprocess.Popen(
        [cmd]+list(args),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    out, err = proc.communicate()
    return out.decode('utf-8'), err.decode('utf-8')

def ts_or(path: Path, default: float = 0) -> float:
    try:
        return path.stat().st_mtime
    except:
        return default

class TU:
    def __init__(self, obj, src, incs):
        self.obj = obj
        self.src = src
        self.incs = incs

    def __str__(self):
        return f'{self.obj}: {self.src}, {self.incs}'

    def __repr__(self):
        return self.__str__()

    @property
    def target(self):
        return f'{BIN}/{self.obj}'

    @property
    def stale(self):
        return self.obj_ts <= self.src_ts or self.obj_ts <= self.max_inc_ts

    @property
    def obj_ts(self):
        return ts_or(Path(self.target), 0)

    @property
    def src_ts(self):
        return ts_or(Path(self.src), 0)

    @property
    def max_inc_ts(self):
        return max(ts_or(Path(p), 0) for p in self.incs)

    def compile(self) -> bool:
        if not self.stale:
            print(f'{self.target} is not stale!')
            return False
        print(f'compiling {self.target} from {self.src}...')
        _, err = run_compile(self.src, self.target)
        if len(err) > 0:
            print(err)
            return False
        return True

def get_deps(file):
    out, _ = run('g++', *FLAGS, '-MM', file)
    obj, deps = out.split(':')
    src, *incs = deps.split()
    src_root = Path(src).parent
    obj = src_root / obj.rstrip(':')
    return obj.as_posix(), src, incs

def get_files(dir, pattern='*'):
    return ROOT.glob(f'{dir}/{pattern}')

def run_compile(src, obj):
    return run('g++', *FLAGS, '-c', src, f'-o{obj}')

def run_link(binary, objects):
    print(f'linking {binary} from {objects}')
    return run('g++', *FLAGS, *objects, f'-o{binary}')

def create_dag(dir):
    tree = {}
    for source in get_files(dir, '*.cc'):
        obj, src, incs = get_deps(source.as_posix())
        tree[obj] = TU(obj, src, incs)
    return tree

# edges in the dag are represented by the header file implemented
# by the required translation unit, so header_to_tu converts the
# header name to the object file that identifies the translation unit
# in the future, each translation unit should have a unique ID, such
# as the translation unit for src/lexer.cc being 'src/lexer' instead
# of 'src/lexer.o'
def header_to_tu(header: str):
    return f'{header.removesuffix(".h")}.o'

def build_main(binary: str):
    dag = create_dag(SRC)
    # TODO: take the target translation unit instead of hardcoding
    main: TU = dag['src/main.o']
    objs = []
    for h in main.incs:
        obj = header_to_tu(h)
        dep_tu: TU = dag[obj]
        dep_tu.compile()
        objs.append(dep_tu.target)

    main.compile()
    objs.append(main.target)
    _, err = run_link(binary, objs)
    if len(err) > 0:
        print(err)

if __name__ == '__main__':
    binary = 'ape'
    if len(sys.argv) > 1:
        binary = sys.argv[1]
    build_main(binary)
