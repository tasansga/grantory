resource "grantory_grant" "database" {
  request_id = grantory_request.database.id

  payload = jsonencode({
    user     = "alice"
    password = "local-runner"
  })
}
