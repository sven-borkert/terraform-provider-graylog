variable "web_endpoint_uri" {
  type        = string
  description = "Graylog API base (e.g., https://graylog.internal.borkert.net/api)"
}

variable "auth_name" {
  type        = string
  description = "Graylog username (admin for test)"
}

variable "auth_password" {
  type        = string
  description = "Graylog password (admin for test)"
  sensitive   = true
}

variable "x_requested_by" {
  type        = string
  description = "Optional X-Requested-By header"
  default     = "terraform-provider-graylog-e2e"
}

variable "api_version" {
  type        = string
  description = "Graylog API version path segment"
  default     = "v3"
}
