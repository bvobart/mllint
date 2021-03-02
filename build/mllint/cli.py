import subprocess
import sys
import os

def mllint() -> int:
  dirname, _ = os.path.split(__file__)
  mllint_exe = os.path.join(dirname, 'mllint-exe')
  return subprocess.run([mllint_exe] + sys.argv[1:]).returncode