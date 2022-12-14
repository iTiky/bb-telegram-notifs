-- Events table: send ack added (true for all existing events)
alter table events
    add column send_ack bool      default false,
    add column send_at  timestamp default null;

update events
set send_ack = true,
    send_at  = now();

alter table events
    alter column send_ack set not null;
