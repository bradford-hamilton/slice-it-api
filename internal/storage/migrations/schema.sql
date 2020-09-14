-- Create table for storing urls --
CREATE TABLE IF NOT EXISTS urls (
  id SERIAL NOT NULL PRIMARY KEY,
  short TEXT NOT NULL,
  long TEXT NOT NULL,
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
