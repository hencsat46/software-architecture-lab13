CREATE TABLE notifications (
    id uuid PRIMARY KEY,
    user_id text NOT NULL,
    message text NOT NULL,
    created_at timestamptz NOT NULL
);
