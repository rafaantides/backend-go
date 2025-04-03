env "local" {
  url = "postgres://admin:admin@localhost:5432/rafa?sslmode=disable"
  dev = "docker://postgres/15/dev"
}
