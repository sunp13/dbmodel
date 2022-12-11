package models

import "github.com/sunp13/dbtool"

type assetPool struct{}

// GetList
func (m *assetPool) GetList() (res []map[string]interface{},err error){
    sql := `
    select * from dps_asset.asset_pool
    where is_deleted = 0
    `
    res, err = dbtool.D.QuerySQL(sql, nil)
    return
}

// GetListByID
func (m *assetPool) GetListByID(id string) (res map[string]interface{}, err error){
    sql := `
    select * from dps_asset.asset_pool
    where pool_id = ?
    `
    params := []interface{}{
        id,
    }

    var result []map[string]interface{}
    result, err = dbtool.D.QuerySQL(sql,params)
    if len(result) > 0{
        res = result[0]
    }
    return
}

// AddList
func (m *assetPool) AddList(poolName,poolDesc,poolLocation,isDeleted string) (res int64, err error) {
	sql := `
	insert into dps_asset.asset_pool(
	pool_name,
	pool_desc,
	pool_location,
	is_deleted
	) values (?,?,?,?)
	`
	params := []interface{}{
	poolName,
	poolDesc,
	poolLocation,
	isDeleted,
	}

	res, err = dbtool.D.UpdateSQL(sql, params)
	return
}

// ModifyList
func (m *assetPool) ModifyList(poolId ,poolName,poolDesc,poolLocation,isDeleted string) (res int64, err error) {
	sql := `
	update dps_asset.asset_pool set
    pool_name,
	pool_desc,
	pool_location,
	is_deleted
	where pool_id = ?
	`
	params := []interface{}{
	poolName,
	poolDesc,
	poolLocation,
	isDeleted,
	}
	res, err = dbtool.D.UpdateSQL(sql, params)
	return
}

// DeleteList
func (m *assetPool) DeleteList(id string) (res int64, err error) {
	sql := `
	update dps_asset.asset_pool set
	is_deleted = 1
	where pool_id = ?
	`
	params := []interface{}{
		id,
	}
	res, err = dbtool.D.UpdateSQL(sql, params)
	return
}
