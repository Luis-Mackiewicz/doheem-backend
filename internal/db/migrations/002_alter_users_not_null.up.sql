UPDATE users SET phone = '' WHERE phone IS NULL;
UPDATE users SET document = '' WHERE document IS NULL;
UPDATE users SET birth_date = '1970-01-01' WHERE birth_date IS NULL;
UPDATE users SET cep = '' WHERE cep IS NULL;

ALTER TABLE users
  ALTER COLUMN phone      SET NOT NULL,
  ALTER COLUMN document   SET NOT NULL,
  ALTER COLUMN birth_date SET NOT NULL,
  ALTER COLUMN cep        SET NOT NULL;
