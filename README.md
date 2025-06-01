# Tool Version Manager (TVM)

A personal tool for quick management of different tools. A tool can be anything like binaries, editors, softwares like node etc.
I have assumed that a tool-version-manager can be aptly managed by:

- ToolDiscovery
    - GetAllLocalVersions
    - GetAllRemoteVersions
    - GetLatestRemoteVersion
- ToolLinker
    - LinkTool
    - UnlinkTool
    - GetLinkInfo
- ToolInstaller
    - InstallToolForVersion
- ToolComparer
    - CompareVersions

For now, I have defined one implementation i.e. script-driven tvm. More can be defined if needed


## TODOs:

- refactor the logic out of `cmd` files
- `table` viewer should be enhanced with more options
- general code cleanups

---

Please note that this is more of my personal project, and the very first one in golang, assisted by copilot. So don't judge :D, and use at your own risk.

# Dev instructions:

1. nix-shell
2. code --no-sandbox $(pwd)


