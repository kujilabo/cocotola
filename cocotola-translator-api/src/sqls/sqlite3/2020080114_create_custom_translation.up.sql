create table `custom_translation` (
 `version` int not null
,`created_at` datetime not null default current_timestamp
,`updated_at` datetime not null default current_timestamp
,`text` varchar(30) not null
,`pos` int not null
,`lang2` varchar(2) not null
,`translated` varchar(100) not null
,`disabled` tinyint(1) not null default 0
,primary key(`text`, `pos`, `lang2`)
);
