# Index Set Templates: Data Sources and Usage

## What’s implemented
- `graylog_index_set_template` (data source): fetches a template by title (built-ins like “7 Days Hot”, “14 Days Hot”, “30 Days Hot”) and returns its `id`, `title`, and `description`.
- `graylog_index_set_templates` (data source): lists all built-in templates; useful for discovery and debugging.
- `graylog_index_set_template` (resource): create/manage custom templates by supplying `index_set_config` JSON.

## Example usage (built-in 30 Days Hot)
```hcl
data "graylog_index_set_template" "hot30" {
  title = "30 Days Hot"
}

resource "graylog_index_set" "tf_e2e" {
  index_set_template_id = data.graylog_index_set_template.hot30.id
  # ...other fields...
}

output "builtin_templates" {
  value = data.graylog_index_set_templates.all.templates
}
```

## Notes
- Graylog 7 requires a template ID when creating index sets; using a data source avoids hardcoding cluster-specific IDs.
- Built-in template IDs (on this cluster): `7 Days Hot` = `692067a9371785137c9b2b73`, `14 Days Hot` = `692067a9371785137c9b2b74`, `30 Days Hot` = `692067a9371785137c9b2b75`.
- Custom templates are managed via the `graylog_index_set_template` resource (`index_set_config` must be valid JSON matching Graylog’s IndexSetTemplateConfig shape).
