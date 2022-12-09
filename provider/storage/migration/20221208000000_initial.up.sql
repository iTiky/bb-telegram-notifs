-- Users table
create table users
(
    id         bigserial primary key,
    tg_id      bigint    not null unique,
    tg_login   varchar(255) unique,
    tg_chat_id bigint    not null unique,
    bb_email   varchar(255) unique,
    active     boolean            default false,
    created_at timestamp not null default now()
);
create index users_tg_id ON users (bb_email);
create index users_bb_email ON users (tg_id);

insert into users (tg_id, tg_login, tg_chat_id, bb_email, active)
values (208585440, 'Tikkirey', 208585440, 'mpkornilov@dev.vtb.ru', true);

-- Repositories table
create table repos
(
    id         bigserial primary key,
    project    varchar(128) not null,
    name       varchar(128) not null,
    created_at timestamp    not null default now(),
    unique (project, name)
);

-- User subscriptions table
create table subscriptions
(
    id         bigserial primary key,
    user_id    bigserial   not null references users (id),
    repo_id    bigserial   not null references repos (id),
    type       varchar(50) not null,
    created_at timestamp   not null default now(),
    updated_at timestamp   not null default now(),
    unique (user_id, repo_id)
);
create index subscriptions_user_id ON subscriptions (user_id);
create index subscriptions_repo_id ON subscriptions (repo_id);

-- Key-value storage table
create table kvs
(
    id         varchar(128) primary key,
    value      varchar(256) not null,
    updated_at timestamp    not null default now()
);

-- Events table
create table events
(
    id                   bigserial primary key,
    hash                 varchar(256) not null unique,
    type                 varchar(50)  not null,
    recipient_tg_id      bigint       not null,
    recipient_tg_chat_id bigint       not null,
    sender_name          varchar(255) not null,
    repo_project         varchar(128) not null,
    repo_name            varchar(128) not null,
    pr_id                bigint       not null,
    pr_title             varchar(255),
    pr_url               varchar(512),
    created_at           timestamp    not null
);