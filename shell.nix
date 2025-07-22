{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    yamllint
    pandoc
    wkhtmltopdf

    # TODO: Add go
  ];

  shellHook = ''
    echo "Resume Wizard Environment Loaded! 🧙"
  '';
}
