# data "graylog_dashboard" "test" {
#   title = "test"
# }

resource "graylog_dashboard" "test" {
  title       = "test"
  description = "test"
}
