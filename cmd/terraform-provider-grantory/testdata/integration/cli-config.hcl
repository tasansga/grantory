provider_installation {
  dev_overrides {
    "tasansga/grantory" = {{ .DevOverridesDir | quote }}
  }
  direct {
    exclude = ["tasansga/grantory"]
  }
}
