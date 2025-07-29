{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    gnumake
    go
    jq
    yamllint
    sqlc
  ];

  shellHook = ''
    echo "Resume Wizard Environment Loaded! 🧙"
  '';
}
