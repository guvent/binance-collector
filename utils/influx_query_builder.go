package utils

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type KeyValueItem struct {
	Key   string
	Value interface{}
}

type InfluxQueryBuilder struct {
	Queries []string
}

func (builder *InfluxQueryBuilder) From(bucket string) *InfluxQueryBuilder {
	query := fmt.Sprintf("from(bucket: \"%s\")", bucket)

	builder.Queries = []string{query}

	return builder
}

func (builder *InfluxQueryBuilder) Range(start, stop string) *InfluxQueryBuilder {
	query := fmt.Sprintf("|> range(start: %s)", start)

	if stop != "" {
		query = fmt.Sprintf("|> range(start: %s, stop: %s)", start, stop)
	}

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Filter(item KeyValueItem) *InfluxQueryBuilder {
	query := fmt.Sprintf("|> filter(fn: (r) => r[\"%s\"] == \"%s\")", item.Key, item.Value)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Filters(items []KeyValueItem, operand string) *InfluxQueryBuilder {
	query := "|> filter( fn: (r) =>"
	opCode := ""

	for _, item := range items {
		query = fmt.Sprintf("%s%s r[\"%s\"] == \"%s\" ", query, opCode, item.Key, item.Value)
		opCode = fmt.Sprintf("%s", operand)
	}

	query = fmt.Sprintf("%s)", query)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Group(fields []string) *InfluxQueryBuilder {
	query := fmt.Sprintf("|> group(columns: [")
	comma := ""

	for _, field := range fields {
		query = fmt.Sprintf("%s%s\"%s\"", query, comma, field)
		comma = ", "
	}

	query = fmt.Sprintf("%s])", query)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Pivot(rowKey, columnKey, valueColumn string) *InfluxQueryBuilder {
	query := fmt.Sprintf(
		"|> pivot(rowKey: [\"%s\"], columnKey: [\"%s\"], valueColumn: \"%s\")", rowKey, columnKey, valueColumn)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Map(expression string) *InfluxQueryBuilder {
	query := fmt.Sprintf("|> map(fn: (r) => (%s))", expression)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Sort(fields []string, isDesc bool) *InfluxQueryBuilder {
	query := fmt.Sprintf("|> sort(columns: [")
	desc := ""
	comma := ""

	for _, field := range fields {
		query = fmt.Sprintf("%s%s\"%s\"", query, comma, field)
		comma = ", "
	}

	if isDesc {
		desc = ", desc: true"
	}

	query = fmt.Sprintf("%s]%s)", query, desc)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Limit(n int) *InfluxQueryBuilder {
	query := fmt.Sprintf("|> limit(n:%d)", n)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) Yield(name string) *InfluxQueryBuilder {
	query := fmt.Sprintf("  |> yield(name: \"%s\")", name)

	builder.Queries = append(builder.Queries, query)

	return builder
}

func (builder *InfluxQueryBuilder) ToString() string {
	queryString := ""
	whitespace := ""

	for _, query := range builder.Queries {
		queryString = fmt.Sprintf("%s%s%s\n", queryString, whitespace, query)
		whitespace = "  "
	}

	return queryString
}

func (builder *InfluxQueryBuilder) Execute(
	queryApi api.QueryAPI,
	onNext func(result *api.QueryTableResult),
) error {
	if result, err := queryApi.Query(context.Background(), builder.ToString()); err != nil {
		return err
	} else {
		for result.Next() {
			onNext(result)
		}
	}
	return nil
}
