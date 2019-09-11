

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
	primary key (id),
	index (parent_id),
	index (name),
	index (sort_order)
)engine innodb charset utf8mb4;

