create table if not exists "group" (
    id uuid primary key,
    specialization_code varchar not null,
    group_number varchar not null,

    unique (specialization_code, group_number)
);
