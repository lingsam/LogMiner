# LogMiner
无需安装oracle客户端环境，同步是基于LogMiner增量同步

## 配置规则
```
# include-table 和 exclude-table 不能同时配置，两者只能配置一个,如果两个都没配置则 Schema 内表全迁移
# include-table 和 exclude-table 支持正则表达式以及通配符（tab_*/tab*）
include-table = ["T_SALE_H","T_DIST_D_XLX","T_STORE_I_XLX"] 
exclude-table = []
```

## 安装运行
```
# go build
# .\gominerlog.exe -mode incr
# .\gominerlog.exe -mode full
# .\gominerlog.exe -mode all

PS E:\LogMiner> .\gominerlog.exe -mode incr
oracle://H2:hyft@192.168.9.18:1521/hyee
firstSCN: 17494430683543 maxSCN: 281474976710655 LOG_FILE: /u01/app/oracle/oradata/HYDEE/redo_6rd
############增量##############
LogMiner\module\migrate\o2m\incr.go:192
[update "H2"."T_STORE_I_XLX" set "PID" = 1709264778 where "PID" = 1709264776]
[update "H2"."T_STORE_I_XLX" set "PID" = 1709264776 where "PID" = 1709264778]
```

## 原理
```
 增量定时 time.Tick(time.Duration(100) * time.Millisecond)

/* 1.查看日志路径 */
SELECT	* FROM	v$logfile;

BEGIN
	dbms_logmnr.add_logfile (
		logfilename => '/opt/oracle/app/oradata/orcl/redo03.log',
		options => dbms_logmnr. NEW
	) ; 
END ;
BEGIN
	dbms_logmnr.add_logfile (
		logfilename => '/opt/oracle/app/oradata/orcl/redo0.log',
		options => dbms_logmnr.ADDFILE
	) ;
END ; 

/* 2.查看添加的日志 */
SELECT	filename FROM 	V$LOGMNR_LOGS ;

/* 3.开始分析 */
BEGIN
	DBMS_LOGMNR.START_LOGMNR (
		OPTIONS => DBMS_LOGMNR.DICT_FROM_ONLINE_CATALOG
	) ;
END ; 

/* 4.获取结果 */
SELECT TABLE_NAME,USERNAME, SCN,TO_CHAR(timestamp,'yyyy-mm-dd hh:mi:ss am'),SQL_REDO FROM V$logmnr_contents where SCN > 0 AND table_name = 'INFOS' ORDER BY timestamp DESC;

select table_name,sql_redo,timestamp,TO_CHAR(timestamp,'yyyy-mm-dd hh:mi:ss am'),SCN,username,session_info from v$logmnr_contents where table_name='INFOS' ORDER BY timestamp DESC;

select sql_redo,timestamp,username,session_info from v$logmnr_contents;

/* 5.结束 */
BEGIN
 DBMS_LOGMNR.END_LOGMNR;
END ;
```

## 注意事项
```
SELECT value FROM v$parameter WHERE name = 'open_cursors';
ALTER SESSION SET OPEN_CURSORS = 3000;
ALTER SYSTEM SET OPEN_CURSORS = 3000 SCOPE=BOTH;
```
