CREATE TABLE note (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL,
    expires_at_timezone TEXT NOT NULL
);
