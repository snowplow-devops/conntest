{ pkgs ? import <nixpkgs> {}}:
let

  inherit (pkgs) go golangci-lint gopls gore gotools go-tools usql cobra-cli gzip;
in pkgs.mkShell {
  buildInputs = [ go golangci-lint gopls gore gotools go-tools usql cobra-cli gzip];
}
