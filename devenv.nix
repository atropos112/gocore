{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: let
  pkgu = import inputs.nixpkgs-unstable {system = pkgs.stdenv.system;};
  # writeShellScript here is identity to cause treesitter to format bash scripts correctly.
  writeShellScript = name: script: script;
  helpScript = writeShellScript "help" ''
    echo
    echo 🦾 Useful project scripts:
    echo 🦾
    ${pkgs.gnused}/bin/sed -e 's| |••|g' -e 's|=| |' <<EOF | ${pkgs.util-linuxMinimal}/bin/column -t | ${pkgs.gnused}/bin/sed -e 's|^|🦾 |' -e 's|••| |g'
    ${lib.generators.toKeyValue {} (lib.mapAttrs (_: value: value.description) config.scripts)}
    EOF
    echo
  '';
in {
  packages = with pkgu; [
    natscli
    nats-top
    nats-server
    gomarkdoc
  ];

  pre-commit = {
    hooks = {
      check-merge-conflicts.enable = true;
      check-added-large-files.enable = true;
      editorconfig-checker.enable = true;
      govet.enable = true;
      gofmt.enable = true;
      gen-doc-refs = {
        enable = true;
        entry = ''gen-doc-refs '';
      };
    };
  };

  enterTest = writeShellScript "test" ''
    go test ./... -race -coverprofile=coverage.out -covermode=atomic
  '';

  scripts = {
    run-docs = {
      exec = writeShellScript "run-docs" ''
        mkdocs serve
      '';
      description = "Run the documentation server";
    };
    gen-doc-refs = {
      exec = writeShellScript "gen-doc-refs" ''
        CURRENT_DIR=$PWD
        cd $CURRENT_DIR/vikunja && gomarkdoc --output ../docs/Vikunja.md
        cd $CURRENT_DIR/utils && gomarkdoc --output ../docs/Utils.md
        cd $CURRENT_DIR
      '';
      description = "Generate the documentation references";
    };
    help = {
      exec = helpScript;
      description = "Show this help message";
    };
  };
  # languages.go = {
  #   enable = true;
  #   package = pkgs.go;
  # };

  enterShell = helpScript;
}
