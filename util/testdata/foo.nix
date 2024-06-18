{ modulesPath, lib, pkgs, ...}:
{
  system.stateVersion = "24.05";
  imports = [
    "${toString modulesPath}/nixpkgs"
    ./bar.nix
  ];
}
