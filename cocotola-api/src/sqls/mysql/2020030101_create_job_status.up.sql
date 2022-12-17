create table `job_status` (
 `id` varchar(26) not null
,`job_name` varchar(40) not null
,`job_parameter` text not null
,`concurrency_key` varchar(40)
,`expiration_datetime` datetime not null
,`created_at` datetime not null default current_timestamp
,primary key(`id`)
,unique(`job_name`, `concurrency_key`)
);
