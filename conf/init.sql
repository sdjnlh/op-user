DROP TABLE IF EXISTS op_user;
CREATE TABLE op_user
(
  id        BIGINT PRIMARY KEY,
  username  VARCHAR(20),
  password  VARCHAR(500),
  nickname  VARCHAR(30),
  email     VARCHAR(50),
  mcc       VARCHAR(5),
  mobile    VARCHAR(15),
  language  VARCHAR(10),
  avatar_id VARCHAR(33),
  role_ids  BIGINT[],
  groups    VARCHAR[],
  org_id    BIGINT,
  ext       JSONB,
  status    INT                      DEFAULT (1),
  crt       TIMESTAMP WITH TIME ZONE DEFAULT now(),
  lut       TIMESTAMP WITH TIME ZONE DEFAULT now(),
  dtd       BOOLEAN                  DEFAULT FALSE
);

--//@formatter:off
COMMENT ON TABLE op_user IS '用户表';
COMMENT ON COLUMN op_user.id IS '主键';
COMMENT ON COLUMN op_user.username IS '用户名（字母、数字、下划线、横线），可用于登录';
COMMENT ON COLUMN op_user.password IS '密码';
COMMENT ON COLUMN op_user.nickname IS '昵称';
COMMENT ON COLUMN op_user.email IS '邮箱';
COMMENT ON COLUMN op_user.mcc IS '手机号码的国家区号';
COMMENT ON COLUMN op_user.mobile IS '手机号码';
COMMENT ON COLUMN op_user.language IS '用户语言';
COMMENT ON COLUMN op_user.avatar_id IS '头像id';
COMMENT ON COLUMN op_user.status IS '状态';
COMMENT ON COLUMN op_user.crt IS '创建时间';
COMMENT ON COLUMN op_user.lut IS '最后更新时间';
--//@formatter:on

--//@formatter:on
DROP TABLE IF EXISTS op_role;
CREATE TABLE op_role
(
  id          BIGINT PRIMARY KEY,
  name        VARCHAR(20),
  code        VARCHAR(20),
  permissions VARCHAR[],
  description VARCHAR(100),
  crt         TIMESTAMP WITH TIME ZONE DEFAULT now(),
  lut         TIMESTAMP WITH TIME ZONE DEFAULT now(),
  dtd         BOOLEAN                  DEFAULT FALSE
);
--//@formatter:off
COMMENT ON TABLE op_role IS '角色表';
COMMENT ON COLUMN op_role.id IS '主键';
COMMENT ON COLUMN op_role.name IS '名称';
COMMENT ON COLUMN op_role.name IS '名称';
COMMENT ON COLUMN op_role.permissions is '权限';
COMMENT ON COLUMN op_role.description IS '简述';
COMMENT ON COLUMN op_role.crt IS '创建时间';
COMMENT ON COLUMN op_role.lut IS '最后更新时间';

--//@formatter:on
DROP TABLE IF EXISTS op_permission;
CREATE TABLE op_permission
(
  id          VARCHAR(30) PRIMARY KEY,
  name        VARCHAR(30),
  parent_id   VARCHAR(30),
  path        VARCHAR[],
  description VARCHAR(100),
  crt         TIMESTAMP WITH TIME ZONE DEFAULT now(),
  lut         TIMESTAMP WITH TIME ZONE DEFAULT now(),
  dtd         BOOLEAN                  DEFAULT FALSE
);
--//@formatter:off
COMMENT ON TABLE op_permission IS '权限表';
COMMENT ON COLUMN op_permission.id IS '主键';
COMMENT ON COLUMN op_permission.name IS '名称';
COMMENT ON COLUMN op_permission.parent_id IS '上级权限ID';
COMMENT ON COLUMN op_permission.path IS '父子关系路径';
COMMENT ON COLUMN op_permission.description IS '简述';
COMMENT ON COLUMN op_permission.crt IS '创建时间';
COMMENT ON COLUMN op_permission.lut IS '最后更新时间';


DROP TABLE IF EXISTS op_group;
CREATE TABLE op_group
(
  id          VARCHAR(30) PRIMARY KEY,
  name        VARCHAR(30),
  tags        VARCHAR[],
  description VARCHAR(500),
  crt         TIMESTAMP WITH TIME ZONE DEFAULT now(),
  lut         TIMESTAMP WITH TIME ZONE DEFAULT now(),
  dtd         BOOLEAN                  DEFAULT FALSE
);
--//@formatter:off
COMMENT ON TABLE op_group IS '工作组表';
COMMENT ON COLUMN op_group.id IS '主键id';
COMMENT ON COLUMN op_group.name IS '工作组名称';
COMMENT ON COLUMN op_group.description IS '工作组详细信息';
--//@formatter:on