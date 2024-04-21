// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package chat_query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"GoChatServer/dal/model/chat_model"
)

func newMessageUser(db *gorm.DB, opts ...gen.DOOption) messageUser {
	_messageUser := messageUser{}

	_messageUser.messageUserDo.UseDB(db, opts...)
	_messageUser.messageUserDo.UseModel(&chat_model.MessageUser{})

	tableName := _messageUser.messageUserDo.TableName()
	_messageUser.ALL = field.NewAsterisk(tableName)
	_messageUser.ID = field.NewInt64(tableName, "id")
	_messageUser.MessageID = field.NewInt64(tableName, "message_id")
	_messageUser.Receiver = field.NewInt64(tableName, "receiver")
	_messageUser.IsRead = field.NewInt32(tableName, "is_read")
	_messageUser.CreatedAt = field.NewTime(tableName, "created_at")
	_messageUser.UpdatedAt = field.NewTime(tableName, "updated_at")
	_messageUser.DeletedAt = field.NewField(tableName, "deleted_at")

	_messageUser.fillFieldMap()

	return _messageUser
}

// messageUser 消息-用户
type messageUser struct {
	messageUserDo

	ALL       field.Asterisk
	ID        field.Int64 // 自增
	MessageID field.Int64 // 消息ID
	Receiver  field.Int64 // 消息接收人
	IsRead    field.Int32 // 是否已读：0-未读；1-已读；
	CreatedAt field.Time  // 创建时间
	UpdatedAt field.Time  // 更新时间
	DeletedAt field.Field // 删除时间

	fieldMap map[string]field.Expr
}

func (m messageUser) Table(newTableName string) *messageUser {
	m.messageUserDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m messageUser) As(alias string) *messageUser {
	m.messageUserDo.DO = *(m.messageUserDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *messageUser) updateTableName(table string) *messageUser {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt64(table, "id")
	m.MessageID = field.NewInt64(table, "message_id")
	m.Receiver = field.NewInt64(table, "receiver")
	m.IsRead = field.NewInt32(table, "is_read")
	m.CreatedAt = field.NewTime(table, "created_at")
	m.UpdatedAt = field.NewTime(table, "updated_at")
	m.DeletedAt = field.NewField(table, "deleted_at")

	m.fillFieldMap()

	return m
}

func (m *messageUser) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *messageUser) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 7)
	m.fieldMap["id"] = m.ID
	m.fieldMap["message_id"] = m.MessageID
	m.fieldMap["receiver"] = m.Receiver
	m.fieldMap["is_read"] = m.IsRead
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
	m.fieldMap["deleted_at"] = m.DeletedAt
}

func (m messageUser) clone(db *gorm.DB) messageUser {
	m.messageUserDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m messageUser) replaceDB(db *gorm.DB) messageUser {
	m.messageUserDo.ReplaceDB(db)
	return m
}

type messageUserDo struct{ gen.DO }

type IMessageUserDo interface {
	gen.SubQuery
	Debug() IMessageUserDo
	WithContext(ctx context.Context) IMessageUserDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IMessageUserDo
	WriteDB() IMessageUserDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IMessageUserDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IMessageUserDo
	Not(conds ...gen.Condition) IMessageUserDo
	Or(conds ...gen.Condition) IMessageUserDo
	Select(conds ...field.Expr) IMessageUserDo
	Where(conds ...gen.Condition) IMessageUserDo
	Order(conds ...field.Expr) IMessageUserDo
	Distinct(cols ...field.Expr) IMessageUserDo
	Omit(cols ...field.Expr) IMessageUserDo
	Join(table schema.Tabler, on ...field.Expr) IMessageUserDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IMessageUserDo
	RightJoin(table schema.Tabler, on ...field.Expr) IMessageUserDo
	Group(cols ...field.Expr) IMessageUserDo
	Having(conds ...gen.Condition) IMessageUserDo
	Limit(limit int) IMessageUserDo
	Offset(offset int) IMessageUserDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IMessageUserDo
	Unscoped() IMessageUserDo
	Create(values ...*chat_model.MessageUser) error
	CreateInBatches(values []*chat_model.MessageUser, batchSize int) error
	Save(values ...*chat_model.MessageUser) error
	First() (*chat_model.MessageUser, error)
	Take() (*chat_model.MessageUser, error)
	Last() (*chat_model.MessageUser, error)
	Find() ([]*chat_model.MessageUser, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*chat_model.MessageUser, err error)
	FindInBatches(result *[]*chat_model.MessageUser, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*chat_model.MessageUser) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IMessageUserDo
	Assign(attrs ...field.AssignExpr) IMessageUserDo
	Joins(fields ...field.RelationField) IMessageUserDo
	Preload(fields ...field.RelationField) IMessageUserDo
	FirstOrInit() (*chat_model.MessageUser, error)
	FirstOrCreate() (*chat_model.MessageUser, error)
	FindByPage(offset int, limit int) (result []*chat_model.MessageUser, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IMessageUserDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (m messageUserDo) Debug() IMessageUserDo {
	return m.withDO(m.DO.Debug())
}

func (m messageUserDo) WithContext(ctx context.Context) IMessageUserDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m messageUserDo) ReadDB() IMessageUserDo {
	return m.Clauses(dbresolver.Read)
}

func (m messageUserDo) WriteDB() IMessageUserDo {
	return m.Clauses(dbresolver.Write)
}

func (m messageUserDo) Session(config *gorm.Session) IMessageUserDo {
	return m.withDO(m.DO.Session(config))
}

func (m messageUserDo) Clauses(conds ...clause.Expression) IMessageUserDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m messageUserDo) Returning(value interface{}, columns ...string) IMessageUserDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m messageUserDo) Not(conds ...gen.Condition) IMessageUserDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m messageUserDo) Or(conds ...gen.Condition) IMessageUserDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m messageUserDo) Select(conds ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m messageUserDo) Where(conds ...gen.Condition) IMessageUserDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m messageUserDo) Order(conds ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m messageUserDo) Distinct(cols ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m messageUserDo) Omit(cols ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m messageUserDo) Join(table schema.Tabler, on ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m messageUserDo) LeftJoin(table schema.Tabler, on ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m messageUserDo) RightJoin(table schema.Tabler, on ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m messageUserDo) Group(cols ...field.Expr) IMessageUserDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m messageUserDo) Having(conds ...gen.Condition) IMessageUserDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m messageUserDo) Limit(limit int) IMessageUserDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m messageUserDo) Offset(offset int) IMessageUserDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m messageUserDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IMessageUserDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m messageUserDo) Unscoped() IMessageUserDo {
	return m.withDO(m.DO.Unscoped())
}

func (m messageUserDo) Create(values ...*chat_model.MessageUser) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m messageUserDo) CreateInBatches(values []*chat_model.MessageUser, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m messageUserDo) Save(values ...*chat_model.MessageUser) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m messageUserDo) First() (*chat_model.MessageUser, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*chat_model.MessageUser), nil
	}
}

func (m messageUserDo) Take() (*chat_model.MessageUser, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*chat_model.MessageUser), nil
	}
}

func (m messageUserDo) Last() (*chat_model.MessageUser, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*chat_model.MessageUser), nil
	}
}

func (m messageUserDo) Find() ([]*chat_model.MessageUser, error) {
	result, err := m.DO.Find()
	return result.([]*chat_model.MessageUser), err
}

func (m messageUserDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*chat_model.MessageUser, err error) {
	buf := make([]*chat_model.MessageUser, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m messageUserDo) FindInBatches(result *[]*chat_model.MessageUser, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m messageUserDo) Attrs(attrs ...field.AssignExpr) IMessageUserDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m messageUserDo) Assign(attrs ...field.AssignExpr) IMessageUserDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m messageUserDo) Joins(fields ...field.RelationField) IMessageUserDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m messageUserDo) Preload(fields ...field.RelationField) IMessageUserDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m messageUserDo) FirstOrInit() (*chat_model.MessageUser, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*chat_model.MessageUser), nil
	}
}

func (m messageUserDo) FirstOrCreate() (*chat_model.MessageUser, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*chat_model.MessageUser), nil
	}
}

func (m messageUserDo) FindByPage(offset int, limit int) (result []*chat_model.MessageUser, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m messageUserDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m messageUserDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m messageUserDo) Delete(models ...*chat_model.MessageUser) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *messageUserDo) withDO(do gen.Dao) *messageUserDo {
	m.DO = *do.(*gen.DO)
	return m
}
