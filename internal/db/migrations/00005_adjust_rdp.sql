-- +goose Up
-- +goose StatementBegin
CREATE TABLE router_dns_provider_new (
  traefik_id INTEGER NOT NULL,
  provider_id INTEGER NOT NULL,
  router_name TEXT NOT NULL,
  FOREIGN KEY (traefik_id) REFERENCES traefik (id) ON DELETE CASCADE,
  FOREIGN KEY (provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE,
  UNIQUE (traefik_id, router_name, provider_id)
);

INSERT INTO
  router_dns_provider_new (traefik_id, provider_id, router_name)
SELECT
  traefik_id,
  provider_id,
  router_name
FROM
  router_dns_provider;

DROP TABLE router_dns_provider;

ALTER TABLE router_dns_provider_new
RENAME TO router_dns_provider;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
CREATE TABLE router_dns_provider_old (
  traefik_id INTEGER NOT NULL,
  provider_id INTEGER NOT NULL,
  router_name TEXT NOT NULL,
  FOREIGN KEY (traefik_id) REFERENCES traefik (id) ON DELETE CASCADE,
  FOREIGN KEY (provider_id) REFERENCES dns_providers (id) ON DELETE CASCADE,
  UNIQUE (traefik_id, router_name)
);

INSERT
OR IGNORE INTO router_dns_provider_old (traefik_id, provider_id, router_name)
SELECT
  traefik_id,
  provider_id,
  router_name
FROM
  router_dns_provider;

DROP TABLE router_dns_provider;

ALTER TABLE router_dns_provider_old
RENAME TO router_dns_provider;

-- +goose StatementEnd
