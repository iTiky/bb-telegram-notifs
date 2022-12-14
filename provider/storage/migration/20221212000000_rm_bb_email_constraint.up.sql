-- User table: remove BitBucket email constraint
alter table users drop constraint users_bb_email_key;
