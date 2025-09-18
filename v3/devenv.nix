{ pkgs, lib, config, ... }: {

  languages.javascript = {
    enable = true;
    bun = {
      enable = true;
    };
  };

  env.HOST = "localhost";
  env.PORT = "8080";

  packages = with pkgs; [
    # sqlite
  ];

  processes = {
    resume-wizard = {
      exec = "bun run dev --host ${config.env.HOST} --port ${config.env.PORT}";
      process-compose = {
        working_dir = "${config.env.DEVENV_ROOT}";
        # log_location = "${config.env.DEVENV_ROOT}/logs/resumewizard.log";
        availability = {
          restart = "on_failure";
          max_restarts = 3;
          backoff_seconds = 2;
        };
        readiness_probe = {
          http_get = {
            host = config.env.HOST;
            port = config.env.PORT;
            path = "/";
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


  # scripts = {
  #   check.exec = ''
  #     ruff check
  #   '';
  #   format.exec = ''
  #     ruff format
  #   '';
  # };
  
  # git-hooks.hooks = {
  #   shellcheck.enable = true;
  #   ruff.enable = true;
  #   unit-tests = {
  #     # TODO: enable when we have tests :)
  #     enable = false;
  #     name = "Unit tests";
  #     entry = "python -m unittest discover .";
  #     # types = [ "python" ];
  #     language = "python";
  #
  #     # Set this to false to not pass the changed files
  #     # to the command (default: true):
  #     pass_filenames = false;
  #   };
  # };

  # enterTest = " python -m unittest discover . ";
}
