data "grantory_hosts" "app" {
  labels = {
    env = "prod"
  }
}
