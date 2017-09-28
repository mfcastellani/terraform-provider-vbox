provider "vbox" {
  user = "Nome do usuário"
  token = "Token do usuário"
}

resource "vbox_disk" "meu_disco" {
  size = 1024
  name = "awesome_disk"
}