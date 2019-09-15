

create database projectA charset utf8mb4;
use projectA;



create user bin@localhost identified by '123456';
grant all privileges on projectA.*to bin@localhost;
flush privileges;




create table if not exists a_categories(
	id int unsigned auto_increment,
	parent_id int unsigned,	//子分类
	name varchar(255),	//产品名字
	logo varchar(255),	//图片
	description varchar(255),	//产品描述
	sort_order int,		//排序
	meta_title varchar(255),	//分类标题
	meta_keywords varchar(255),	//分类关键字
	meta_description varchar(255),	//分类相关描述
	create_at timestamp null,	//添加时间
	updated_at timestamp null,	//更新时间
	deleted_at timestamp null,	//删除时间
	primary key (id),
	index (parent_id),
	index (name),
	index (sort_order),
	index (create_at),
	index (updated_at),
	index (deleted_at)
)engine innodb charset utf8mb4;



 
 

alter table a_categories add column create_at timestamp,add column updated_at timestamp,add column deleted_at timestamp;

insert into a_categories (id, name, parent_id) values (1, '未分类', 0);
insert into a_categories (id, name, parent_id) values (2, '图书', 0);
insert into a_categories (id, name, parent_id) values (3, '数码产品', 0);
insert into a_categories (id, name, parent_id) values (4, '纸质书', 2);
insert into a_categories (id, name, parent_id) values (5, '电子书', 2);
insert into a_categories (id, name, parent_id) values (6, '笔记本', 3);
insert into a_categories (id, name, parent_id) values (7, '平板', 3);
insert into a_categories (id, name, parent_id) values (8, '一体机', 3);
insert into a_categories (id, name, parent_id) values (9, '联想', 6);
insert into a_categories (id, name, parent_id) values (10, '华硕', 6);
insert into a_categories (id, name, parent_id) values (11, '戴尔', 6);
insert into a_categories (id, name, parent_id) values (12, '惠普', 6);

