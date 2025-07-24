{
  description = "";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs";
  };
  outputs =
    { self, nixpkgs, ... }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
      ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          termpet = pkgs.buildGoModule {
            pname = "termpet";
            version = "0.0.2";
            src = ./.;
            rev = "v0.0.2-a";
            
            vendorHash = "0rwhqavn771gkw94sjhwf6jfcapngbvqfyhrban9xz43v6hvqap4";

            buildInputs = [
              pkgs.go 
              pkgs.neo-cowsay
            ];

          };
          default = termpet;
        }
      );

      devShells = forAllSystems (system:
      let 
        pkgs = nixpkgsFor.${system};
      in 
      {
        default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.neo-cowsay
          ];
        };
      }
      );
    };
}