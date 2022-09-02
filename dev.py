#!python3

import multiprocessing
import subprocess
import sys

from pathlib import Path

FLAGS    = ['-std=c++17']
SRC      = './src'
TESTS    = './test'
BUILD    = './build'
BIN      = './bin'
ROOT     = Path('./').resolve()
PARALLEL = True

def setup():
    # make build directory
    Path(BIN).mkdir(exist_ok=True)
    build_dir = Path(BUILD)
    build_dir.mkdir(exist_ok=True)
    for d in [SRC, TESTS]:
        dir = build_dir / d
        dir.mkdir(exist_ok=True)

def resolve_path(path: str):
    return Path(path).resolve().relative_to(ROOT).as_posix()

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
    def stale(self):
        return self.obj_ts <= self.src_ts or self.obj_ts <= self.max_inc_ts

    @property
    def obj_ts(self):
        return ts_or(Path(self.obj), 0)

    @property
    def src_ts(self):
        return ts_or(Path(self.src), 0)

    @property
    def max_inc_ts(self):
        return max(ts_or(Path(p), 0) for p in self.incs)

    def compile_tu(self) -> bool:
        if not self.stale:
            print(f'{self.obj} is not stale!')
            return False
        print(f'compiling {self.obj} from {self.src}...')
        _, err = run_compile(self.src, self.obj)
        if len(err) > 0:
            print(err)
            return False
        return True

    def compile_binary(self, binary, dag, parallel: bool) -> bool:
        tus = [self]
        for h in self.incs:
            tus.append(dag[header_to_tu(h)])

        if parallel:
            cores = multiprocessing.cpu_count()
            pool = multiprocessing.Pool(cores)
            res = pool.map(TU.compile_tu, tus)
        else:
            res = [tu.compile_tu() for tu in tus]

        _, err = run_link(binary, [tu.obj for tu in tus])
        if len(err) > 0:
            print(err)

def parse_dep_output(output: str):
    obj, deps = output.split(':')
    src, *incs = deps.split()               # source file and included headers
    src = resolve_path(src)
    name = src.removesuffix('.cc')          # tu identifier
    target = BUILD / Path(src).parent / obj # object output
    return name, target.as_posix(), src, [resolve_path(inc) for inc in incs]

def get_deps(file):
    out, _ = run('g++', *FLAGS, '-MM', file)
    return parse_dep_output(out)

def get_files(dir, pattern='*'):
    return Path(dir).glob(f'{pattern}')

def run_compile(src, obj):
    return run('g++', *FLAGS, '-c', src, f'-o{obj}')

def run_link(binary, objects):
    print(f'linking {binary} from {objects}')
    return run('g++', *FLAGS, *objects, f'-o{binary}')

def create_dag(dirs):
    tree = {}
    for dir in dirs:
        for source in get_files(dir, '*.cc'):
            name, obj, src, incs = get_deps(source.as_posix())
            tree[name] = TU(obj, src, incs)
    return tree

def header_to_tu(header: str):
    return header.removesuffix(".h")

def build_executable(target: str, output: str, dirs):
    dag = create_dag(dirs)
    tu: TU = dag[target]
    tu.compile_binary(output, dag, PARALLEL)

def build_main(name: str):
    build_executable('src/main', name, [SRC])

def build_test(name: str):
    build_executable(f'test/{name}', f'{BIN}/{name}', [TESTS, SRC])

if __name__ == '__main__':
    setup()
    match sys.argv[1:]:

        case ['test', name, *opts]:
            print(f'building test {name}...')
            build_test(name)

        case ['main', *opts]:
            print(f'compiling ape...')
            build_main('ape')

        case other:
            print(f'invalid input: {other}')
