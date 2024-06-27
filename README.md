# aws-announcer
Terraform module for deploying a AWS serverles notifications center to send messages in multiple communication channels

# aws-announcer
The _aws-announcer is a generic Terraform module within the pippi.io family, maintained by Tech Chapter. The pippi.io modules are build to support common use cases often seen at Tech Chapters clients. They are created with best practices in mind and battle tested at scale. All modules are free and open-source under the Mozilla Public License Version 2.0.

The aws-announcer module is made to provision a AWS serverles notifications center to send messages in multiple communication channels such as, sms, email, slack, teams etc.

## Example usage
```hcl

```

<!-- BEGIN_TF_DOCS -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | >= 1.6.0 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 4.0 |

## Resources

No resources.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_default_tags"></a> [default\_tags](#input\_default\_tags) | A map of default tags, that will be applied to all resources applicable. | `map(string)` | `{}` | no |
| <a name="input_name_prefix"></a> [name\_prefix](#input\_name\_prefix) | A prefix that will be used on all named resources. | `string` | `"pippi-"` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->
