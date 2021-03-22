package dumper

import "regexp"

var (
	//regular expressions used for DBMS recognition based on error message response
	dbmsErrors map[string][]*regexp.Regexp = map[string][]*regexp.Regexp{
		"MySQL": {regexp.MustCompile("SQL syntax.*MySQL"), regexp.MustCompile("Warning.*mysql_.*"), regexp.MustCompile("valid MySQL result"), regexp.MustCompile("MySqlClient\\.")},

		"PostgreSQL": {regexp.MustCompile("PostgreSQL.*ERROR"), regexp.MustCompile("Warning.*\\Wpg_.*"), regexp.MustCompile("valid PostgreSQL result"), regexp.MustCompile("Npgsql\\.")},

		"MSSQL": {regexp.MustCompile("Driver.* SQL[\\-\\_\\ ]*Server"), regexp.MustCompile("OLE DB.* SQL Server"), regexp.MustCompile("(\\W|\\A)SQL Server.*Driver"),
			regexp.MustCompile("Warning.*mssql_.*"), regexp.MustCompile("(\\W|\\A)SQL Server.*[0-9a-fA-F]{8}"), regexp.MustCompile("(?s)Exception.*\\WSystem\\.Data\\.SqlClient\\."),
			regexp.MustCompile("(?s)Exception.*\\WRoadhouse\\.Cms\\."), regexp.MustCompile("(?i)Ole DB"), regexp.MustCompile("(?i)ODBC"),
			regexp.MustCompile("(?i)ADODB"), regexp.MustCompile("(?i)Microsoft VBScript")},

		"Access": {regexp.MustCompile("Microsoft Access Driver"), regexp.MustCompile("JET Database Engine"), regexp.MustCompile("Access Database Engine"), regexp.MustCompile("(?i)OleDb"), regexp.MustCompile("(?i)runtime error")},

		"Oracle": {regexp.MustCompile("\bORA-[0-9][0-9][0-9][0-9]"), regexp.MustCompile("Oracle error"), regexp.MustCompile("Oracle.*Driver"), regexp.MustCompile("Warning.*\\Woci_.*"), regexp.MustCompile("Warning.*\\Wora_.*")},

		"IBM DB2": {regexp.MustCompile("CLI Driver.*DB2"), regexp.MustCompile("DB2 SQL error"), regexp.MustCompile("\bdb2_\\w+\\(")},

		"SQLite": {regexp.MustCompile("SQLite/JDBCDriver"), regexp.MustCompile("SQLite.Exception"), regexp.MustCompile("System.Data.SQLite.SQLiteException"), regexp.MustCompile("Warning.*sqlite_.*"),
			regexp.MustCompile("Warning.*SQLite3::"), regexp.MustCompile("\\[SQLITE_ERROR\\]")},

		"Sybase": {regexp.MustCompile("(?i)Warning.*sybase.*"), regexp.MustCompile("Sybase message"), regexp.MustCompile("Sybase.*Server message.*")},

		"MariaDB": {regexp.MustCompile("(?i)MariaDB server"), regexp.MustCompile("(?i)MariaDB")},
	}
)
