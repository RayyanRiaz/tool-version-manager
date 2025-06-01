{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [ pkgs.go_1_24 ];
  hardeningDisable = [ "fortify" ];
  shellHook = ''
    export GOPATH=$PWD/.gopath
    export GOBIN=$GOPATH/bin
    export PATH=$GOBIN:$PATH
    mkdir -p $GOBIN
  '';
}
