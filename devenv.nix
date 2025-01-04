{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: let
  pkgu = import inputs.nixpkgs-unstable {system = pkgs.stdenv.system;};
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

  enterTest = ''
    go test ./... -race -coverprofile=coverage.out -covermode=atomic
  '';

  scripts = {
    run-docs = {
      exec = ''
        mkdocs serve
      '';
      description = "Run the documentation server";
    };
    gen-doc-refs = {
      exec = ''
        CURRENT_DIR=$PWD
        cd $CURRENT_DIR/vikunja && gomarkdoc --output ../docs/Vikunja.md
        cd $CURRENT_DIR/utils && gomarkdoc --output ../docs/Utils.md
        cd $CURRENT_DIR
      '';
      description = "Generate the documentation references";
    };
  };
  # languages.go = {
  #   enable = true;
  #   package = pkgs.go;
  # };

  enterShell = ''
    echo
    echo 🦾 Useful project scripts:
    echo 🦾
    ${pkgs.gnused}/bin/sed -e 's| |••|g' -e 's|=| |' <<EOF | ${pkgs.util-linuxMinimal}/bin/column -t | ${pkgs.gnused}/bin/sed -e 's|^|🦾 |' -e 's|••| |g'
    ${lib.generators.toKeyValue {} (lib.mapAttrs (_: value: value.description) config.scripts)}
    EOF
    echo
  '';
}
