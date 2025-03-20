CREATE TABLE documents (
    id uuid primary key default gen_random_uuid(),
    "text" text not null,
    source text not null,
    embedding vector(768)
);
