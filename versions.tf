terraform {
  required_version = ">= 0.13"
  required_providers {
    github-runners = {
      source  = "yurii-kysil/azure-github-runners"
      version = "~> 1.0"
    }
  }
}
