resource "grantory_host" "app" {
  host_id = "app01"
  labels = {
    env = "prod"
  }
}
