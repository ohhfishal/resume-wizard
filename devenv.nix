{ pkgs, inputs, lib, config, ... }:
let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in
{
  languages.go = {
    enable = true;
    package = pkgs-unstable.go;
  };

  env.HOST = "localhost";
  env.PORT = "8080";

  # packages = [
  #   pkgs.gnumake
  #   pkgs.jq
  #   pkgs.yamllint
  # ];

  git-hooks.hooks = {
    # Shell
    shellcheck.enable = true;

    # Golang
    govet.enable = true;
    gotest.enable = true;
    gofmt.enable = true;
  };

  tasks = {
    "go:templ" = {
      exec = "go tool templ generate";
      execIfModified = [
        "*.templ"
      ];
      cwd = "./templates";
    };
  };

  processes = {
    templ = {
      exec = "go tool templ generate -watch -lazy";
    };
    resume-wizard = {
      exec = "go run . serve";
      process-compose = {
        working_dir = "${config.env.DEVENV_ROOT}";
        log_location = "${config.env.DEVENV_ROOT}/dev/logs/wizard.log";
        availability = {
          restart = "on_failure";
          max_restarts = 3;
          backoff_seconds = 2;
        };
        readiness_probe = {
          http_get = {
            host = config.env.HOST;
            port = config.env.PORT;
            path = "/health";
          };
          initial_delay_seconds = 5;
          period_seconds = 10;
          timeout_seconds = 3;
          success_threshold = 1;
          failure_threshold = 3;
        };
      };
    };
  };
}
