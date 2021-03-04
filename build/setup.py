import os
import platform
import setuptools
import shutil

#-------------------------------------------------------------------------------

def get_mllint_exe() -> str:
  """
  Get the platform-specific filename of the compiled mllint executable,
  or raise an exception if the platform is unsupported.
  """
  system, _, _, _, machine, _ = platform.uname()

  # Windows
  if system == 'Windows' and (machine == 'i386' or machine == 'i686'):
    return os.path.join('bin', 'mllint-windows-386')

  elif system == 'Windows' and machine == 'AMD64':
    return os.path.join('bin', 'mllint-windows-amd64')

  # MacOS
  elif system == 'Darwin' and machine == 'x86_64':
    return os.path.join('bin', 'mllint-darwin-amd64')

  # Linux
  elif system == 'Linux' and (machine == 'i386' or machine == 'i686'):
    return os.path.join('bin', 'mllint-linux-386')

  elif system == 'Linux' and machine == 'x86_64':
    return os.path.join('bin', 'mllint-linux-amd64')

  # Other OSes are not supported right now, might be able to support more if the Go compiler supports it and people want it.
  else:
    print()
    print('Sorry, your OS is not supported. mllint currently supports:')
    print('- Linux (32 or 64-bit x86)')
    print('- Windows (32 or 64-bit x86)')
    print('- MacOS (only 64-bit x86)')
    print()
    print(f'Your OS: {system} ({machine})')
    print()
    raise Exception(f'unsupported OS: {system} ({machine})')

#-------------------------------------------------------------------------------

def patch_distutils():
  """
  distutils.util.change_root() has a bug on Windows where it fails with a string index out of range error
  when the pathname is empty. To work around this, we need to monkey-patch change_root,
  which is what this function does.

  Also, distutils.command.install checks for truthy ext_modules instead of calling has_ext_modules()
  on the distribution that we're passing into setup() to force platform specific builds.
  That causes the mllint-exe to be put in `purelib` data, which audit_wheel later complains about.
  So we monkey-patch that too to force installing to `platlib` data.
  """
  import distutils.util
  original_change_root = distutils.util.change_root
  from distutils.command.install import install

  def change_root(new_root, pathname):
    if os.name != 'nt': # if not Windows, just use the original change_root
      return original_change_root(new_root, pathname)
    
    # else, if we're on Windows:
    (_, path) = os.path.splitdrive(pathname)
    if path and path[0] == '\\':
      path = path[1:]
    return os.path.join(new_root, path)
  
  # From: https://github.com/bigartm/bigartm/issues/840#issuecomment-342825690
  class InstallPlatlib(install):
    def finalize_options(self):
      install.finalize_options(self)
      if self.distribution.has_ext_modules():
        self.install_lib = self.install_platlib
  
  distutils.util.change_root = change_root
  distutils.command.install.install = InstallPlatlib


class PlatformSpecificDistribution(setuptools.Distribution):
  """Distribution which always forces a platform-specific package"""
  def has_ext_modules(self):
      return True

patch_distutils()

#-------------------------------------------------------------------------------

# Copy mllint-exe into the package.
exe_path = get_mllint_exe()
if not os.path.exists(exe_path):
  print()
  print(f'Expected to find a compiled mllint binary at {exe_path} but it did not exist!')
  print("> If you're compiling mllint from source, run 'make build-all' before building this package, or just run 'make package'")
  print("> If you're installing mllint using 'pip install', then it seems pip downloaded the source package, instead of a platform-specific wheel.")
  print()
  raise Exception(f'Expected to find a compiled mllint binary at {exe_path} but it did not exist!')

shutil.copy2(exe_path, os.path.join('mllint', 'mllint-exe'))

# Include ReadMe as long description
with open("ReadMe.md", "r", encoding="utf-8") as fh:
  long_description = fh.read()

setuptools.setup(
  name="mllint",
  version="0.1.2",
  author="Bart van Oort",
  author_email="bart@vanoort.is",
  description="Linter for Machine Learning projects",
  license='MIT',
  long_description=long_description,
  long_description_content_type="text/markdown",
  url="https://gitlab.com/bvobart/mllint",
  project_urls={
      "Bug Tracker": "https://gitlab.com/bvobart/mllint/-/issues",
  },
  classifiers=[
      "Development Status :: 2 - Pre-Alpha",
      "Environment :: Console",
      "Intended Audience :: Developers",
      "Intended Audience :: Information Technology",
      "Intended Audience :: Science/Research",
      "License :: OSI Approved :: MIT License",
      "Natural Language :: English",
      "Operating System :: MacOS",
      "Operating System :: Microsoft :: Windows",
      "Operating System :: POSIX :: BSD :: FreeBSD",
      "Operating System :: POSIX :: Linux",
      "Programming Language :: Python :: 3",
      "Topic :: Scientific/Engineering :: Artificial Intelligence",
      "Topic :: Scientific/Engineering :: Information Analysis",
      "Topic :: Software Development :: Build Tools",
      "Topic :: Software Development :: Pre-processors",
      "Topic :: Software Development :: Quality Assurance",
      "Topic :: Software Development :: Version Control :: Git",
  ],
  packages=['mllint'],
  package_data={'mllint': ['mllint-exe']},
  python_requires=">=3.6",
  entry_points={
    'console_scripts': [
      'mllint=mllint.cli:main'
    ],
  },
  distclass=PlatformSpecificDistribution,
)