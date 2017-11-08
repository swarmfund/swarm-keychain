-- +migrate Up

create table keys (
  id bigserial,
  account_id varchar(64) not null,
  filename varchar(255) not null,
  key varchar(64) not null,
  constraint unique_account_id_filename unique (account_id, filename)
);

-- +migrate Down

drop table keys;