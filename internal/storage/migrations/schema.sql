-- Create table for storing urls --
CREATE TABLE IF NOT EXISTS urls (
  id SERIAL NOT NULL PRIMARY KEY,
  short TEXT NOT NULL,
  long TEXT NOT NULL,
  view_count INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create function for keeping `updated_at` updated with current timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS trigger AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$
    LANGUAGE 'plpgsql';

-- Create trigger to call the function we created to keep updated_at updated
DROP trigger IF EXISTS urls_updated_at ON urls;
CREATE trigger urls_updated_at
    BEFORE UPDATE ON urls FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();

-- Create a unique index across the appropriate columns
CREATE UNIQUE INDEX unique_url_idx
    ON urls (long, short);

-- Add unique constraint to the urls table so that we only have one copy of each deterministic URL
ALTER TABLE urls
    ADD CONSTRAINT unique_url_constraint UNIQUE
    USING INDEX unique_url_idx;

-- Create index on short column, as it is used for search constantly
CREATE index IF NOT EXISTS idx_urls_short
    ON urls(short);
