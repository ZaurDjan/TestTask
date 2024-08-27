ALTER TABLE sessions
ADD COLUMN if not exists ip_address text not null;