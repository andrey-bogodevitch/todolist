create table public.users
(
    id         bigserial
        primary key,
    name       text,
    role       text,
    created_at timestamp with time zone,
    login      text,
    password   text,
    deleted_at timestamp with time zone
);


create table public.sessions
(
    id         uuid,
    user_id    bigint,
    created_at timestamp with time zone,
    expired_at timestamp with time zone
);


create table public.tasks
(
    user_id    bigint,
    task       text,
    status     text,
    deleted_at timestamp with time zone,
    id         bigserial
        primary key
);

