CREATE TABLE orders (
    id uuid PRIMARY KEY,
    user_id text NOT NULL,
    product text NOT NULL,
    quantity integer NOT NULL CHECK (quantity > 0),
    created_at timestamptz NOT NULL
);
