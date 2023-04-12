resource "forwardnetworks_proxy" "example" {
  network_id           = "your_network_id"
  protocol             = "https"
  host                 = "proxy.example.com"
  port                 = 8080
  username             = "proxyuser"
  password             = "proxypassword"
  disable_cert_checking = false
}
