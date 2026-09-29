{
  description = "A terminal application that simulates running Claude Code. It doesn't run anything.";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "run-nothing";
          version = "2.1.139";
          src = ./.;
          vendorHash = null;
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/run-nothing";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls ];
        };
      }
    );
}
