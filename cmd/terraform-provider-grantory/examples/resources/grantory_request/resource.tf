resource "grantory_request" "database" {
  host_id = grantory_host.app.host_id

  payload = jsonencode({
    db = "things"
  })

  labels = {
    team = "operations"
  }
}
