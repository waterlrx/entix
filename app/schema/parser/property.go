package parser

import (
	"github.com/graphql-go/graphql"
	"rxdrag.com/entify/model/graph"
	"rxdrag.com/entify/model/meta"
	"rxdrag.com/entify/scalars"
)

func (p *ModelParser) InputPropertyType(property graph.Propertier) graphql.Type {
	if property.GetType() == meta.FILE {
		return scalars.UploadType
	}
	return p.PropertyType(property.GetType())
}

func (p *ModelParser) PropertyType(propType string) graphql.Output {
	switch propType {
	case meta.ID:
		return graphql.ID
	case meta.INT:
		return graphql.Int
	case meta.FLOAT:
		return graphql.Float
	case meta.BOOLEAN:
		return graphql.Boolean
	case meta.STRING:
		return graphql.String
	case meta.DATE:
		return graphql.DateTime
	case
		meta.JSON,
		meta.VALUE_OBJECT,
		meta.ENTITY,
		meta.ID_ARRAY,
		meta.INT_ARRAY,
		meta.FLOAT_ARRAY,
		meta.STRING_ARRAY,
		meta.DATE_ARRAY,
		meta.ENUM_ARRAY,
		meta.VALUE_OBJECT_ARRAY,
		meta.ENTITY_ARRAY:
		return scalars.JSONType
	case meta.ENUM:
		// 方便输入，改为字符串
		// enum := property.GetEumnType()
		// if enum == nil {
		// 	panic("Can not find enum entity")
		// }
		// return p.EnumType(enum.Name)
		return graphql.String
	case meta.FILE:
		//return graphql.String
		return fileOutputType
	}

	panic("No column type:" + propType)
}

func (p *ModelParser) AttributeExp(column *graph.Attribute) *graphql.InputObjectFieldConfig {
	switch column.Type {
	case meta.INT:
		return &IntComparisonExp
	case meta.FLOAT:
		return &FloatComparisonExp
	case meta.BOOLEAN:
		return &BooleanComparisonExp
	case meta.STRING:
		return &StringComparisonExp
	case meta.DATE:
		return &DateTimeComparisonExp
	case
		meta.JSON,
		meta.VALUE_OBJECT,
		meta.ENTITY,
		meta.ID_ARRAY,
		meta.INT_ARRAY,
		meta.FLOAT_ARRAY,
		meta.STRING_ARRAY,
		meta.DATE_ARRAY,
		meta.ENUM_ARRAY,
		meta.VALUE_OBJECT_ARRAY,
		meta.ENTITY_ARRAY,
		meta.FILE:
		return nil
	case meta.ID:
		return &IdComparisonExp
	case meta.ENUM:
		return p.EnumComparisonExp(column)
	}

	panic("No column type: " + column.Type)
}

func (p *ModelParser) AttributeOrderBy(column *graph.Attribute) *graphql.Enum {
	switch column.Type {
	case
		meta.JSON,
		meta.VALUE_OBJECT,
		meta.BOOLEAN,
		meta.ENTITY,
		meta.ID_ARRAY,
		meta.INT_ARRAY,
		meta.FLOAT_ARRAY,
		meta.STRING_ARRAY,
		meta.DATE_ARRAY,
		meta.ENUM_ARRAY,
		meta.VALUE_OBJECT_ARRAY,
		meta.ENTITY_ARRAY:
		return nil
	}

	return EnumOrderBy
}