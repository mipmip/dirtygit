{ lib, buildGoModule, go }:
buildGoModule rec {
  pname = "dirtygit";

  version = "1.0";

  src = ./.;

  doCheck = false;

  vendorHash = "sha256-KBu77tQfZjZsAcUatXZj+sHa+5uUNN5PuFaSk1rzIkQ=";

  meta = with lib; {
    description = ''
      Convert aws config and credential files into a single JSON object
    '';
    homepage = "https://github.com/mipmip/dirtygit";
    license = licenses.bsd2;
  };

}
