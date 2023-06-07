alter table op_role
    add path varchar(255);
comment on column op_role.path is '配置角色默认跳转页(可为空)';