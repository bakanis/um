create sequence um_users_id_seq start 1;

create table um_users (
	id bigint not null default nextval('um_users_id_seq'),
	email_addr varchar(128) not null,
	display_name varchar(64) not null default '',
	status integer not null default 0,
	hash varchar(128) not null default '',
	salt varchar(32) not null default '',
	created_on timestamp with time zone not null default 'now',
	last_login timestamp with time zone not null default 'epoch',

	constraint pk_um_users primary key(id),
	constraint ux_um_users_email_addr unique(email_addr)
);

create index ix_um_users_created_on on um_users(created_on);
create index ix_um_users_last_login on um_users(last_login);

