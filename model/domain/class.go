package domain

import "rxdrag.com/entify/model/meta"

type Class struct {
	Uuid          string
	InnerId       uint64
	StereoType    string
	Name          string
	Description   string
	Root          bool
	SoftDelete    bool
	Attributes    []*Attribute
	Methods       []*Method
	parents       []*Class
	Children      []*Class
	AppId         uint64
	IdNoShift     bool
	QueryScript   string
	MuationScript string
	EnableVersion bool
	PackageUuid   string
}

func NewClass(c *meta.ClassMeta) *Class {
	cls := Class{
		Uuid:          c.Uuid,
		InnerId:       c.InnerId,
		StereoType:    c.StereoType,
		Name:          c.Name,
		Description:   c.Description,
		Root:          c.Root,
		SoftDelete:    c.SoftDelete,
		Attributes:    make([]*Attribute, len(c.Attributes)),
		Methods:       make([]*Method, len(c.Methods)),
		parents:       []*Class{},
		Children:      []*Class{},
		QueryScript:   c.QueryScript,
		MuationScript: c.MuationScript,
		EnableVersion: c.EnableVersion,
		AppId:         c.AppId,
		PackageUuid:   c.PackageUuid,
		IdNoShift:     c.IdNoShift,
	}

	for i := range c.Attributes {
		cls.Attributes[i] = NewAttribute(&c.Attributes[i], &cls)
	}

	for i := range c.Methods {
		cls.Methods[i] = NewMethod(&c.Methods[i], &cls)
	}

	return &cls
}

func (c *Class) HasChildren() bool {
	return len(c.Children) > 0
}

func (c *Class) AllParents() []*Class {
	parents := []*Class{}
	for i := range c.parents {
		parent := c.parents[i]
		parents = append(parents, parent)
		parents = append(parents, parent.AllParents()...)
	}

	return parents
}
