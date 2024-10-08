# dirtygit

Do you find yourself context-switching between a bunch of different git repos?

Have you ever accidentally discovered that changes you've made locally have not
been committed or pushed to your git server?

`dirtygit` is a text-mode UI tool to find git repos that have uncommitted files or which have not
been pushed to a remote.

## Source-mode installation

```bash
go install github.com/boyvinall/dirtygit@master
```

## Configuration

Copy [.dirtygit.yml](.dirtygit.yml) to `~/.dirtygit.yml` and edit to your needs.

## Running

```bash
dirtygit [ <directories...> ]
```

If one/more directories are specified as `<directories>`, then this will override the
`scandirs.include` from your config file.

![demo](demo.gif)

## UI

Simple key navigation in the UI as follows:

| Key               | Action                                           |
| ----------------- | ------------------------------------------------ |
| `<up>` / `<down>` | Navigation inside repositories or diff views     |
| `<tab>`           | switch focus between repositories and diff views |
| `e`               | Open selected repo in editor (current vscode)    |
| `s`               | Rescan directories                               |
| `q` / `ctrl-C`    | quit                                             |

Inside the "diff" view, a list of dirty files is shown, with the git status
for both staged changes (`S`) and working directory (`W`).

## Development

```bash
make lint
```

## Future

- Also scan for local changes which have not been pushed
- [x] Allow configuration of editor
- Improve error logging and presentation
- Possibly show git diff
- Allow deletion of files / repositories
- Add tests!

## Pim First

- [x] ignore dir errors
- [x] simple build task
- [x] dir wildcards
- [x] exclude part of paths
- [x] custom edit command

## Pim Future

- [ ] flake
- [ ] refact into src
- [ ] output to json
- [ ] deamon mode
- [ ] menubar https://github.com/fyne-io/systray
- [ ] add regex to exclude
- [ ] if no starting slash, compare from right to left
- [ ] help popup
- [ ] enable global gitignore conf which is already visible in config
