{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    yamllint
    # TODO: Add go
  ];

  shellHook = ''
    echo "Resume Wizard Environment Loaded! 🧙"
  '';
}
