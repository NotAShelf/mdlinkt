{
  callPackage,
  gopls,
  go,
}: let
  mainPkg = callPackage ./package.nix {};
in
  mainPkg.overrideAttrs (oa: {
    nativeBuildInputs =
      [
        gopls
        go
      ]
      ++ (oa.nativeBuildInputs or []);
  })
