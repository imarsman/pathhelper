# pathhelper

Helper for building a path for MacOS. Most of my effort in writing this has been tied to making sure I did things in the
proper order and did checking of paths. The code is iterative with no goroutines or channels since the paths need to be
read in a predictable and repeatable order.

For any `zsh` specific examples in this README there are straightforward `bash` and `csh` equivalents.

I wrote this to better understand how the `PATH` and `MANPATH` variables are set in MacOS. I ended up learning a lot
about how shell start up and how they define variables for a login session. The `/usr/libexec/path_helper` binary runs
on my laptop in about 6 msec. pathhelper takes about 12 msec to run on my laptop. Given that the load time for my `zsh`
environment is 500 msec this is at the moment acceptable. `path_helper` is [written in
C](https://opensource.apple.com/source/shell_cmds/shell_cmds-162/path_helper/path_helper.c.auto.html) and does not from
what I can see do any validation of the paths found. Not checking for the existence of the directories pointed bo by the
path elements found could partly explain its faster execution.

Linux distributions do not store system paths as lists of files in `/etc/paths`, `/etc/manpaths` or in `/etc/paths.d/`,
or `/etc/manpaths.d`, so this would not really work to set the path in a `.zshrc` file. It would work as a user path
setting tool but that is limited. In theory the system parts would gracefully fail in Linux but the user paths could be
helpful to add to the `PATH` variable. The issue would be that one would have to mess around with the system path and
that is not recommented. A modification of this would just give a list of additional path elements and leave it up to
the user to append them to the main `PATH`.

In addition to the work that `/usr/libexec/path_helper` does on MacOS `pathhelper` also looks in `~/.config/pathhelper/`
for files in `paths.d` and `manpaths.d`. 

Here is the list of entries for `/etc/paths.d`

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

`/opt/homebrew/bin:/opt/homebrew/sbin:/usr/local/bin:...`

The convention is that there will be one path per line in the main `paths` and `manpaths` and in files in the `paths.d`
and `manpaths.d` directories. More than one line in files is permitted, and lines started with `#` are ignored.

## What is processed

`path_helper` does not verify the paths it processes. It does escape characters, including `"`, `'`, `$`, and `\`.
`pathhelper` aslso escapes these characters as well as `!`. I initially set `pathhelper` to by default check for the
existence of all paths that it comes across. I have switched that to being done if the flag, `-V` is used. If you want
to check paths and get a report on what paths do not exist, use `pathhelper -z -V -v` to verify all paths and report
things found. Skipping the verification of each path saves about a millisecond.

Here is an error for a path entry put in `/etc/paths.d` for a tool no longer installed. That also strikes me as
something that should not be in `/etc/paths.d`. This invalid path was highlighted by running `pathhelper -z -V -v`.

`ERROR stat /Users/ian/.dotnet/tools: no such file or directory`

## Use in .zshenv

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

Paths made as a result of a call, such as with go when calling `go env GOPATH` in `.zshrc` must be added to the main
`PATH` separately. If you put them in a path file the entry will be rejected as it would not evaluate to a valid path.
Here is an example.

```sh
export GOPATH="$(go env GOPATH)"
path+=("$GOPATH/bin")
export PATH
```

Note that in both `bash` and `zsh` there are several files that are run at the start of a new terminal session. For
`bash` this is `bash_profile` and for `zsh` this is `.zshprofile`. Homebrew installs a line to `~/.zprofile` that runs
`eval $(/opt/homebrew/bin/brew shellenv)`. This adds homebrew dirs to the path.

If desired, user directories can be put first using the `-u` flag. This could be to allow overriding of system commands.
This is essentially what homebrew does with its invocation (`eval $(/opt/homebrew/bin/brew shellenv)`).

## Usage

Here is the help output for `pathhelper`

```
 pathhelper -h
Usage: pathhelper [--bash] [--zsh] [--csh] [--verify] 
                  [--verbose] [--trace] [--init] [--user-first]

Options:
  --bash, -s             get bash format path settings
  --zsh, -z              get zsh format path settings
  --csh, -c              get csh format path settings
  --verify, -V           verify paths' existence before adding
  --trace, -t            show paths evaluated and time do evaluate
  --init, -i             check and build user path dirs if necessary
  --verbose, -v          display issues as paths evaluated
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
INFO --------------------------------------------------------------------------------
INFO processing /etc/paths
INFO checking /etc/paths
INFO checking /usr/local/bin
INFO checking /usr/bin
INFO checking /bin
INFO checking /usr/sbin
INFO checking /sbin
INFO processing /etc/paths took 706.334µs
INFO --------------------------------------------------------------------------------
INFO processing /etc/paths.d
INFO checking /etc/paths.d/100-rvictl
INFO checking /Library/Apple/usr/bin
INFO checking /etc/paths.d/40-XQuartz
INFO checking /opt/X11/bin
INFO checking /etc/paths.d/TeX
INFO checking /Library/TeX/texbin
INFO checking /etc/paths.d/dotnet-cli-tools
INFO checking /Users/ian/.dotnet/tools
ERROR stat /Users/ian/.dotnet/tools: no such file or directory
INFO checking /etc/paths.d/exiftool
INFO checking /usr/local/bin
INFO checking /etc/paths.d/go
INFO checking /usr/local/go/bin
INFO processing /etc/paths.d took 3.160583ms
INFO --------------------------------------------------------------------------------
INFO processing /Users/ian/.config/pathhelper/paths.d
INFO checking /Users/ian/.config/pathhelper/paths.d/homebrew
ERROR skipping line in homebrew "# Note that brew sets its paths in ~/.zprofile or ~/.bash_profile"
INFO checking /Users/ian/.config/pathhelper/paths.d/localbin
INFO checking /Users/ian/bin
INFO checking /Users/ian/.config/pathhelper/paths.d/ops
INFO checking /Users/ian/.ops/bin
INFO processing /Users/ian/.config/pathhelper/paths.d took 1.101458ms
INFO --------------------------------------------------------------------------------
INFO processing /etc/manpaths
INFO checking /etc/manpaths
INFO checking /usr/share/man
INFO checking /usr/local/share/man
INFO processing /etc/manpaths took 147.375µs
INFO --------------------------------------------------------------------------------
INFO processing /etc/manpaths.d
INFO checking /etc/manpaths.d/40-XQuartz
INFO checking /opt/X11/share/man
INFO checking /etc/manpaths.d/TeX
INFO checking /Library/TeX/Distributions/.DefaultTeX/Contents/Man
INFO processing /etc/manpaths.d took 980µs
INFO --------------------------------------------------------------------------------
INFO processing /Users/ian/.config/pathhelper/manpaths.d
INFO processing /Users/ian/.config/pathhelper/manpaths.d took 270.75µs
INFO --------------------------------------------------------------------------------
```

The original `path_helper` is about 182 lines of code. This project's code is

```
$ gocloc cmd README.md
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                               5             84             44            381
Markdown                         1             26              0            122
-------------------------------------------------------------------------------
TOTAL                            6            110             44            503
-------------------------------------------------------------------------------
```