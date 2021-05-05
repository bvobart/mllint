import subprocess
import sys
import os

def mllint() -> int:
  """Runs the mllint executable and returns its exit code"""
  dirname, _ = os.path.split(__file__)
  mllint_exe = os.path.join(dirname, 'mllint-exe')
  os.chmod(mllint_exe, os.stat(mllint_exe).st_mode | 0o111) # Ensures mllint-exe is executable, equivalent to `chmod +x`
  return subprocess.run([mllint_exe] + sys.argv[1:], check=False).returncode

def main():
  """Runs the mllint executable and exits with the exit code it returned"""
  sys.exit(mllint())
