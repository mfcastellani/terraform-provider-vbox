provider "vbox" {
  user = "Nome do usuário"
  token = "Token do usuário"
}

resource "vbox_disk" "meu_disco" {
  size = 1024
  name = "awesome_disk"
}

resource "vbox_instance" "minha_imagem" {
  name = "awesome_image"
  osname = "Ubuntu"
  diskname = "${vbox_disk.meu_disco.name}"
}

