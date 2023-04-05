{
  description = "connection testing utility for snowplow destinations";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
            config.allowUnsupportedSystem = true;
          }; in {
            devShell = import ./shell.nix { inherit pkgs; };
            defaultPackage = pkgs.buildGoModule {
              pname = "conntest";
              version = self.shortRev or "${self.lastModifiedDate}-dirty";
              src = self;
              vendorSha256 = "NQGsnuy3SdDvVxlSaOdOXFjekdP29oEPBQKnbUndTRc=";
            };
          }
    );
}
