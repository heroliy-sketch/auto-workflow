# auto-workflow

`auto-workflow` is a small Go HTTP service with the same deployment shape as the
reference project: GitHub Actions builds a container image, pushes it to
Amazon ECR, and deploys it to Amazon EKS with Helm.

## Local development

```bash
go mod download
make test
make build
make run
```

The service listens on port `8080` by default. The health endpoint is:

```text
GET /health
```

Set `AUTO_WORKFLOW_CONFIG` to load a different YAML configuration file.

## Deployment configuration

The workflow in `.github/workflows/deploy.yml` deploys only the `devin`, `uat`,
and `prod` environments. It runs for pushes to those branches and can also be
started manually. It expects these GitHub Actions values:

- Secrets: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
- Variables: `AWS_REGION`, `EKS_CLUSTER_NAME`, `ECR_REPOSITORY`, optional
  `KUBE_NAMESPACE`, `INGRESS_ENABLED`, and `INGRESS_HOST`

The AWS principal must be able to push to ECR and update the target EKS
cluster. Configure the repository values before pushing to a deployment branch.
