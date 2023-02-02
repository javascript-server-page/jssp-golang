package db

/*
func GenerateObjJsdo(js *engine.VM) *otto.Object {
	cache := NewDBCache()
	obj := js.CreateObjectValue().Object()
	obj.Set("mysql", func(dataSourceName string) otto.Value {
		// panic(js.CreateError(errors.New("dsadsa")))
		return otto.Value{}
		//db, err := cache.GetDB("mysql", dataSourceName)
		//if err != nil {
		//	return *js.CreateError(err)
		//}
		//return *def_db(js, db)
	})
	obj.Set("mssql", func(dataSourceName string) otto.Value {
		db, err := cache.GetDB("mssql", dataSourceName)
		if err != nil {
			return *js.CreateError(err)
		}
		return *def_db(js, db)
	})
	obj.Set("postgres", func(dataSourceName string) otto.Value {
		db, err := cache.GetDB("postgres", dataSourceName)
		if err != nil {
			return *js.CreateError(err)
		}
		return *def_db(js, db)
	})
	return obj
}

func def_db(js *engine.VM, db *sql.DB) *otto.Value {
	val := js.CreateObjectValue()
	obj := val.Object()
	obj.Set("queryOne", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("query", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	obj.Set("exec", func(key *string, val *string) {

	})
	obj.Set("begin", func(call otto.FunctionCall) otto.Value {
		return otto.Value{}
	})
	return val
}

//func build_rows(js *engine.VM, rows *sql.Rows) *otto.Value {
//	keys, err := rows.Columns()
//	if err != nil {
//		return nil, err
//	}
//	arr := make([]map[string]interface{}, 0)
//	if len(keys) == 0 {
//		return rows, nil
//	}
//	tmp := make([]interface{}, len(keys))
//	for i := range tmp {
//		var str sql.NullString
//		tmp[i] = &str
//	}
//	for rows.Next() {
//		err = rows.Scan(tmp...)
//		if err != nil {
//			return rows, err
//		}
//		row := make(map[string]interface{})
//		for i, key := range keys {
//			str := tmp[i].(*sql.NullString)
//			if str.Valid {
//				row[key] = str.String
//			} else {
//				row[key] = ""
//			}
//		}
//		arr = append(arr, row)
//	}
//	err = rows.Err()
//	if err != nil {
//		return arr, err
//	}
//	return arr, rows.Close()
//}
*/
