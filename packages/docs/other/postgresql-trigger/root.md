# PostgreSQL 触发器

触发器是 PostgreSQL 里的一种自动执行机制：当表或视图上发生 `INSERT`、`UPDATE`、`DELETE`、`TRUNCATE` 等事件时，数据库会自动调用指定函数。

它适合放那些“必须始终执行”的数据库规则，比如：

- 自动维护 `updated_at`
- 写审计日志
- 校验或修正写入数据
- 在主表写入后自动补齐关联数据

<!-- truncate -->

## 基本结构

PostgreSQL 触发器通常由两部分组成：

1. 触发器函数
2. 触发器定义

最常见写法：

```sql
CREATE OR REPLACE FUNCTION trigger_function_name()
RETURNS TRIGGER AS $$
BEGIN
	-- 处理逻辑
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_name
BEFORE INSERT OR UPDATE ON table_name
FOR EACH ROW
EXECUTE FUNCTION trigger_function_name();
```

## 关键字解释

先看函数部分：

- `CREATE`：创建一个数据库对象，这里创建的是函数。
- `OR REPLACE`：如果同名函数已经存在，就用新的定义覆盖它。
- `FUNCTION`：说明要创建的是函数。
- `trigger_function_name()`：函数名，后面的 `()` 是参数列表，这里表示无参数。
- `RETURNS`：声明函数返回值类型。
- `TRIGGER`：表示这个函数返回触发器类型，能被触发器调用。
- `AS`：后面开始写函数体。
- `$$`：函数体分隔符，对应一整段被包起来的函数代码。
- `BEGIN`：函数逻辑开始。
- `RETURN NEW`：返回新行数据，常见于 `BEFORE INSERT` 和 `BEFORE UPDATE`。
- `END`：函数逻辑结束。
- `LANGUAGE`：声明函数使用的语言。
- `plpgsql`：表示函数体使用 PostgreSQL 的过程式语言 `plpgsql` 编写。

再看触发器定义部分：

- `CREATE TRIGGER`：创建触发器。
- `trigger_name`：触发器名称。
- `BEFORE`：在数据真正写入前执行。
- `INSERT OR UPDATE`：在插入或更新时触发。
- `ON table_name`：把触发器挂到指定表上。
- `FOR EACH ROW`：每影响一行就执行一次。
- `EXECUTE FUNCTION`：满足触发条件后执行指定函数。

可以把这套结构理解成两层：

- 函数定义“要做什么”
- 触发器定义“什么时候、在哪张表上执行”

## 触发器怎么分类

### 按执行时机分

- `BEFORE`：在数据真正写入前执行，适合校验、补默认值、修改 `NEW`
- `AFTER`：在数据写入后执行，适合写日志、补关联表数据
- `INSTEAD OF`：主要用于视图，替代原本的写入动作

### 按触发事件分

- `INSERT`
- `UPDATE`
- `DELETE`
- `TRUNCATE`

### 按执行粒度分

- `FOR EACH ROW`：每一行都会触发一次，最常用
- `FOR EACH STATEMENT`：整条 SQL 只触发一次，不关心影响了多少行

## 触发器里常用变量

在行级触发器里，最常用的是这几个：

- `NEW`：新数据，常见于 `INSERT` 和 `UPDATE`
- `OLD`：旧数据，常见于 `UPDATE` 和 `DELETE`
- `TG_OP`：当前操作类型，比如 `INSERT`、`UPDATE`、`DELETE`

常见规则：

- `BEFORE INSERT` / `BEFORE UPDATE` 通常返回 `NEW`
- `BEFORE DELETE` 通常返回 `OLD`
- `AFTER` 触发器一般也要返回一个值，但通常不会再影响已写入的数据

## 最小示例

下面这个例子会在更新数据时自动刷新 `updated_at`：

```sql
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = NOW();
	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
```

效果是：只要 `users` 表发生更新，数据库就会自动把当前行的 `updated_at` 改成当前时间。

## 一个更贴近业务的例子

如果你想在创建订单后自动写一条审计日志，可以这样做：

```sql
CREATE OR REPLACE FUNCTION log_order_created()
RETURNS TRIGGER AS $$
BEGIN
	INSERT INTO order_logs (order_id, action, created_at)
	VALUES (NEW.id, 'created', NOW());

	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER orders_after_insert
AFTER INSERT ON orders
FOR EACH ROW
EXECUTE FUNCTION log_order_created();
```

这种场景更适合 `AFTER INSERT`，因为它依赖订单已经成功写入。

## 查看和删除

查看某张表上的触发器：

```sql
\dS table_name
```

删除触发器：

```sql
DROP TRIGGER trigger_name ON table_name;
```

删除触发器函数：

```sql
DROP FUNCTION function_name();
```

临时禁用触发器：

```sql
ALTER TABLE table_name DISABLE TRIGGER trigger_name;
ALTER TABLE table_name ENABLE TRIGGER trigger_name;
```

## 什么时候适合用

适合放进触发器的逻辑通常有两个特点：

- 它是数据一致性规则，不应该依赖某个应用入口记得执行
- 它离数据库很近，用 SQL 处理更稳定

比如：

- 自动维护时间戳
- 审计日志
- 同步统计字段
- 创建主记录后补默认从表数据

## 不要滥用

触发器很方便，但也有成本：

- 隐式执行，排查问题时不如应用代码直观
- 写得太重会拖慢写入性能
- 多个触发器叠加后，执行顺序和副作用会变复杂

实践上建议：

- 只放强一致性、数据库层必须兜底的逻辑
- 函数保持短小，避免塞太多业务流程
- 重要触发器要配文档和测试
- 能用约束表达的规则，优先考虑约束而不是触发器

## 相关链接

- [PostgreSQL 官方文档: Trigger Definition](https://www.postgresql.org/docs/current/sql-createtrigger.html)
- [PostgreSQL 官方文档: Trigger Functions](https://www.postgresql.org/docs/current/plpgsql-trigger.html)
- [workspace-kit 中触发器的实际用法](../../rust/workspace-kit/wk-1-init-project/postgresql-trigger/root.md)
