# Release Checklist

## Pre-Release Steps

### 1. Repository Setup
- [ ] Rename repository to `terraform-provider-azure-github-runners`
- [ ] Ensure repository is public
- [ ] Verify repository follows naming convention

### 2. GPG Key Setup
- [ ] Generate GPG key for signing releases
- [ ] Export public key: `gpg --armor --export "your-email@example.com"`
- [ ] Add public key to Terraform Registry (User Settings > Signing Keys)
- [ ] Add GPG private key to GitHub Secrets as `GPG_PRIVATE_KEY`
- [ ] Add GPG passphrase to GitHub Secrets as `PASSPHRASE`

### 3. GitHub Secrets Configuration
Add the following secrets to your repository:
- [ ] `GPG_PRIVATE_KEY` - ASCII-armored GPG private key
- [ ] `PASSPHRASE` - GPG key passphrase
- [ ] `GITHUB_TOKEN` - Personal Access Token with `public_repo` scope

### 4. Code Quality
- [ ] Run tests: `make test`
- [ ] Run linter: `make lint`
- [ ] Format code: `make fmt`
- [ ] Generate documentation: `make docs`
- [ ] Verify all features work correctly

### 5. Version Management
- [ ] Update version in `main.go` if needed
- [ ] Update CHANGELOG.md with release notes
- [ ] Ensure semantic versioning (e.g., v1.0.0)

## Release Process

### 1. Create Release
- [ ] Create and push version tag: `git tag v1.0.0 && git push origin v1.0.0`
- [ ] Verify GitHub Actions workflow runs successfully
- [ ] Check that release is created on GitHub with all required assets

### 2. Verify Release Assets
The release should contain:
- [ ] Binary files for multiple OS/arch combinations
- [ ] `terraform-provider-azure-github-runners_1.0.0_manifest.json`
- [ ] `terraform-provider-azure-github-runners_1.0.0_SHA256SUMS`
- [ ] `terraform-provider-azure-github-runners_1.0.0_SHA256SUMS.sig`

### 3. Publish to Terraform Registry
- [ ] Sign in to Terraform Registry with GitHub account
- [ ] Go to Publish > Provider
- [ ] Select organization and repository
- [ ] Verify provider appears in registry
- [ ] Test provider installation: `terraform init`

## Post-Release

### 1. Documentation
- [ ] Update README.md with registry installation instructions
- [ ] Verify all documentation links work
- [ ] Test example configurations

### 2. Community
- [ ] Announce release on relevant channels
- [ ] Update any related documentation or tutorials
- [ ] Monitor for issues and feedback

## Troubleshooting

### Common Issues
1. **Repository naming**: Must be `terraform-provider-{NAME}`
2. **GPG signing**: Ensure key is properly configured
3. **GitHub Actions**: Check workflow permissions and secrets
4. **Binary naming**: Must match `terraform-provider-{NAME}_v{VERSION}`

### Support
- Terraform Registry: terraform-registry@hashicorp.com
- GitHub Issues: Create issue in repository
