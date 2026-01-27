terraform {
  required_providers {
    grantory = {
      source  = "tasansga/grantory"
      version = "{{ .ProviderVersion }}"
    }
  }
}

provider "grantory" {
  server = "{{ .ServerURL }}"
}

data "grantory_requests" "ungranted" {
{{- if .RequestLabels }}
  labels = {
{{- range .RequestLabels }}
    {{ .Key }} = "{{ .Value }}"
{{- end }}
  }
{{- end }}
}

data "grantory_registers" "integration" {
{{- if .RegisterLabels }}
  labels = {
{{- range .RegisterLabels }}
    {{ .Key }} = "{{ .Value }}"
{{- end }}
  }
{{- end }}
}

locals {
  outstanding_request = data.grantory_requests.ungranted.requests[0]
}

resource "grantory_grant" "integration" {
  request_id = local.outstanding_request.request_id

  payload = jsonencode({
    granted = true
  })
}

output "register_snapshot" {
  value = data.grantory_registers.integration.registers[0].register_id
}
