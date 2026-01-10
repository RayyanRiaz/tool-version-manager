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

## Configuration

TVM uses a YAML configuration file to define tools and their management scripts.

### Config File Location

The configuration file path is resolved in the following order of priority:

1. **CLI flag**: `--config` or `-c`
   ```bash
   tvm --config /path/to/my-tools.yaml list local rg
   ```

2. **Environment variable**: `TVM_CONFIG`
   ```bash
   export TVM_CONFIG=/path/to/my-tools.yaml
   tvm list local rg
   ```

3. **Default**: `tools.yaml` in the current working directory

## TODOs:

- refactor the logic out of `cmd` files
- `table` viewer should be enhanced with more options
- general code cleanups

---

Please note that this is more of my personal project, and the very first one in golang, assisted by copilot. So don't judge :D, and use at your own risk.

# Dev instructions:

1. nix-shell
2. code --no-sandbox $(pwd)


