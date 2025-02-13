create table to_do_list.role (
	id			bigserial primary key,
	created_at	timestamp default localtimestamp not null,
	valid_from	timestamp default localtimestamp not null,
	valid_to	timestamp default '9999-12-31 23:59:59.999'::timestamp not null, 
	code		int2	unique not null,
	name		varchar(260) unique not null
);

comment on table to_do_list.role is 'Перечень доступных для пользователей ролей';

comment on column to_do_list.role.id is 'Идентификатор роли';
comment on column to_do_list.role.created_at is 'Дата и время создания роли';
comment on column to_do_list.role.valid_from is 'Дата и время вступления роли в силу';
comment on column to_do_list.role.valid_to is 'Дата и время окончания действия роли';
comment on column to_do_list.role.code is 'Код роли';
comment on column to_do_list.role.name is 'Наименование роли';

create table to_do_list.user (
	id			bigserial primary key,
	created_at	timestamp default localtimestamp not null,
	role_code	int2 not null,
	name		varchar(100) not null,
	login		varchar(100) unique not null,
	email		varchar(100) unique not null,
	password	varchar(100) not null,
	foreign key (role_code) references to_do_list.role (code)
);

comment on table to_do_list.user is 'Перечень пользователей приложения';

comment on column to_do_list.user.id is 'Идентификатор пользователя';
comment on column to_do_list.user.created_at is 'Дата и время создания пользователя';
comment on column to_do_list.user.role_code is 'Код роли пользователя';
comment on column to_do_list.user.name is 'Имя пользователя';
comment on column to_do_list.user.login is 'Логин пользователя';
comment on column to_do_list.user.email is 'Электронная почта пользователя';
comment on column to_do_list.user.password is 'Пароль пользователя';

create table to_do_list.note_type (
	id			bigserial primary key,
	created_at	timestamp default localtimestamp not null,
	creator_id	int8 not null,
	name		varchar(100) not null,
	foreign key (creator_id	) references to_do_list.user (id)
);

comment on table to_do_list.note_type is 'Перечень типов заметок';

comment on column to_do_list.note_type.id is			'Идентификатор типа заметки';
comment on column to_do_list.note_type.created_at is	'Дата и время создания типа заметки';
comment on column to_do_list.note_type.creator_id is	'Создатель типа заметки';
comment on column to_do_list.note_type.name	 is			'Наименование типа заметки';

create table to_do_list.note (
	id				bigserial primary key,
	created_at		timestamp default localtimestamp not null,
	creator_id		int8 not null,
	type			int8 not null,
	content			text not null,
	is_completed	bool not null default false,
	foreign key (creator_id) references to_do_list.user (id),
	foreign key (type) references to_do_list.note_type (id)
)

comment on table to_do_list.note is 'Перечень заметок'

comment on column to_do_list.note.id is 'Идентификатор заметки';
comment on column to_do_list.note.created_at is 'Дата и время создания заметки';
comment on column to_do_list.note.creator_id is 'Идентификатор автора заметки';
comment on column to_do_list.note.type is 'Тип заметки';
comment on column to_do_list.note.content is 'Содержание заметки';
comment on column to_do_list.note.is_completed is 'Признак актуальности заметки';