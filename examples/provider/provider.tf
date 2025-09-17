# Configuration-based authentication with Personal Access Token
provider "azure-github-runners" {
  token        = "ghp_xxxxxxxxxxxxxxxxxxxx"
  organization = "my-org"
  base_url     = "https://api.github.com"
}

# Configuration-based authentication with GitHub App
provider "azure-github-runners" {
  organization = "my-org"

  app_auth {
    id              = 123456
    installation_id = 789012
    pem_file        = "/path/to/private-key.pem"
  }
}
