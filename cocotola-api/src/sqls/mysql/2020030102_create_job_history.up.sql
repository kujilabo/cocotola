create table `job_history` (
 `job_status_id` varchar(36) not null
,`job_name` varchar(40) not null
,`job_parameter` text not null
,`status` varchar(20) not null
,`created_at` datetime not null
);
