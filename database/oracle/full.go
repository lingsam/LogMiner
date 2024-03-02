package oracle

import (
	"fmt"
	"gominerlog/common"

	"github.com/shopspring/decimal"
	"github.com/thinkeridea/go-extend/exstrings"
)

func (o *Oracle) GetOracleCurrentSnapshotSCN() (uint64, error) {
	// 获取当前 SCN 号
	_, res, err := Query(o.Ctx, o.OracleDB, "select min(current_scn) CURRENT_SCN from gv$database")
	var globalSCN uint64
	if err != nil {
		return globalSCN, err
	}
	globalSCN, err = common.StrconvUintBitSize(res[0]["CURRENT_SCN"], 64)
	if err != nil {
		return globalSCN, fmt.Errorf("get oracle current snapshot scn %s utils.StrconvUintBitSize failed: %v", res[0]["CURRENT_SCN"], err)
	}
	return globalSCN, nil
}

// 获取表字段名以及行数据 -> 用于 FULL/ALL
func (o *Oracle) GetOracleTableRowsData(querySQL string, insertBatchSize int) ([]string, []string, error) {
	var (
		err          error
		rowsResult   []string
		rowsTMP      []string
		batchResults []string
		cols         []string
	)
	rows, err := o.OracleDB.QueryContext(o.Ctx, querySQL)
	if err != nil {
		return []string{}, batchResults, err
	}
	defer rows.Close()

	tmpCols, err := rows.Columns()
	if err != nil {
		return cols, batchResults, err
	}

	// 字段名关键字反引号处理
	for _, col := range tmpCols {
		cols = append(cols, common.StringsBuilder("`", col, "`"))
	}

	// 用于判断字段值是数字还是字符
	var (
		columnNames   []string
		columnTypes   []string
		databaseTypes []string
	)
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return cols, batchResults, err
	}

	for _, ct := range colTypes {
		columnNames = append(columnNames, ct.Name())
		// 数据库字段类型 DatabaseTypeName() 映射 go 类型 ScanType()
		columnTypes = append(columnTypes, ct.ScanType().String())
		databaseTypes = append(databaseTypes, ct.DatabaseTypeName())
	}

	// 数据 Scan
	columns := len(cols)
	rawResult := make([][]byte, columns)
	dest := make([]interface{}, columns)
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	// 表行数读取
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return cols, batchResults, err
		}

		for i, raw := range rawResult {
			// 注意 Oracle/Mysql NULL VS 空字符串区别
			// Oracle 空字符串与 NULL 归于一类，统一 NULL 处理 （is null 可以查询 NULL 以及空字符串值，空字符串查询无法查询到空字符串值）
			// Mysql 空字符串与 NULL 非一类，NULL 是 NULL，空字符串是空字符串（is null 只查询 NULL 值，空字符串查询只查询到空字符串值）
			// 按照 Oracle 特性来，转换同步统一转换成 NULL 即可，但需要注意业务逻辑中空字符串得写入，需要变更
			// Oracle/Mysql 对于 'NULL' 统一字符 NULL 处理，查询出来转成 NULL,所以需要判断处理
			if raw == nil {
				rowsResult = append(rowsResult, fmt.Sprintf("%v", `NULL`))
			} else if string(raw) == "" {
				rowsResult = append(rowsResult, fmt.Sprintf("%v", `NULL`))
			} else {
				switch columnTypes[i] {
				case "int64":
					r, err := common.StrconvIntBitSize(string(raw), 64)
					if err != nil {
						return cols, batchResults, fmt.Errorf("column [%s] strconv failed, %v", columnNames[i], err)
					}
					rowsResult = append(rowsResult, fmt.Sprintf("%v", r))
				case "uint64":
					r, err := common.StrconvUintBitSize(string(raw), 64)
					if err != nil {
						return cols, batchResults, fmt.Errorf("column [%s] strconv failed, %v", columnNames[i], err)
					}
					rowsResult = append(rowsResult, fmt.Sprintf("%v", r))
				case "float32":
					r, err := common.StrconvFloatBitSize(string(raw), 32)
					if err != nil {
						return cols, batchResults, fmt.Errorf("column [%s] strconv failed, %v", columnNames[i], err)
					}
					rowsResult = append(rowsResult, fmt.Sprintf("%v", r))
				case "float64":
					r, err := common.StrconvFloatBitSize(string(raw), 64)
					if err != nil {
						return cols, batchResults, fmt.Errorf("column [%s] strconv failed, %v", columnNames[i], err)
					}
					rowsResult = append(rowsResult, fmt.Sprintf("%v", r))
				case "rune":
					r, err := common.StrconvRune(string(raw))
					if err != nil {
						return cols, batchResults, fmt.Errorf("column [%s] strconv failed, %v", columnNames[i], err)
					}
					rowsResult = append(rowsResult, fmt.Sprintf("%v", r))
				case "godror.Number":
					r, err := decimal.NewFromString(string(raw))
					if err != nil {
						return cols, rowsResult, fmt.Errorf("column [%s] NewFromString strconv failed, %v", columnNames[i], err)
					}
					rowsResult = append(rowsResult, fmt.Sprintf("%v", r.String()))
				default:
					// 特殊字符
					rowsResult = append(rowsResult, fmt.Sprintf("'%v'", common.SpecialLettersUsingMySQL(raw)))
				}
			}
		}

		rowsTMP = append(rowsTMP, common.StringsBuilder("(", exstrings.Join(rowsResult, ","), ")"))

		// 数组清空
		rowsResult = rowsResult[0:0]

		// batch 批次
		if len(rowsTMP) == insertBatchSize {
			batchResults = append(batchResults, exstrings.Join(rowsTMP, ","))
			fmt.Println(rowsTMP)

			// 数组清空
			rowsTMP = rowsTMP[0:0]
		}
	}

	if err = rows.Err(); err != nil {
		return cols, batchResults, err
	}

	// 非 batch 批次
	if len(rowsTMP) > 0 {
		batchResults = append(batchResults, exstrings.Join(rowsTMP, ","))
	}

	return cols, batchResults, nil
}
