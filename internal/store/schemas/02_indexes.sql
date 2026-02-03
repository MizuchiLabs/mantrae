CREATE INDEX idx_http_middlewares_profile_name ON http_middlewares (profile_id, name);

CREATE INDEX idx_http_routers_profile_name ON http_routers (profile_id, name);

CREATE INDEX idx_http_servers_transports_profile_name ON http_servers_transports (profile_id, name);

CREATE INDEX idx_http_services_profile_name ON http_services (profile_id, name);

CREATE INDEX idx_tcp_middlewares_profile_name ON tcp_middlewares (profile_id, name);

CREATE INDEX idx_tcp_routers_profile_name ON tcp_routers (profile_id, name);

CREATE INDEX idx_tcp_servers_transports_profile_name ON tcp_servers_transports (profile_id, name);

CREATE INDEX idx_tcp_services_profile_name ON tcp_services (profile_id, name);

CREATE INDEX idx_udp_routers_profile_name ON udp_routers (profile_id, name);

CREATE INDEX idx_udp_services_profile_name ON udp_services (profile_id, name);
