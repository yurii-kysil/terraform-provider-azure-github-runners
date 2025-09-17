variable "github_token" {
  description = "GitHub personal access token"
  type        = string
  sensitive   = true
  default     = null
}

variable "github_app_id" {
  description = "GitHub App ID"
  type        = number
  default     = null
}

variable "github_app_installation_id" {
  description = "GitHub App installation ID"
  type        = number
  default     = null
}

variable "github_app_pem_file" {
  description = "Path to GitHub App private key PEM file"
  type        = string
  default     = null
}

variable "organization" {
  description = "GitHub organization name"
  type        = string
  default     = "my-org"
}
