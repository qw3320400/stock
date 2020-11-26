CREATE DATABASE `stock`

-- 按年分表
CREATE TABLE `stock_k_data_%d` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `create_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `code` varchar(64) NOT NULL COMMENT '股票代码',
  `time_cst` datetime NOT NULL COMMENT '时间',
  `frequency` varchar(64) NOT NULL COMMENT '周期',
  `adjust_flag` varchar(64) NOT NULL COMMENT '复权',
  `open` varchar(64) NOT NULL COMMENT '开盘价',
  `high` varchar(64) NOT NULL COMMENT '最高价',
  `low` varchar(64) NOT NULL COMMENT '最低价',
  `close` varchar(64) NOT NULL COMMENT '收盘价',
  `preclose` varchar(64) NOT NULL COMMENT '收盘价',
  `volume` varchar(64) NOT NULL COMMENT '成交量',
  `amount` varchar(64) NOT NULL COMMENT '成交额',
  `turn` varchar(64) NOT NULL COMMENT '换手率x100%',
  `trade_status` varchar(64) NOT NULL COMMENT '交易状态',
  `pct_chg` varchar(64) NOT NULL COMMENT '涨跌幅x100%',
  `pe_ttm` varchar(64) NOT NULL COMMENT '滚动市盈率x100%',
  `pb_mrq` varchar(64) NOT NULL COMMENT '市净率x100%',
  `ps_ttm` varchar(64) NOT NULL COMMENT '滚动市销率x100%',
  `pcf_ncf_ttm` varchar(64) NOT NULL COMMENT '滚动市现率x100%',
  `is_st` varchar(64) NOT NULL COMMENT '是否ST股',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_stock_k_data_01` (`code`,`time_cst`,`frequency`,`adjust_flag`),
  KEY `idx_stock_k_data_02` (`code`),
  KEY `idx_stock_k_data_03` (`time_cst`),
  KEY `idx_stock_k_data_04` (`frequency`),
  KEY `idx_stock_k_data_05` (`adjust_flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='股票k线数据';

CREATE TABLE `stock_all_code` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `create_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `code` varchar(64) NOT NULL COMMENT '股票代码',
  `name` varchar(256) NOT NULL COMMENT '股票名称',
  `industry` varchar(256) NOT NULL COMMENT '所属行业',
  `industry_classification` varchar(256) NOT NULL COMMENT '所属行业类别',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_stock_all_code_01` (`code`),
  KEY `idx_stock_all_code_02` (`industry`),
  KEY `idx_stock_all_code_03` (`industry_classification`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='股票代码';

CREATE TABLE `stock_trade_date` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `create_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `date_cst` datetime NOT NULL COMMENT '日期',
  `is_trading_day` varchar(64) NOT NULL COMMENT '是否交易日',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_stock_trade_date_01` (`date_cst`),
  KEY `idx_stock_trade_date_02` (`is_trading_day`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='股票交易日';

CREATE TABLE `stock_stategy_result` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `create_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `code` varchar(64) NOT NULL COMMENT '股票代码',
  `tag` varchar(64) NOT NULL COMMENT '策略标记',
  `start_time_cst` datetime NOT NULL COMMENT '开始时间',
  `end_time_cst` datetime NOT NULL COMMENT '结束时间',
  `anual_return_rate` varchar(64) NOT NULL COMMENT '年化收益率',
  `draw_down` varchar(64) NOT NULL COMMENT '最大回撤',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_stock_stategy_result_01` (`code`,`tag`,`start_time_cst`,`end_time_cst`),
  KEY `idx_stock_stategy_result_02` (`code`),
  KEY `idx_stock_stategy_result_03` (`tag`),
  KEY `idx_stock_stategy_result_04` (`start_time_cst`),
  KEY `idx_stock_stategy_result_05` (`end_time_cst`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='策略主要指标';

CREATE TABLE `stock_stategy_result_data` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `create_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time_utc` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  `stock_stategy_result_id` bigint(20) NOT NULL COMMENT '策略结果id',
  `code` varchar(64) NOT NULL COMMENT '股票代码',
  `tag` varchar(64) NOT NULL COMMENT '策略标记',
  `time_cst` datetime NOT NULL COMMENT '时间',
  `value` varchar(64) NOT NULL COMMENT '值',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_stock_stategy_data_01` (`stock_stategy_result_id`,`time_cst`),
  KEY `idx_stock_stategy_data_02` (`stock_stategy_result_id`),
  KEY `idx_stock_stategy_data_03` (`code`),
  KEY `idx_stock_stategy_data_04` (`tag`),
  KEY `idx_stock_stategy_data_05` (`time_cst`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='策略数据';
