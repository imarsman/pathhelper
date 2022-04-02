# pathhelper
Helper for building a path for MacOS.

Linux distributions do not store system paths as lists of files in `/etc/paths`, `/etc/manpaths` or in `/etc/paths.d/`,
or `/etc/manpaths.d`, so this would not really work to set the path in a `.zshrc` file.

In addition to the work that `/usr/libexec/path_helper` does on MacOS `pathhelper` also looks in `~/.config/pathhelper/`
for files in `paths.d` and `manpaths.d`.

The converntion is that there will be one path per line in the main `paths` and `manpaths` and in files in the `paths.d`
and `manpaths.d` directories. More than one line in files is permitted, and lines started with `#` are ignored.

When looking for paths in files, each path encountered is checked to ensure it exists or these non-existent files files
are ignonred.

I wrote this to better understand how the `PATH` and `MANPATH` variables are set in MacOS. The
`/usr/libexec/path_helper` binary runs on my laptop in about 6 msec. pathhelper takes about 11 msec to run on my laptop.

Here is a sample setup in `.zshrc`.

```
# Use pathhelper to get system paths plus user paths
eval $(~/bin/pathhelper -z)
```

Paths made as a result of a call in `.zshrc` must be added to the main `PATH` separately.

Here is the help output for `pathhelper`

```
$ pathhelper -h
Usage: pathhelper [--bash] [--csh] [--zsh] [--verbose] [--init]

Options:
  --bash, -s
  --csh, -c
  --zsh, -z
  --verbose, -v
  --init, -i
  --help, -h             display this help and exit
```

Note that formats compatible with `csh`, `bash`, and `zsh` are available using flags. There is also an init method which
will make `.config/..` directories if they do not exist. Using the `-v` (verbose) flag will show entries that were
rejected because they did not exist.