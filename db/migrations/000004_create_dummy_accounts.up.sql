INSERT INTO users (id, email, password_hash, role, created_at)
VALUES (
  gen_random_uuid(),
  'admin@example.com',
  '$2a$12$H2PlXUqGzIR8x5nBSGS2g.YNrk6bO85A2rv10lWsCg8l8Mid53bza',
  'admin',
  NOW()
);

INSERT INTO users (id, email, password_hash, role, created_at)
VALUES (
  gen_random_uuid(),
  'user@example.com',
  '$2a$12$H2PlXUqGzIR8x5nBSGS2g.YNrk6bO85A2rv10lWsCg8l8Mid53bza',
  'user',
  NOW()
);
