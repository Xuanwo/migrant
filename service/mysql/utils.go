package mysql

const schemaMigrationTable = `CREATE TABLE IF NOT EXISTS schema_migrations (
  id         VARCHAR(255) NOT NULL,
  type       VARCHAR(255) NOT NULL,
  applied_at BIGINT       NOT NULL  DEFAULT '0',

  PRIMARY KEY (id)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
`