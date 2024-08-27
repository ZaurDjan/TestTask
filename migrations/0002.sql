ALTER TABLE users
    ALTER COLUMN id TYPE text USING id::text;


ALTER TABLE sessions
    ALTER COLUMN uid TYPE text USING uid::text;


ALTER TABLE assets
    ALTER COLUMN uid TYPE text USING uid::text;

UPDATE users
SET password_hash = '2bb80d537b1da3e38bd30361aa855686bde0eacd7162fef6a25fe97bf527a25b'
WHERE login = 'alice';