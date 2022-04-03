# pathhelper

Helper for building a path for MacOS. Most of my effort in writing this has been tied to making sure I did things in the
proper order and did checking of paths. The code is iterative with no goroutines or channels since the paths need to be
read in a predictable and repeatable order.

For any `zsh` specific examples in this README there are straightforward `bash` and `csh` equivalents.

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

Looking at several of the above entries, it can be seen that some installers such as the installer for Golang add
to this list. This likely makes for many fewer complaints that "go can't be found". Homebrew takes a different approach
and adds a call with an exact path to the end of `.zprofile` that prepends to the path.

The convention is that there will be one path per line in the main `paths` and `manpaths` and in files in the `paths.d`
and `manpaths.d` directories. More than one line in files is permitted, and lines started with `#` are ignored.

When looking for paths in files, each path encountered is checked to ensure it exists and is readable or these
non-existent files files are ignonred. These invalid items are silently rejected but can be seen if the `-v` flag is
used.

Here is a sample setup in `.zshrc` with an attempt at a failsafe if pathhelper fails.

```sh
# Use pathhelper to get system paths plus user paths
if [ -x ~/bin/pathhelper ]; then
    eval $(~/bin/pathhelper -z)
elif [ -x /usr/libexec/path_helper ]; then
    echo "using fallback path building"
    eval $(/usr/libexec/path_helper -s)
fi

# Just in case use system path_helper if things went wrong
if [ "$PATH" = "" ]; then
    if [ -x /usr/libexec/path_helper ]; then
        eval `/usr/libexec/path_helper -s`
    fi
fi
```

## Notes

Paths made as a result of a call in `.zshrc` must be added to the main `PATH` separately. If you put them in a path file
the entry will be rejected as it would not evaluate to a valid path. Here is an example.

```sh
export GOPATH="$(go env GOPATH)"
path+=("$GOPATH/bin")
export PATH
```

Note that in both `bash` and `zsh` there are several files that are run at the start of a new terminal session. For `bash`
this is `bash_profile` and for `zsh` this is `.zshprofile`. Things like homebrew install a line to `~/.zprofile` that
runs `eval $(/opt/homebrew/bin/brew shellenv)`. This adds homebrew dirs to the path.

If desired, user directories can be put first using the `-u` flag. This could be to allow overriding of system commands.
This is essentially what homebrew does with its invocation (`eval $(/opt/homebrew/bin/brew shellenv)`).

## Usage

Here is the help output for `pathhelper`

```
$ pathhelper -h
Usage: pathhelper [--bash] [--zsh] [--csh] [--verbose] [--trace] [--init] [--user-first]

Options:
  --bash, -s             get bash format path settings
  --zsh, -z              get zsh format path settings
  --csh, -c              get csh format path settings
  --verbose, -v          display issues as paths evaluated
  --trace, -t            display very detailed activity
  --init, -i             check and build user path dirs if necessary
  --user-first, -u       put user directories first
  --help, -h             display this help and exit
```

Note that formats compatible with `csh`, `bash`, and `zsh` are available using flags. If there is no specifier `bash`
format will be used, which also works for `zsh`. There is also an init (`-i`) method which will make `.config/..`
directories if they do not exist. Using the `-v` (verbose) flag will show entries that were rejected because they did
not exist.

Here is a sample of the extra output to stderr when using the `-v` flag. It is reasonably straightforward to use any
errors in the output in combination with the previous lines to track down invalid paths.

```
$ pathhelper -z -v
INFO evaluating /etc/paths
INFO checking /usr/local/bin
INFO checking /usr/bin
INFO checking /bin
INFO checking /usr/sbin
INFO checking /sbin
INFO evaluating /etc/paths.d
INFO checking /Library/Apple/usr/bin
INFO checking /opt/X11/bin
INFO checking /Library/TeX/texbin
INFO checking /Users/ian/.dotnet/tools
ERROR stat /Users/ian/.dotnet/tools: no such file or directory
INFO checking /usr/local/bin
INFO checking /usr/local/go/bin
INFO evaluationg /Users/ian/.config/pathhelper/paths.d
INFO checking /opt/homebrew/bin
INFO checking /Users/ian/bin
INFO checking /Users/ian/.ops/bin
INFO evaluating /etc/manpaths
INFO checking /usr/share/man
INFO checking /usr/local/share/man
INFO evaluating /etc/manpaths.d
INFO checking /opt/X11/share/man
INFO checking /Library/TeX/Distributions/.DefaultTeX/Contents/Man
INFO evaluationg /Users/ian/.config/pathhelper/manpaths.d
```

The original `path_helper` is about 182 lines of code. This project's code is

```
$ gocloc cmd README.md
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                               5             82             43            356
Markdown                         1             25              0            119
-------------------------------------------------------------------------------
TOTAL                            6            107             43            475
-------------------------------------------------------------------------------
```