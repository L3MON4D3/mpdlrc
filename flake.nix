{
  description = "mpd-lrc";

  # Nixpkgs / NixOS version to use.
  inputs.nixpkgs.url = "nixpkgs/nixos-24.11";

  outputs = { self, nixpkgs }:
    let
      # Generate a user-friendly version number.
      version = "1.0.0";

      # System types to support.
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      # Helper function to generate an attrset '{ x86_64-linux = f "x86_64-linux"; ... }'.
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      # Nixpkgs instantiated for supported system types.
      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      # Provide some binary packages for selected system types.
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          mpd-lrc = pkgs.buildGoModule {
            pname = "mpd-lrc";
            inherit version;
            # In 'nix develop', we don't need a copy of the source tree
            # in the Nix store.
            src = ./.;

            preBuild = ''
              export GOWORK=off
            '';

            vendorHash = "sha256-OE1U3dbsp9v05aNQex3o73tY5CsvKttPsmjWHPdEeqs=";
          };
        });

      # Add dependencies that are only needed for development
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go gopls gotools go-tools gomod2nix.outputs.packages.${system}.default ];
          };
        });

      # The default package for 'nix build'. This makes sense if the
      # flake provides only one package or there is a clear "main"
      # package.
      defaultPackage = forAllSystems (system: self.packages.${system}.mpd-lrc);
    };
}
