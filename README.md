# pathhelper
Helper for building a path for MacOS.

I wrote this to better understand how the `PATH` and `MANPATH` variables are set in MacOS. The
`/usr/libexec/path_helper` binary runs on my laptop in about 6 msec. pathhelper takes about 12 msec to run on my laptop.
Given that the load time for my `zsh` environment is 500 msec this is at the moment acceptable. `path_helper` is
[written in C](https://opensource.apple.com/source/shell_cmds/shell_cmds-162/path_helper/path_helper.c.auto.html) and
does not from what I can see do any validation of the paths found. Not checking for the existence of the directories
pointed bo by the path elements found could partly explain its faster execution.

Linux distributions do not store system paths as lists of files in `/etc/paths`, `/etc/manpaths` or in `/etc/paths.d/`,
or `/etc/manpaths.d`, so this would not really work to set the path in a `.zshrc` file. It would work as a user path
setting tool but that is limited.

In addition to the work that `/usr/libexec/path_helper` does on MacOS `pathhelper` also looks in `~/.config/pathhelper/`
for files in `paths.d` and `manpaths.d`. Here is the list of entries for `/etc/paths.d`

```sh
$ ls -1 /etc/paths.d
100-rvictl
40-XQuartz
TeX
exiftool
go
```

Looking at the several off the above entries, it can be seen that some installers such as the installer for Golang add
to this list. This likely makes for many fewer complaints that "go can't be found".

The convention is that there will be one path per line in the main `paths` and `manpaths` and in files in the `paths.d`
and `manpaths.d` directories. More than one line in files is permitted, and lines started with `#` are ignored.

When looking for paths in files, each path encountered is checked to ensure it exists and is readable or these
non-existent files files are ignonred. These invalid items are silently rejected but can be seen if the `-v` flag is
used.

Here is a sample setup in `.zshrc`.

```sh
# Use pathhelper to get system paths plus user paths
# This should fail over to using the system one so there will be some sort of path
if [ -x ~/bin/pathhelper ]; then
    eval $(~/bin/pathhelper -z)
elif [ -x /usr/libexec/path_helper ]; then
    echo "using fallback path building"
    eval $(/usr/libexec/path_helper -s)
    path +=('/opt/homebrew/bin' '~/bin' '/Users/ian/.ops/bin')
    export PATH
fi
```

Paths made as a result of a call in `.zshrc` must be added to the main `PATH` separately. If you put them in a path file
the entry will be rejected as it would not evaluate to a valid path.

Here is the help output for `pathhelper`

```
$ pathhelper -h
Usage: pathhelper [--bash] [--zsh] [--csh] [--verbose] [--init]

Options:
  --bash, -s             get bash format path settings
  --zsh, -z              get zsh format path settings
  --csh, -c              get csh format path settings
  --verbose, -v          display issues as paths evaluated
  --init, -i             check and build user path dirs if necessary
  --help, -h             display this help and exit
```

Note that formats compatible with `csh`, `bash`, and `zsh` are available using flags. If there is no specifier `bash`
format will be used, which also works for `zsh`. There is also an init (`-i`) method which will make `.config/..`
directories if they do not exist. Using the `-v` (verbose) flag will show entries that were rejected because they did
not exist.