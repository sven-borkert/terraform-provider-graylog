# Claude Code Instructions

## Refactoring Process for Provider Components

When refactoring a resource or data source for Graylog 7 compatibility, follow this process:

### 1. Gather Information
- Identify the component name (e.g., `dashboard`, `stream`, `event_definition`)
- Read the existing documentation in `docs/resources/todo/` or `docs/data-sources/todo/`
- Read the related API documentation in `docs/api-docs/` (e.g., `dashboards.json`, `views.json`)
- Read the current source code in `graylog/resource/<component>/` or `graylog/datasource/<component>/`
- Read the client code in `graylog/client/<component>/`

### 2. Refactor the Component
- Update the schema in `data_source.go` or `resource.go` if needed
- Update the read/create/update/delete logic to work with Graylog 7 API
- Key changes for Graylog 7:
  - List endpoints return paginated responses with `elements` key instead of the resource name
  - Single resource GET may use `/views/{id}` instead of `/dashboards/{id}` for dashboards
  - Create requests may need entity wrapping via `util.WrapEntityForCreation()`
  - Update requests may need computed field removal via `util.RemoveComputedFields()`

### 3. Test with Local Build
- Build the provider: `make build`
- Add test configuration to `examples/graylog7-e2e/main.tf`
- Test with local provider: `cd examples/graylog7-e2e && ../../bin/terraform-dev plan`
- Apply changes: `../../bin/terraform-dev apply`
- Verify outputs and check for errors
- Repeat until fully working without errors

### 4. Update Documentation
- Create new documentation file in `docs/resources/` or `docs/data-sources/`
- Include:
  - Description of the resource/data source
  - Link to source code
  - Example usage with realistic examples
  - Argument reference with descriptions
  - Attributes reference with descriptions
- Remove the old documentation from the `todo/` subfolder

### 5. Request Manual Testing
- Ask the user to manually verify the changes work as expected
- Only commit after user confirmation

### 6. Commit Changes
- Stage all related files
- Create a descriptive commit message
- Push to remote

## Project Structure

- `graylog/resource/` - Resource implementations
- `graylog/datasource/` - Data source implementations
- `graylog/client/` - API client implementations
- `docs/api-docs/` - Graylog API documentation (JSON)
- `docs/resources/` - Resource documentation (tested)
- `docs/resources/todo/` - Resource documentation (needs refactoring)
- `docs/data-sources/` - Data source documentation (tested)
- `docs/data-sources/todo/` - Data source documentation (needs refactoring)
- `examples/graylog7-e2e/` - End-to-end test configuration

## Local Development

```bash
# Build the provider
make build

# Test with local provider (from examples/graylog7-e2e/)
../../bin/terraform-dev plan
../../bin/terraform-dev apply
../../bin/terraform-dev destroy

# Run linter
make lint

# Run tests
make test
```

## Credentials

Graylog credentials for testing are in `examples/graylog7-e2e/graylog.auto.tfvars` (gitignored).
