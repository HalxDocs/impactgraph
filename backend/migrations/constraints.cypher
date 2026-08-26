CREATE CONSTRAINT developer_id IF NOT EXISTS FOR (d:Developer) REQUIRE d.id IS UNIQUE;
CREATE CONSTRAINT repository_id IF NOT EXISTS FOR (r:Repository) REQUIRE r.id IS UNIQUE;
CREATE CONSTRAINT package_id IF NOT EXISTS FOR (p:Package) REQUIRE p.id IS UNIQUE;
CREATE CONSTRAINT service_id IF NOT EXISTS FOR (s:Service) REQUIRE s.id IS UNIQUE;
CREATE CONSTRAINT team_id IF NOT EXISTS FOR (t:Team) REQUIRE t.id IS UNIQUE;
CREATE CONSTRAINT organization_id IF NOT EXISTS FOR (o:Organization) REQUIRE o.id IS UNIQUE;
CREATE INDEX package_name IF NOT EXISTS FOR (p:Package) ON (p.name);
CREATE INDEX service_name IF NOT EXISTS FOR (s:Service) ON (s.name);
CREATE INDEX repository_name IF NOT EXISTS FOR (r:Repository) ON (r.name);
