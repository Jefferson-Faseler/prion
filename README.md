# Prion

A package manager for vim, designed for use with [pathogen](https://github.com/tpope/vim-pathogen)

### Installation
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
prion install [repo-url]
```

#### Updating
```
prion update [package name]
```

#### Removing
```
prion rm [package name]
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
At it's core, prion, is simplifying the normal pathogen method of cloning into `~/.vim/bundle/` and then over time updating the packages as you go.

Prion also offers an easier way to update your `.vimrc` configuration. In a normal workflow especially when you are trying something new whether it be a package, framework, or language. You are constantly opening, editing, saving, and reloading your vim editor to make and check changes. Prion cuts that cycle into two steps for most cases.


## Contributing
Knock yourself out
