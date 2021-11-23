# Prion

A package manager for vim, designed for use with [pathogen](https://github.com/tpope/vim-pathogen)

### Installation

#### Homebrew
```
brew install Jefferson-Faseler/taps/prion
```

#### Go install
```
go get -u github.com/Jefferson-Faseler/prion
cd $GOPATH/src/github.com/Jefferson-Faseler/prion
go install .
```

### Usage
Prion is designed to work like any other package manager from the cli

#### Adding a package
This will work for either https or ssh remote urls

```
prion install [repo-urls]
```

#### Updating
```
prion update [package names]
```

or
```
prion update --all
```

#### Reinstalling
```
prion reinstall [package names]
```

or
```
prion install [repo-urls]
# and confirm in prompt
```

#### Removing
```
prion rm [package names]
```

#### Listing
```
prion ls
```

#### Editing .vimrc
```
prion config add [configuration string]
```

Ex: `prion config add "set autoindent"`

#### Opening .vimrc
Open your .vimrc in your default editor if set, otherwise it will open it in vim.

```
prion config edit
```


### What it does
At it's core, prion is simplifying the normal pathogen method of cloning into `~/.vim/bundle/` and then over time updating the packages as you go.

Prion also offers an easier way to update your `.vimrc` configuration. In a normal workflow, you are constantly opening, editing, saving, and reloading your vim editor to make and check changes. Prion cuts that cycle into two steps for most cases.


### Overriding defaults
Prion is looking for two env variables you may have set:
- `VIMRC_PATH` (defaults to `$HOME/.vimrc`)
- `VIM_BUNDLE_DIR` (defaults to `$HOME/.vim/bundle`)

To override defaults you can you can export the vars, in-line them, or make a `.prion.env` in your home directory.

## Contributing
Knock yourself out
