{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    yamllint
    gnumake
    go
  ];

  shellHook = ''
    echo "Resume Wizard Environment Loaded! 🧙"
  '';
}
