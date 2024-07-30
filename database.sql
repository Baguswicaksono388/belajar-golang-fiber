create table sample( id varchar(255) not NULL, name varchar(255) not NULL, primary key (id)) engine=INNODB;

create table users ( id varchar(100) not NULL, password varchar(100) not NULL, name varchar(255) not NULL, created_at timestamp not NULL default current_timestamp, updated_at timestamp not NULL default current_timestamp on update current_timestamp, primary key (id)) engine=INNODB;

alter table users rename column name to first_name;
alter table users add column middle_name varchar(100) NULL after first_name;
alter table users add column last_name varchar(100) NULL after middle_name;

create table user_logs(id int auto_increment, user_id varchar(100) not NULL, action varchar(100) not NULL, created_at timestamp not null default current_timestamp, updated_at timestamp not null default current_timestamp on update current_timestamp, primary key (id)) engine=InnoDB;