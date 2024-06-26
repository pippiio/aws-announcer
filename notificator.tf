module "notificator" {
  source = "github.com/pippiio/aws-serverless?ref=main"

  name_prefix = local.name_prefix
  default_tags = local.default_tags

  function = {
    "test-frontend" = {
      source = {
        experimental_ecr_cache = true

        type = "container"
        path = "ghcr.io/pippiio/aws-serverless/frontend:latest"
      }
    }
  }
}
