#创建数据库
drop database if exists FishMoney;
create database FishMoney;
use FishMoney;

#创建用户表
create table t_user(
	userId integer not null auto_increment,
	name char(32) not null,
	password char(48) not null,
	type integer not null,
	createTime timestamp not null default CURRENT_TIMESTAMP,
	modifyTime timestamp not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP, 
	primary key( userId )
)engine=innodb default charset=utf8 auto_increment = 10001;

alter table t_user add index nameIndex(name,password);

#创建类目表
create table t_category(
	categoryId integer not null auto_increment,	
	userId integer not null,
	name char(32) not null,
	remark varchar(128) not null,
	createTime timestamp not null default CURRENT_TIMESTAMP,
	modifyTime timestamp not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP, 
	primary key( categoryId )
)engine=innodb default charset=utf8 auto_increment = 10001;

alter table t_category add index userIdIndex(userId);

#创建银行卡表
create table t_card(
	cardId integer not null auto_increment,	
	userId integer not null,
	name char(32) not null,
	bank char(32) not null,
	card varchar(32) not null,
	money integer not null,
	remark varchar(128) not null,
	createTime timestamp not null default CURRENT_TIMESTAMP,
	modifyTime timestamp not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP, 
	primary key( cardId )
)engine=innodb default charset=utf8 auto_increment = 10001;

alter table t_card add index userIdIndex(userId);

#创建记账表
create table t_account(
	accountId integer not null auto_increment,
	userId integer not null,
	name char(32) not null,
	money integer not null,
	remark varchar(128) not null,
	categoryId integer not null,
	cardId integer not null,
	type integer not null,
	createTime timestamp not null default CURRENT_TIMESTAMP,
	modifyTime timestamp not null default CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP, 
	primary key( accountId )
)engine=innodb default charset=utf8 auto_increment = 10001;

alter table t_account add index userIdIndex(userId);

#创建session表
create table ci_sessions (
  session_id char(40) default '0' not null,
  ip_address char(45) default '0' not null,
  user_agent char(120) not null,
  last_activity int(10) unsigned default 0 not null,
  user_data text default '' not null,
  primary key (session_id),
  key `last_activity_idx` (`last_activity`)
);

#创建初始数据
insert into t_user(userId,name,password,type) values
(10001,"fish",SHA1("123"),0);

insert into t_category(categoryId,userId,name,remark) values
(10001,10001,"日常收支",''),
(10002,10001,"衣着服装",''),
(10003,10001,"理财投资",''),
(10004,10001,"薪酬工资",'');

insert into t_card(cardId,userId,name,bank,card,money,remark) values
(10001,10001,'工资卡',"农业银行卡",'',0,''),
(10002,10001,'消费卡',"工商银行卡",'',0,''),
(10003,10001,'理财卡',"工商银行卡",'',0,'');

insert into t_account(accountId,userId,name,money,remark,categoryId,cardId,type) values
(10001,10001,"日常支出",100,'',10001,10002,1),
(10002,10001,"日常收入",100,'',10001,10002,1);

#显示一下所有数据
select * from t_user;
select * from t_category;
select * from t_card;
select * from t_account;