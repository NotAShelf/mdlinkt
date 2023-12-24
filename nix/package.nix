{
  buildGoModule,
  lib,
  ...
}: let
  inherit (lib) hasSuffix;

  pname = "mdlinkt";
  version = "0.1.0";
in
  buildGoModule {
    inherit pname version;

    src = builtins.filterSource (path: type: type != "directory" || hasSuffix path != ".nix") ../.;
    vendorHash = null;
    ldflags = ["-s" "-w"];

    doCheck = false;

    meta = {
      description = "A CLI tool for checking for dead links in a markdown file";
      license = lib.licenses.gpl3Only;
      mainProgram = pname;
      maintainers = with lib.maintainers; [NotAShelf];
    };
  }
