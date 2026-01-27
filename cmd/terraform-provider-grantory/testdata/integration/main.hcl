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

resource "grantory_host" "integration" {
  labels = {
    env = "integration"
  }
}

resource "grantory_request" "integration" {
  host_id    = grantory_host.integration.host_id
  payload = jsonencode({
    request = "integration"
  })
{{- if .RequestLabels }}
  labels = {
{{- range .RequestLabels }}
    {{ .Key }} = "{{ .Value }}"
{{- end }}
  }
{{- end }}
}

resource "grantory_register" "integration" {
  host_id     = grantory_host.integration.host_id
  payload = jsonencode({
    source = "{{ .RegisterDataSource }}"
  })
{{- if .RegisterLabels }}
  labels = {
{{- range .RegisterLabels }}
    {{ .Key }} = "{{ .Value }}"
{{- end }}
  }
{{- end }}
}
